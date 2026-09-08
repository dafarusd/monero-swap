package swap2

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/athanorlabs/atomic-swap/common"
	mcrypto "github.com/athanorlabs/atomic-swap/crypto/monero"
	"github.com/athanorlabs/atomic-swap/monero"
)

// TakerConfig configures a taker (pays ETH, receives Monero).
type TakerConfig struct {
	Env       common.Environment
	Chain     *Chain
	Wallet    monero.WalletClient // used for one-off wallets and as default payout
	Store     *Store
	PayoutXMR *mcrypto.Address // nil = Wallet's primary address
	Poll      time.Duration
}

// Taker takes offers and sees them through.
type Taker struct {
	cfg TakerConfig
}

// NewTaker returns a taker.
func NewTaker(cfg TakerConfig) *Taker {
	if cfg.Poll == 0 {
		cfg.Poll = 20 * time.Second
	}
	if cfg.PayoutXMR == nil {
		cfg.PayoutXMR = cfg.Wallet.PrimaryAddress()
	}
	return &Taker{cfg: cfg}
}

// Take locks `amountWei` of ETH against an offer and runs the swap to completion.
func (t *Taker) Take(ctx context.Context, offerID [32]byte, amountWei *big.Int) (string, error) {
	o, err := t.cfg.Chain.Contract.Offers(&bind.CallOpts{Context: ctx}, offerID)
	if err != nil {
		return "", fmt.Errorf("read offer: %w", err)
	}
	_, now, err := t.cfg.Chain.Head(ctx)
	if err != nil {
		return "", err
	}
	switch {
	case !o.Active:
		return "", fmt.Errorf("offer is not active")
	case o.Expiry <= now:
		return "", fmt.Errorf("offer expired")
	case o.Asset != (ethcommon.Address{}):
		return "", fmt.Errorf("this command handles ETH offers only (offer asset %s)", o.Asset.Hex())
	case amountWei.Cmp(o.MinAmount) < 0 || amountWei.Cmp(o.MaxAmount) > 0:
		return "", fmt.Errorf("amount %s ETH outside offer range %s-%s", FmtETH(amountWei), FmtETH(o.MinAmount), FmtETH(o.MaxAmount))
	}
	xmr := new(big.Int).Mul(amountWei, o.XmrPerAsset)
	xmr.Div(xmr, big.NewInt(1e18))
	log.Printf("taker: taking offer %s: pay %s ETH, receive %s XMR", short(hex.EncodeToString(offerID[:])), FmtETH(amountWei), FmtXMR(xmr.Uint64()))

	kp, err := mcrypto.GenerateKeys()
	if err != nil {
		return "", err
	}
	spendPriv, viewPriv, spendPub := HexKeys(kp)
	pending := "pending:" + spendPub
	if err := t.cfg.Store.Update(func(st *State) {
		st.Swaps[pending] = &SwapRec{SwapID: pending, Role: "taker", State: "pending_take", OfferID: hex.EncodeToString(offerID[:]),
			MySpendPriv: spendPriv, MyViewPriv: viewPriv}
	}); err != nil {
		return "", err
	}

	height, err := t.cfg.Wallet.GetHeight()
	if err != nil {
		return "", fmt.Errorf("xmr height: %w", err)
	}
	restore := height
	if restore > 20 {
		restore -= 20
	}

	receipt, err := t.cfg.Chain.Send(ctx, "takeOffer", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.Value = amountWei
		return t.cfg.Chain.Contract.TakeOffer(opts, offerID, amountWei,
			Bytes32(kp.PublicKeyPair().SpendKey().Bytes()), Bytes32(kp.ViewKey().Bytes()))
	})
	if err != nil {
		return "", err
	}
	var ev *xmrSwapCreatedEvent
	for _, l := range receipt.Logs {
		if e, err := t.cfg.Chain.Contract.ParseSwapCreated(*l); err == nil {
			ev = e
			break
		}
	}
	if ev == nil {
		return "", fmt.Errorf("takeOffer receipt has no SwapCreated event")
	}
	swapID := hex.EncodeToString(ev.SwapId[:])

	makerPub, err := mcrypto.NewPublicKeyFromBytes(o.MakerSpendPub[:])
	if err != nil {
		return "", err
	}
	makerView, err := mcrypto.NewPrivateViewKeyFromBytes(o.MakerViewPriv[:])
	if err != nil {
		return "", err
	}
	addr, _ := SharedAddress(t.cfg.Env, makerPub, kp.PublicKeyPair().SpendKey(), makerView, kp.ViewKey())

	if err := t.cfg.Store.Update(func(st *State) {
		rec := st.Swaps[pending]
		delete(st.Swaps, pending)
		rec.SwapID = swapID
		rec.State = StateETHLocked
		rec.Asset = ethcommon.Address{}.Hex()
		rec.ValueWei = amountWei.String()
		rec.XmrPiconero = ev.XmrPiconero.Uint64()
		rec.TheirSpendPub = hex.EncodeToString(o.MakerSpendPub[:])
		rec.TheirViewPriv = hex.EncodeToString(o.MakerViewPriv[:])
		rec.SharedAddress = addr.String()
		rec.RestoreHeight = restore
		rec.Timeout1 = ev.Timeout1
		rec.Timeout2 = ev.Timeout2
		rec.CreatedBlock = receipt.BlockNumber.Uint64()
		rec.EthTxHash = receipt.TxHash.Hex()
		st.Swaps[swapID] = rec
	}); err != nil {
		return "", err
	}
	log.Printf("taker: swap %s created; ETH locked. Shared Monero address %s. t1=%s t2=%s",
		short(swapID), addr, time.Unix(int64(ev.Timeout1), 0).Format(time.Kitchen), time.Unix(int64(ev.Timeout2), 0).Format(time.Kitchen))
	return swapID, t.Resume(ctx, swapID)
}

