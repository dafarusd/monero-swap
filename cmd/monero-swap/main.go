// monero-swap runs one side of a Monero <-> ETH atomic swap against the XmrSwap contract.
//
//	monero-swap maker run   keeps an offer live and serves every buyer that takes it
//	monero-swap taker take  pays ETH into an offer and collects the Monero
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
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"

	"github.com/athanorlabs/atomic-swap/common"
	mcrypto "github.com/athanorlabs/atomic-swap/crypto/monero"
	"github.com/athanorlabs/atomic-swap/monero"
	"github.com/athanorlabs/atomic-swap/swap2"
)

func main() {
	app := &cli.App{
		Name:  "monero-swap",
		Usage: "trustless Monero <-> ETH swaps through the XmrSwap contract",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "env", Value: "stagenet", Usage: "stagenet or mainnet (Monero network)"},
			&cli.StringFlag{Name: "eth-rpc", Required: true, Usage: "Ethereum-compatible JSON-RPC URL"},
			&cli.StringFlag{Name: "contract", Required: true, Usage: "XmrSwap contract address"},
			&cli.StringFlag{Name: "eth-key", Usage: "file with the hex private key that pays gas (default DATA_DIR/eth.key, created if missing)"},
			&cli.StringFlag{Name: "data-dir", Value: filepath.Join(os.Getenv("HOME"), ".monero-swap"), Usage: "state, keys and wallet files"},
			&cli.StringFlag{Name: "monerod-host", Required: true, Usage: "Monero node host"},
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
						&cli.StringFlag{Name: "payout", Required: true, Usage: "ETH address that receives payments (your own wallet)"},
						&cli.StringFlag{Name: "min", Value: "0.02", Usage: "smallest payment in ETH"},
						&cli.StringFlag{Name: "max", Value: "0.05", Usage: "largest payment in ETH"},
						&cli.StringFlag{Name: "price", Required: true, Usage: "XMR the buyer gets per 1 ETH, e.g. 15.5"},
						&cli.DurationFlag{Name: "offer-ttl", Value: 24 * time.Hour, Usage: "how long each offer stays open"},
						&cli.DurationFlag{Name: "t1", Value: time.Hour, Usage: "time the buyer has to confirm the Monero before you may claim anyway"},
						&cli.DurationFlag{Name: "t2", Value: time.Hour, Usage: "time after t1 before the buyer may refund regardless"},
						&cli.BoolFlag{Name: "manual-xmr", Usage: "do not hold Monero in the program; print the address and pay from your own wallet"},
					},
					Action: runMaker,
				}},
			},
			{
				Name:  "taker",
				Usage: "buy Monero with ETH",
				Subcommands: []*cli.Command{
					{
						Name:  "take",
						Usage: "take an offer and run the swap to the end",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "offer", Required: true, Usage: "offer id (hex)"},
							&cli.StringFlag{Name: "amount", Required: true, Usage: "ETH to pay, e.g. 0.02"},
							&cli.StringFlag{Name: "payout-xmr", Usage: "Monero address that receives the coins (default: the program's wallet)"},
						},
						Action: runTake,
					},
					{
						Name:  "resume",
						Usage: "continue a swap after a restart",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "swap", Required: true, Usage: "swap id (hex)"},
							&cli.StringFlag{Name: "payout-xmr"},
						},
						Action: runResume,
					},
				},
			},
			{Name: "offers", Usage: "list open offers on the contract", Flags: []cli.Flag{&cli.Uint64Flag{Name: "blocks", Value: 50000, Usage: "how far back to scan"}}, Action: runOffers},
			{Name: "status", Usage: "show this program's offers and swaps", Action: runStatus},
			{Name: "addresses", Usage: "show this program's ETH and Monero addresses", Action: runAddresses},
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
	chain  *swap2.Chain
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
			log.Printf("created new gas wallet %s (key in %s); fund it with a little ETH", k, keyFile)
		}
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	chain, err := swap2.NewChain(ctx, c.String("eth-rpc"), ethcommon.HexToAddress(c.String("contract")), keyFile)
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
		wallet, err = monero.NewWalletClient(&monero.WalletClientConf{
			Env:            e,
			WalletFilePath: filepath.Join(dataDir, "wallet", "swap-wallet"),
			WalletPassword: c.String("xmr-wallet-password"),
			MonerodNodes:   []*common.MoneroNode{{Host: c.String("monerod-host"), Port: c.Uint("monerod-port")}},
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
		return err
	}
	maxWei, err := swap2.ParseWei(c.String("max"))
	if err != nil {
		return err
	}
	price, err := swap2.ParsePiconero(c.String("price"))
	if err != nil {
		return err
	}
	m := swap2.NewMaker(swap2.MakerConfig{
		Env:         e.env,
		Chain:       e.chain,
		Wallet:      e.wallet,
		Store:       e.store,
		Payout:      ethcommon.HexToAddress(c.String("payout")),
		Asset:       ethcommon.Address{},
		MinWei:      minWei,
		MaxWei:      maxWei,
		XmrPerAsset: new(big.Int).SetUint64(price),
		OfferTTL:    c.Duration("offer-ttl"),
		Timeout1:    c.Duration("t1"),
		Timeout2:    c.Duration("t2"),
		ManualXMR:   c.Bool("manual-xmr"),
	})
	err = m.Run(e.ctx)
	if err == context.Canceled {
		return nil
	}
	return err
}

