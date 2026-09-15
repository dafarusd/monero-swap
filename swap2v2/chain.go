// Package swap2v2 implements the maker (Monero seller) side of XmrSwapV2: Monero <-> ETH
// atomic swaps with the offer board on-chain and no peer-to-peer layer.
//
// It is a separate package from swap2 on purpose. swap2 drives the v1 contract, which is still
// deployed and still works; v2 changed enough of the ABI that sharing one code path would mean
// branching on the contract version in every call. The contract-agnostic pieces — the state
// file, the Monero key math, the shared-address handling — are imported from swap2 rather than
// copied.
//
// What a v2 seller must know that a v1 seller did not: posting an offer locks a bond, 5% of the
// offer's ceiling, paid from the same key that pays gas. The bond always comes back, but it comes
// back as a *credit* inside the contract, not as a transfer. Nothing collects it for you, so this
// package withdraws it automatically once there is a balance worth the gas.
package swap2v2

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
	"github.com/athanorlabs/atomic-swap/xmrswap2"
)

// logChunk keeps eth_getLogs queries small enough for public RPC endpoints.
const logChunk = 2000

// Chain wraps an Ethereum client, a signing key, and the XmrSwapV2 contract.
type Chain struct {
	EC       *ethclient.Client
	Contract *xmrswap2.XmrSwapV2
	Addr     ethcommon.Address
	ChainID  *big.Int
	From     ethcommon.Address
	key      *ecdsa.PrivateKey
	abi      *abi.ABI
}

// Terms are the contract's immutable economics, read once at startup. They cannot change for
// the life of the contract, so a seller can be shown them and trust them.
type Terms struct {
	FeeBps     uint16
	BondBps    uint16
	MinTimeout time.Duration
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
	c, err := xmrswap2.NewXmrSwapV2(contract, ec)
	if err != nil {
		return nil, err
	}
	parsed, err := xmrswap2.XmrSwapV2MetaData.GetAbi()
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

// Terms reads the contract's fixed fee, bond rate and minimum timeout.
func (c *Chain) Terms(ctx context.Context) (Terms, error) {
	opts := &bind.CallOpts{Context: ctx}
	fee, err := c.Contract.FeeBps(opts)
	if err != nil {
		return Terms{}, fmt.Errorf("feeBps: %w", err)
	}
	bond, err := c.Contract.BondBps(opts)
	if err != nil {
		return Terms{}, fmt.Errorf("bondBps: %w", err)
	}
	minT, err := c.Contract.MINTIMEOUT(opts)
	if err != nil {
		return Terms{}, fmt.Errorf("MIN_TIMEOUT: %w", err)
	}
	return Terms{FeeBps: fee, BondBps: bond, MinTimeout: time.Duration(minT) * time.Second}, nil
}

// BondFor asks the contract what an offer with this ceiling costs to list.
func (c *Chain) BondFor(ctx context.Context, maxWei *big.Int) (*big.Int, error) {
	return c.Contract.BondFor(&bind.CallOpts{Context: ctx}, maxWei)
}

// Owed is what the contract is holding for us: returned bonds, and fees if we are a recipient.
func (c *Chain) Owed(ctx context.Context) (*big.Int, error) {
	return c.Contract.Withdrawable(&bind.CallOpts{Context: ctx}, c.From)
}

// Balance is the gas wallet's ETH. A v2 seller needs this to cover the bond as well as gas.
func (c *Chain) Balance(ctx context.Context) (*big.Int, error) {
	return c.EC.BalanceAt(ctx, c.From, nil)
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

// Send runs fn to build and submit a transaction, then waits for its receipt and checks it
// succeeded. value, if non-nil, is attached to the transaction (used for the offer bond).
func (c *Chain) Send(
	ctx context.Context,
	what string,
	value *big.Int,
	fn func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Receipt, error) {
	opts, err := c.Opts(ctx)
	if err != nil {
		return nil, err
	}
	if value != nil {
		opts.Value = value
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
