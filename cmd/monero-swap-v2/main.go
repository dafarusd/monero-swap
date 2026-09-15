// monero-swap-v2 sells Monero for ETH against the XmrSwapV2 contract.
//
//	monero-swap-v2 maker run   keeps an offer live and serves every buyer that takes it
//	monero-swap-v2 offers      lists the open offers, read straight from contract storage
//	monero-swap-v2 withdraw    collects bonds and fees the contract is holding for you
//
// The buyer side is the web page; there is no taker command here. v1's seller is the separate
// `monero-swap` binary and still drives the v1 contract.
package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/urfave/cli/v2"

	"github.com/athanorlabs/atomic-swap/common"
	"github.com/athanorlabs/atomic-swap/monero"
	"github.com/athanorlabs/atomic-swap/swap2"
	"github.com/athanorlabs/atomic-swap/swap2v2"
)

func main() {
	app := &cli.App{
		Name:  "monero-swap-v2",
		Usage: "sell Monero for ETH through the XmrSwapV2 contract",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "env", Value: "mainnet", Usage: "stagenet or mainnet (Monero network)"},
			&cli.StringFlag{Name: "eth-rpc", Required: true, Usage: "Ethereum-compatible JSON-RPC URL"},
			&cli.StringFlag{Name: "contract", Required: true, Usage: "XmrSwapV2 contract address"},
			&cli.StringFlag{Name: "eth-key", Usage: "file with the hex private key that pays gas and the bond (default DATA_DIR/eth.key, created if missing)"},
			&cli.StringFlag{Name: "data-dir", Value: filepath.Join(os.Getenv("HOME"), ".monero-swap-v2"), Usage: "state, keys and wallet files"},
			&cli.StringFlag{Name: "monerod-host", Usage: "Monero node host (required for maker and addresses)"},
			&cli.UintFlag{Name: "monerod-port", Value: 18081, Usage: "Monero node RPC port (stagenet nodes often use 38081)"},
			&cli.StringFlag{Name: "xmr-wallet-password", Value: "", Usage: "password of the program's Monero wallet"},
		},
		Commands: []*cli.Command{
			{
				Name:  "maker",
				Usage: "sell Monero",
				Subcommands: []*cli.Command{{
					Name:  "run",
					Usage: "keep one offer live and serve takers (Ctrl-C to stop; swaps resume on restart)",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "payout", Required: true, Usage: "ETH address that receives payments (your own wallet, never the gas key)"},
						&cli.StringFlag{Name: "min", Value: "0.02", Usage: "smallest payment, in ETH"},
						&cli.StringFlag{Name: "max", Value: "0.05", Usage: "largest payment, in ETH. The bond is a percentage of THIS."},
						&cli.StringFlag{Name: "price", Required: true, Usage: "XMR the buyer gets per 1 ETH, e.g. 15.5"},
						&cli.DurationFlag{Name: "offer-ttl", Value: 24 * time.Hour, Usage: "how long each offer stays open"},
						&cli.DurationFlag{Name: "t1", Value: 24 * time.Hour, Usage: "time the buyer has to confirm the Monero before you may claim anyway (contract minimum applies)"},
						&cli.DurationFlag{Name: "t2", Value: 24 * time.Hour, Usage: "time after t1 before the buyer may refund regardless (contract minimum applies)"},
						&cli.BoolFlag{Name: "manual-xmr", Usage: "do not hold Monero in the program; print the address and pay from your own wallet"},
						&cli.StringFlag{Name: "min-withdraw", Value: "0.0001", Usage: "smallest credited balance worth spending gas to collect, in ETH"},
					},
					Action: runMaker,
				}},
			},
			{Name: "offers", Usage: "list open offers, read straight from contract storage", Action: runOffers},
			{Name: "terms", Usage: "show the contract's immutable fee, split, bond and timeout", Action: runTerms},
			{Name: "withdraw", Usage: "collect bonds and fees the contract is holding for you", Action: runWithdraw},
			{Name: "status", Usage: "show this program's offers and swaps", Action: runStatus},
			{Name: "addresses", Usage: "show this program's ETH and Monero addresses and balances", Action: runAddresses},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type env struct {
	ctx    context.Context
	cancel context.CancelFunc
	env    common.Environment
	chain  *swap2v2.Chain
	wallet monero.WalletClient
	store  *swap2.Store
}