func takerFrom(c *cli.Context, e *env) (*swap2.Taker, error) {
	var payout *mcrypto.Address
	if p := c.String("payout-xmr"); p != "" {
		var err error
		payout, err = mcrypto.NewAddress(p, e.env)
		if err != nil {
			return nil, fmt.Errorf("payout-xmr: %w", err)
		}
	}
	return swap2.NewTaker(swap2.TakerConfig{Env: e.env, Chain: e.chain, Wallet: e.wallet, Store: e.store, PayoutXMR: payout}), nil
}

func runTake(c *cli.Context) error {
	e, err := setup(c, true)
	if err != nil {
		return err
	}
	defer e.close()
	t, err := takerFrom(c, e)
	if err != nil {
		return err
	}
	amount, err := swap2.ParseWei(c.String("amount"))
	if err != nil {
		return err
	}
	id, err := hex.DecodeString(strings.TrimPrefix(c.String("offer"), "0x"))
	if err != nil || len(id) != 32 {
		return fmt.Errorf("--offer must be a 32-byte hex id")
	}
	var offerID [32]byte
	copy(offerID[:], id)
	swapID, err := t.Take(e.ctx, offerID, amount)
	if swapID != "" {
		log.Printf("swap id: %s", swapID)
	}
	return err
}

func runResume(c *cli.Context) error {
	e, err := setup(c, true)
	if err != nil {
		return err
	}
	defer e.close()
	t, err := takerFrom(c, e)
	if err != nil {
		return err
	}
	return t.Resume(e.ctx, strings.TrimPrefix(c.String("swap"), "0x"))
}

func runOffers(c *cli.Context) error {
	e, err := setup(c, false)
	if err != nil {
		return err
	}
	defer e.close()
	head, now, err := e.chain.Head(e.ctx)
	if err != nil {
		return err
	}
	from := uint64(0)
	if head > c.Uint64("blocks") {
		from = head - c.Uint64("blocks")
	}
	logs, err := e.chain.Logs(e.ctx, from, head, [][]ethcommon.Hash{{e.chain.EventID("OfferPosted")}})
	if err != nil {
		return err
	}
	n := 0
	for _, l := range logs {
		ev, err := e.chain.Contract.ParseOfferPosted(l)
		if err != nil {
			continue
		}
		o, err := e.chain.Contract.Offers(&bind.CallOpts{Context: e.ctx}, ev.OfferId)
		if err != nil || !o.Active || o.Expiry <= now {
			continue
		}
		n++
		fmt.Printf("offer %s\n  sells XMR for %s\n  buyer pays %s-%s, gets %s XMR per unit\n  expires %s, t1 %s, t2 %s, maker %s\n",
			hex.EncodeToString(ev.OfferId[:]), assetLabel(o.Asset),
			swap2.FmtETH(o.MinAmount), swap2.FmtETH(o.MaxAmount), swap2.FmtXMR(o.XmrPerAsset.Uint64()),
			time.Unix(int64(o.Expiry), 0).Format(time.RFC3339), time.Duration(o.Timeout1Duration)*time.Second,
			time.Duration(o.Timeout2Duration)*time.Second, o.Maker.Hex())
	}
	if n == 0 {
		fmt.Println("no open offers")
	}
	return nil
}

func assetLabel(a ethcommon.Address) string {
	if a == (ethcommon.Address{}) {
		return "ETH"
	}
	return "token " + a.Hex()
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
	ethBal, err := e.chain.EC.BalanceAt(e.ctx, e.chain.From, nil)
	if err != nil {
		return err
	}
	fmt.Printf("ETH gas wallet:  %s  (%s ETH)\nMonero wallet:   %s  (%s XMR, %s unlocked)\n",
		e.chain.From.Hex(), swap2.FmtETH(ethBal), e.wallet.PrimaryAddress(), swap2.FmtXMR(bal.Balance), swap2.FmtXMR(bal.UnlockedBalance))
	return nil
}
