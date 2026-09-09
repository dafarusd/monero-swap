package swap2

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	contracts "github.com/athanorlabs/atomic-swap/ethereum"
)

// Asset is ETH or an ERC-20 token on the swap chain.
type Asset struct {
	Addr     ethcommon.Address // zero for ETH
	Symbol   string
	Decimals uint8
}

// ETH is the native coin.
var ETH = Asset{Symbol: "ETH", Decimals: 18}

func (a Asset) IsETH() bool { return a.Addr == (ethcommon.Address{}) }

// LoadAsset reads symbol and decimals from the chain (or returns ETH for the zero address).
func LoadAsset(ctx context.Context, ec *ethclient.Client, addr ethcommon.Address) (Asset, error) {
	if addr == (ethcommon.Address{}) {
		return ETH, nil
	}
	c, err := contracts.NewIERC20Caller(addr, ec)
	if err != nil {
		return Asset{}, err
	}
	opts := &bind.CallOpts{Context: ctx}
	dec, err := c.Decimals(opts)
	if err != nil {
		return Asset{}, fmt.Errorf("token %s: decimals(): %w (is this an ERC-20?)", addr.Hex(), err)
	}
	sym, err := c.Symbol(opts)
	if err != nil {
		sym = addr.Hex()[:10]
	}
	return Asset{Addr: addr, Symbol: sym, Decimals: dec}, nil
}

func (a Asset) unit() *big.Int { return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(a.Decimals)), nil) }

// Parse turns "12.5" into base units, exactly.
func (a Asset) Parse(s string) (*big.Int, error) { return ParseDecimal(s, uint(a.Decimals)) }

// Fmt prints base units as a decimal with the symbol.
func (a Asset) Fmt(v *big.Int) string {
	f := new(big.Float).SetPrec(256).SetInt(v)
	f.Quo(f, new(big.Float).SetInt(a.unit()))
	return f.Text('f', int(min(a.Decimals, 6))) + " " + a.Symbol
}

// XmrPerAsset converts "piconero per whole unit" into the contract's "piconero per 1e18 base units".
func (a Asset) XmrPerAsset(piconeroPerWhole uint64) *big.Int {
	v := new(big.Int).SetUint64(piconeroPerWhole)
	v.Mul(v, big.NewInt(1e18))
	return v.Div(v, a.unit())
}

func min(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}