func setup(c *cli.Context, needWallet bool) (*env, error) {
	dataDir := c.String("data-dir")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, err
	}
	var e common.Environment
	switch c.String("env") {
	case "stagenet":
		e = common.Stagenet
	case "mainnet":
		e = common.Mainnet
	default:
		return nil, fmt.Errorf("--env must be stagenet or mainnet")
	}
	keyFile := c.String("eth-key")
	if keyFile == "" {
		keyFile = filepath.Join(dataDir, "eth.key")
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			k, err := swap2.NewEthKeyFile(keyFile)
			if err != nil {
				return nil, err
			}
			log.Printf("created new gas wallet %s (key in %s)", k, keyFile)
			log.Printf("fund it with ETH for gas AND for the offer bond — see `monero-swap-v2 terms`")
		}
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	chain, err := swap2v2.NewChain(ctx, c.String("eth-rpc"), ethcommon.HexToAddress(c.String("contract")), keyFile)
	if err != nil {
		cancel()
		return nil, err
	}
	store, err := swap2.OpenStore(filepath.Join(dataDir, "state.json"))
	if err != nil {
		cancel()
		return nil, err
	}
	var wallet monero.WalletClient
	if needWallet {
		host := c.String("monerod-host")
		if host == "" {
			cancel()
			return nil, fmt.Errorf("--monerod-host is required for this command")
		}
		wallet, err = monero.NewWalletClient(&monero.WalletClientConf{
			Env:            e,
			WalletFilePath: filepath.Join(dataDir, "wallet", "swap-wallet"),
			WalletPassword: c.String("xmr-wallet-password"),
			MonerodNodes:   []*common.MoneroNode{{Host: host, Port: c.Uint("monerod-port")}},
		})
		if err != nil {
			cancel()
			return nil, fmt.Errorf("monero wallet: %w", err)
		}
	}
	return &env{ctx: ctx, cancel: cancel, env: e, chain: chain, wallet: wallet, store: store}, nil
}

func (e *env) close() {
	if e.wallet != nil {
		e.wallet.Close()
	}
	e.cancel()
}

func runMaker(c *cli.Context) error {
	e, err := setup(c, true)
	if err != nil {
		return err
	}
	defer e.close()

	minWei, err := swap2.ParseWei(c.String("min"))
	if err != nil {
		return fmt.Errorf("--min: %w", err)
	}
	maxWei, err := swap2.ParseWei(c.String("max"))
	if err != nil {
		return fmt.Errorf("--max: %w", err)
	}
	if maxWei.Cmp(minWei) < 0 {
		return fmt.Errorf("--max must be at least --min")
	}
	minWithdraw, err := swap2.ParseWei(c.String("min-withdraw"))
	if err != nil {
		return fmt.Errorf("--min-withdraw: %w", err)
	}
	price, err := swap2.ParsePiconero(c.String("price"))
	if err != nil {
		return fmt.Errorf("--price: %w", err)
	}
	// The contract prices in piconero per 1e18 wei, which for ETH is piconero per whole ETH.
	xmrPerAsset := new(big.Int).SetUint64(price)

	log.Printf("selling XMR for ETH: %s to %s per swap, %s XMR per 1 ETH",
		swap2.FmtETH(minWei), swap2.FmtETH(maxWei), swap2.FmtXMR(price))

	m := swap2v2.NewMaker(swap2v2.MakerConfig{
		Env:         e.env,
		Chain:       e.chain,
		Wallet:      e.wallet,
		Store:       e.store,
		Payout:      ethcommon.HexToAddress(c.String("payout")),
		MinWei:      minWei,
		MaxWei:      maxWei,
		XmrPerAsset: xmrPerAsset,
		OfferTTL:    c.Duration("offer-ttl"),
		Timeout1:    c.Duration("t1"),
		Timeout2:    c.Duration("t2"),
		ManualXMR:   c.Bool("manual-xmr"),
		MinWithdraw: minWithdraw,
	})
	err = m.Run(e.ctx)
	if err == context.Canceled {
		return nil
	}
	return err
}

func runTerms(c *cli.Context) error {
	e, err := setup(c, false)
	if err != nil {
		return err
	}
	defer e.close()
	t, err := e.chain.Terms(e.ctx)
	if err != nil {
		return err
	}
	opts := &bind.CallOpts{Context: e.ctx}
	share, err := e.chain.Contract.FounderShareBps(opts)
	if err != nil {
		return err
	}
	founder, err := e.chain.Contract.Founder(opts)
	if err != nil {
		return err
	}
	devFund, err := e.chain.Contract.DevFund(opts)
	if err != nil {
		return err
	}
	oneEth, _ := swap2.ParseWei("1")
	bond1, err := e.chain.BondFor(e.ctx, oneEth)
	if err != nil {
		return err
	}
	fmt.Printf(`contract %s on chain %s

  fee              %d bps (%.2f%%) of the buyer's payment, taken when you claim
  fee split        %d/%d between %s and %s
  bond             %d bps of an offer's ceiling — %s ETH on a 1 ETH ceiling
  minimum timeout  %s for each of the two windows

The bond is always returned, whichever way the swap ends. It is returned as a credit
inside the contract, not as a transfer: run 'withdraw', or let 'maker run' collect it.
These values are immutable and cannot be changed by anyone, including the deployer.
`,
		e.chain.Addr.Hex(), e.chain.ChainID,
		t.FeeBps, float64(t.FeeBps)/100,
		share, 10000-share, founder.Hex(), devFund.Hex(),
		t.BondBps, swap2.FmtETH(bond1),
		t.MinTimeout)
	return nil
}

