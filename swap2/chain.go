package swap2

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/athanorlabs/atomic-swap/ethereum/block"
	"github.com/athanorlabs/atomic-swap/xmrswap"
)

// logChunk keeps eth_getLogs queries small enough for public RPC endpoints.
const logChunk = 2000

// Chain wraps an Ethereum client, a signing key, and the XmrSwap contract.
type Chain struct {
	EC       *ethclient.Client
	Contract *xmrswap.XmrSwap
	Addr     ethcommon.Address
	ChainID  *big.Int
	From     ethcommon.Address
	key      *ecdsa.PrivateKey
	abi      *abi.ABI
}

// NewChain connects to rpcURL and loads the signing key from keyFile (hex, 0x optional).
func NewChain(ctx context.Context, rpcURL string, contract ethcommon.Address, keyFile string) (*Chain, error) {
	raw, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("read eth key: %w", err)
	}
	key, err := ethcrypto.HexToECDSA(strings.TrimPrefix(strings.TrimSpace(string(raw)), "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse eth key: %w", err)
	}
	ec, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial eth rpc: %w", err)
	}
	chainID, err := ec.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	c, err := xmrswap.NewXmrSwap(contract, ec)
	if err != nil {
		return nil, err
	}
	parsed, err := xmrswap.XmrSwapMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	code, err := ec.CodeAt(ctx, contract, nil)
	if err != nil {
		return nil, err
	}
	if len(code) == 0 {
		return nil, fmt.Errorf("no contract code at %s on chain %s", contract.Hex(), chainID)
	}
	return &Chain{
		EC:       ec,
		Contract: c,
		Addr:     contract,
		ChainID:  chainID,
		From:     ethcrypto.PubkeyToAddress(key.PublicKey),
		key:      key,
		abi:      parsed,
	}, nil
}

// Opts returns fresh transaction options bound to ctx.
func (c *Chain) Opts(ctx context.Context) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(c.key, c.ChainID)
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	return opts, nil
}

// Send runs fn to build and submit a transaction, then waits for its receipt and
// checks it succeeded.
func (c *Chain) Send(
	ctx context.Context,
	what string,
	fn func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	opts, err := c.Opts(ctx)
	if err != nil {
		return nil, err
	}
	tx, err := fn(opts)
	if err != nil {
		return nil, fmt.Errorf("%s: send: %w", what, err)
	}
	log.Printf("eth: %s sent, tx %s", what, tx.Hash().Hex())
	receipt, err := block.WaitForReceipt(ctx, c.EC, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("%s: wait for receipt: %w", what, err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("%s: transaction %s reverted", what, tx.Hash().Hex())
	}
	log.Printf("eth: %s confirmed in block %d, gas used %d", what, receipt.BlockNumber.Uint64(), receipt.GasUsed)
	return receipt, nil
}

// Head returns the latest block number and its timestamp.
func (c *Chain) Head(ctx context.Context) (number uint64, ts uint64, err error) {
	h, err := c.EC.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	return h.Number.Uint64(), h.Time, nil
}

// EventID returns the topic hash for a contract event.
func (c *Chain) EventID(name string) ethcommon.Hash {
	return c.abi.Events[name].ID
}

// Logs fetches contract logs in [from, to] with the given topic filter, in chunks.
func (c *Chain) Logs(ctx context.Context, from, to uint64, topics [][]ethcommon.Hash) ([]types.Log, error) {
	var out []types.Log
	for start := from; start <= to; start += logChunk {
		end := start + logChunk - 1
		if end > to {
			end = to
		}
		q := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(start),
			ToBlock:   new(big.Int).SetUint64(end),
			Addresses: []ethcommon.Address{c.Addr},
			Topics:    topics,
		}
		logs, err := c.EC.FilterLogs(ctx, q)
		if err != nil {
			return nil, fmt.Errorf("filter logs %d-%d: %w", start, end, err)
		}
		out = append(out, logs...)
	}
	return out, nil
}

// LogWatcher polls for new logs matching topics, starting after fromBlock, and calls
// handle for each. It returns when ctx ends or handle returns true (done).
func (c *Chain) LogWatcher(
	ctx context.Context,
	fromBlock uint64,
	topics [][]ethcommon.Hash,
	interval time.Duration,
	handle func(types.Log) (done bool, err error),
) error {
	next := fromBlock
	for {
		head, _, err := c.Head(ctx)
		if err != nil {
			log.Printf("eth: head fetch failed (will retry): %s", err)
		} else if head >= next {
			logs, err := c.Logs(ctx, next, head, topics)
			if err != nil {
				log.Printf("eth: log fetch failed (will retry): %s", err)
			} else {
				for _, l := range logs {
					if l.Removed {
						continue
					}
					done, err := handle(l)
					if err != nil {
						return err
					}
					if done {
						return nil
					}
				}
				next = head + 1
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

// ParseWei turns a decimal ETH string into wei.
func ParseWei(eth string) (*big.Int, error) {
	f, ok := new(big.Float).SetPrec(256).SetString(eth)
	if !ok {
		return nil, fmt.Errorf("bad amount %q", eth)
	}
	f.Mul(f, big.NewFloat(1e18))
	wei, _ := f.Int(nil)
	return wei, nil
}

// ParsePiconero turns a decimal XMR string into piconero.
func ParsePiconero(xmr string) (uint64, error) {
	f, ok := new(big.Float).SetPrec(128).SetString(xmr)
	if !ok {
		return 0, fmt.Errorf("bad amount %q", xmr)
	}
	f.Mul(f, big.NewFloat(1e12))
	p, _ := f.Uint64()
	return p, nil
}

// FmtXMR formats piconero as XMR.
func FmtXMR(p uint64) string {
	return fmt.Sprintf("%d.%012d", p/1e12, p%1e12)
}

// FmtETH formats wei as ETH.
func FmtETH(wei *big.Int) string {
	f := new(big.Float).SetInt(wei)
	f.Quo(f, big.NewFloat(1e18))
	return f.Text('f', 6)
}

// NewEthKeyFile generates a fresh key, writes it (hex, mode 0600) and returns the address.
func NewEthKeyFile(path string) (ethcommon.Address, error) {
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return ethcommon.Address{}, err
	}
	hexKey := fmt.Sprintf("%x", ethcrypto.FromECDSA(key))
	if err := os.WriteFile(path, []byte(hexKey+"\n"), 0o600); err != nil {
		return ethcommon.Address{}, err
	}
	return ethcrypto.PubkeyToAddress(key.PublicKey), nil
}