func (t *Taker) swap(swapID string) *SwapRec {
	var cp SwapRec
	t.cfg.Store.View(func(st *State) { cp = *st.Swaps[swapID] })
	return &cp
}

func (t *Taker) setState(swapID, state string, mut func(*SwapRec)) error {
	log.Printf("taker: swap %s -> %s", short(swapID), state)
	return t.cfg.Store.Update(func(st *State) {
		s := st.Swaps[swapID]
		s.State = state
		s.LastError = ""
		if mut != nil {
			mut(s)
		}
	})
}

// Resume continues a swap from its stored state until it is swept or refunded.
func (t *Taker) Resume(ctx context.Context, swapID string) error {
	for {
		s := t.swap(swapID)
		var err error
		switch s.State {
		case StateETHLocked:
			err = t.awaitXMR(ctx, s)
		case StateXMRSeen:
			err = t.setReady(ctx, s)
		case StateReady:
			err = t.awaitClaim(ctx, s)
		case StateSwept, StateRefunded:
			return nil
		default:
			err = fmt.Errorf("unknown state %q", s.State)
		}
		if err != nil {
			_ = t.cfg.Store.Update(func(st *State) { st.Swaps[swapID].LastError = err.Error() })
			return err
		}
	}
}

func (t *Taker) keys(s *SwapRec) (mySpend *mcrypto.PrivateSpendKey, viewAB *mcrypto.PrivateViewKey, addr *mcrypto.Address, err error) {
	mySpend, err = ParseSpendPriv(s.MySpendPriv)
	if err != nil {
		return
	}
	myView, err := ParseViewPriv(s.MyViewPriv)
	if err != nil {
		return
	}
	theirView, err := ParseViewPriv(s.TheirViewPriv)
	if err != nil {
		return
	}
	viewAB = mcrypto.SumPrivateViewKeys(myView, theirView)
	addr, err = mcrypto.NewAddress(s.SharedAddress, t.cfg.Env)
	return
}