func runOffers(c *cli.Context) error {
	e, err := setup(c, false)
	if err != nil {
		return err
	}
	defer e.close()
	_, now, err := e.chain.Head(e.ctx)
	if err != nil {
		return err
	}
	opts := &bind.CallOpts{Context: e.ctx}
	total, err := e.chain.Contract.OfferCount(opts)
	if err != nil {
		return err
	}
	n := 0
	const page = 50
	for off := int64(0); off < total.Int64(); off += page {
		offers, err := e.chain.Contract.ListOffers(opts, big.NewInt(off), big.NewInt(page))
		if err != nil {
			return err
		}
		for _, o := range offers {
			if !o.Active || o.Expiry <= now {
				continue
			}
			price := o.XmrPerAsset
			priced := "fixed"
			if o.Oracle != (ethcommon.Address{}) {
				p, perr := e.chain.Contract.PriceOf(opts, o.Id)
				if perr != nil {
					fmt.Printf("offer %s\n  priced from a feed that will not resolve right now (%s); not takeable\n",
						hex.EncodeToString(o.Id[:]), perr)
					continue
				}
				price, priced = p, "feed "+o.Oracle.Hex()
			}
			n++
			fmt.Printf("offer %s\n  buyer pays %s to %s ETH, gets %s XMR per 1 ETH (%s)\n  expires %s, t1 %s, t2 %s\n  maker %s, bond %s ETH\n",
				hex.EncodeToString(o.Id[:]),
				swap2.FmtETH(o.MinAmount), swap2.FmtETH(o.MaxAmount), swap2.FmtXMR(price.Uint64()), priced,
				time.Unix(int64(o.Expiry), 0).Format(time.RFC3339),
				time.Duration(o.Timeout1Duration)*time.Second, time.Duration(o.Timeout2Duration)*time.Second,
				o.Maker.Hex(), swap2.FmtETH(o.Bond))
		}
	}
	if n == 0 {
		fmt.Println("no open offers")
	}
	return nil
}

func runWithdraw(c *cli.Context) error {
	e, err := setup(c, false)
	if err != nil {
		return err
	}
	defer e.close()
	owed, err := e.chain.Owed(e.ctx)
	if err != nil {
		return err
	}
	if owed.Sign() == 0 {
		fmt.Printf("the contract is holding nothing for %s\n", e.chain.From.Hex())
		return nil
	}
	fmt.Printf("collecting %s ETH for %s\n", swap2.FmtETH(owed), e.chain.From.Hex())
	_, err = e.chain.Send(e.ctx, "withdraw", nil, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return e.chain.Contract.Withdraw(opts)
	})
	return err
}

func runStatus(c *cli.Context) error {
	e, err := setup(c, false)
	if err != nil {
		return err
	}
	defer e.close()
	e.store.View(func(st *swap2.State) {
		type view struct {
			Offers map[string]*swap2.OfferRec `json:"offers"`
			Swaps  map[string]*swap2.SwapRec  `json:"swaps"`
		}
		v := view{Offers: map[string]*swap2.OfferRec{}, Swaps: map[string]*swap2.SwapRec{}}
		for k, o := range st.Offers {
			cp := *o
			cp.SpendPriv, cp.ViewPriv = "<hidden>", "<hidden>"
			v.Offers[k] = &cp
		}
		for k, s := range st.Swaps {
			cp := *s
			cp.MySpendPriv, cp.MyViewPriv, cp.TheirViewPriv = "<hidden>", "<hidden>", "<hidden>"
			v.Swaps[k] = &cp
		}
		out, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(out))
	})
	return nil
}

func runAddresses(c *cli.Context) error {
	e, err := setup(c, true)
	if err != nil {
		return err
	}
	defer e.close()
	bal, err := e.wallet.GetBalance(0)
	if err != nil {
		return err
	}
	ethBal, err := e.chain.Balance(e.ctx)
	if err != nil {
		return err
	}
	owed, err := e.chain.Owed(e.ctx)
	if err != nil {
		return err
	}
	fmt.Printf("ETH gas+bond wallet:  %s  (%s ETH)\nHeld by the contract: %s ETH (bonds and fees waiting; run `withdraw`)\nMonero wallet:        %s  (%s XMR, %s unlocked)\n",
		e.chain.From.Hex(), swap2.FmtETH(ethBal), swap2.FmtETH(owed),
		e.wallet.PrimaryAddress(), swap2.FmtXMR(bal.Balance), swap2.FmtXMR(bal.UnlockedBalance))
	return nil
}