func (t *Taker) awaitXMR(ctx context.Context, s *SwapRec) error {
	_, viewAB, addr, err := t.keys(s)
	if err != nil {
		return err
	}
	deadline := time.Unix(int64(s.Timeout1), 0).Add(-10 * time.Minute)
	log.Printf("taker: swap %s: waiting for %s XMR at %s (until %s)", short(s.SwapID), FmtXMR(s.XmrPiconero), addr, deadline.Format(time.Kitchen))
	ok, err := WaitForXMR(ctx, t.cfg.Wallet, "taker-watch-"+short(s.SwapID), viewAB, addr, s.RestoreHeight, s.XmrPiconero, deadline)
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("taker: swap %s: XMR did not arrive in time; refunding ETH", short(s.SwapID))
		return t.refund(ctx, s)
	}
	return t.setState(s.SwapID, StateXMRSeen, nil)
}

func (t *Taker) setReady(ctx context.Context, s *SwapRec) error {
	id := hashFromHex(s.SwapID)
	if _, err := t.cfg.Chain.Send(ctx, "setReady", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return t.cfg.Chain.Contract.SetReady(opts, id)
	}); err != nil {
		return err
	}
	return t.setState(s.SwapID, StateReady, nil)
}

// awaitClaim waits for the maker to claim (revealing its secret), then sweeps the Monero.
func (t *Taker) awaitClaim(ctx context.Context, s *SwapRec) error {
	id := hashFromHex(s.SwapID)
	for {
		head, now, err := t.cfg.Chain.Head(ctx)
		if err == nil {
			logs, err := t.cfg.Chain.Logs(ctx, s.CreatedBlock, head, [][]ethcommon.Hash{{t.cfg.Chain.EventID("Claimed")}, {id}})
			if err != nil {
				log.Printf("taker: log fetch failed (will retry): %s", err)
			} else if len(logs) > 0 {
				ev, err := t.cfg.Chain.Contract.ParseClaimed(logs[0])
				if err != nil {
					return err
				}
				secret := hex.EncodeToString(ev.Secret[:])
				log.Printf("taker: swap %s: maker claimed ETH in tx %s and revealed its secret", short(s.SwapID), ev.Raw.TxHash.Hex())
				_ = t.cfg.Store.Update(func(st *State) { st.Swaps[s.SwapID].TheirSecret = secret })
				return t.sweep(ctx, t.swap(s.SwapID))
			}
			if now >= s.Timeout2 {
				log.Printf("taker: swap %s: maker never claimed before t2; refunding ETH", short(s.SwapID))
				return t.refund(ctx, s)
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(t.cfg.Poll):
		}
	}
}

func (t *Taker) sweep(ctx context.Context, s *SwapRec) error {
	mySpend, viewAB, _, err := t.keys(s)
	if err != nil {
		return err
	}
	theirSpend, err := ParseSpendPriv(s.TheirSecret)
	if err != nil {
		return err
	}
	spendAB := mcrypto.SumPrivateSpendKeys(mySpend, theirSpend)
	log.Printf("taker: swap %s: sweeping XMR to %s", short(s.SwapID), t.cfg.PayoutXMR)
	ids, err := SweepShared(ctx, t.cfg.Wallet, "taker-sweep-"+short(s.SwapID), spendAB, viewAB, s.RestoreHeight, t.cfg.PayoutXMR)
	if err != nil {
		return err
	}
	log.Printf("taker: swap %s: XMR sweep sent: %v", short(s.SwapID), ids)
	return t.setState(s.SwapID, StateSwept, func(r *SwapRec) { r.SweepTxIDs = ids })
}

// refund takes the ETH back, revealing our secret so the maker can recover any XMR it locked.
func (t *Taker) refund(ctx context.Context, s *SwapRec) error {
	mySpend, _, _, err := t.keys(s)
	if err != nil {
		return err
	}
	id := hashFromHex(s.SwapID)
	for {
		receipt, err := t.cfg.Chain.Send(ctx, "refund", func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return t.cfg.Chain.Contract.Refund(opts, id, Bytes32(mySpend.Bytes()))
		})
		if err == nil {
			return t.setState(s.SwapID, StateRefunded, func(r *SwapRec) { r.EthTxHash = receipt.TxHash.Hex() })
		}
		// The contract may refuse until t2 (e.g. ready was set). Keep trying on the poll interval.
		log.Printf("taker: refund not accepted yet (%s); retrying", err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(t.cfg.Poll):
		}
	}
}
