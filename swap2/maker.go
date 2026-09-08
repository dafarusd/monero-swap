package swap2

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/athanorlabs/atomic-swap/coins"
	"github.com/athanorlabs/atomic-swap/common"
	mcrypto "github.com/athanorlabs/atomic-swap/crypto/monero"
	"github.com/athanorlabs/atomic-swap/monero"
)

// Contract stages.
const (
	StageInvalid   = 0
	StagePending   = 1
	StageReady     = 2
	StageCompleted = 3
)

// MakerConfig configures a maker (Monero seller).
type MakerConfig struct {
	Env         common.Environment
	Chain       *Chain
	Wallet      monero.WalletClient // primary Monero wallet; holds inventory unless ManualXMR
	Store       *Store
	Payout      ethcommon.Address // receives ETH / tokens
	Asset       ethcommon.Address // zero address = ETH
	MinWei      *big.Int
	MaxWei      *big.Int
	XmrPerAsset *big.Int // piconero per 1e18 base units of asset
	OfferTTL    time.Duration
	Timeout1    time.Duration
	Timeout2    time.Duration
	ManualXMR   bool // print the shared address and wait for the human to pay
	Poll        time.Duration
}

// Maker keeps one live offer on the contract and runs every swap that takes it.
type Maker struct {
	cfg     MakerConfig
	wg      sync.WaitGroup
	mu      sync.Mutex
	running map[string]bool
}

// NewMaker returns a maker; call Run.
func NewMaker(cfg MakerConfig) *Maker {
	if cfg.Poll == 0 {
		cfg.Poll = 20 * time.Second
	}
	return &Maker{cfg: cfg, running: map[string]bool{}}
}

// Run blocks until ctx ends.
func (m *Maker) Run(ctx context.Context) error {
	log.Printf("maker: eth %s, payout %s, contract %s, chain %s", m.cfg.Chain.From.Hex(), m.cfg.Payout.Hex(), m.cfg.Chain.Addr.Hex(), m.cfg.Chain.ChainID)
	log.Printf("maker: xmr wallet %s", m.cfg.Wallet.PrimaryAddress())

	// resume unfinished swaps
	var resume []string
	m.cfg.Store.View(func(st *State) {
		for id, s := range st.Swaps {
			if s.Role == "maker" && s.State != StateClaimed && s.State != StateRecovered {
				resume = append(resume, id)
			}
		}
	})
	for _, id := range resume {
		log.Printf("maker: resuming swap %s", id)
		m.startSwap(ctx, id)
	}

	for {
		if err := m.reconcileOffers(ctx); err != nil {
			log.Printf("maker: offer reconcile failed (will retry): %s", err)
		}
		select {
		case <-ctx.Done():
			m.wg.Wait()
			return ctx.Err()
		case <-time.After(m.cfg.Poll):
		}
	}
}

// reconcileOffers detects takes of our offers, retires expired ones, and keeps one live.
func (m *Maker) reconcileOffers(ctx context.Context) error {
	var active []*OfferRec
	m.cfg.Store.View(func(st *State) {
		for _, o := range st.Offers {
			if o.Active {
				cp := *o
				active = append(active, &cp)
			}
		}
	})

	head, now, err := m.cfg.Chain.Head(ctx)
	if err != nil {
		return err
	}

	live := 0
	for _, o := range active {
		id := hashFromHex(o.OfferID)
		onchain, err := m.cfg.Chain.Contract.Offers(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			return fmt.Errorf("read offer %s: %w", o.OfferID, err)
		}
		if onchain.Active {
			if onchain.Expiry <= now {
				log.Printf("maker: offer %s expired, cancelling", short(o.OfferID))
				if _, err := m.cfg.Chain.Send(ctx, "cancelOffer", func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return m.cfg.Chain.Contract.CancelOffer(opts, id)
				}); err != nil {
					return err
				}
				_ = m.cfg.Store.Update(func(st *State) { st.Offers[o.OfferID].Active = false })
				continue
			}
			live++
			continue
		}
		// not active on chain: taken or cancelled. Look for the take.
		logs, err := m.cfg.Chain.Logs(ctx, o.PostedBlock, head,
			[][]ethcommon.Hash{{m.cfg.Chain.EventID("SwapCreated")}, {id}})
		if err != nil {
			return err
		}
		if len(logs) == 0 {
			log.Printf("maker: offer %s inactive with no take (cancelled?)", short(o.OfferID))
			_ = m.cfg.Store.Update(func(st *State) { st.Offers[o.OfferID].Active = false })
			continue
		}
		ev, err := m.cfg.Chain.Contract.ParseSwapCreated(logs[0])
		if err != nil {
			return err
		}
		if err := m.recordTake(ctx, o, ev); err != nil {
			return err
		}
	}

	if live == 0 {
		return m.postOffer(ctx)
	}
	return nil
}

// recordTake stores the new swap and starts working it.
func (m *Maker) recordTake(ctx context.Context, o *OfferRec, ev *xmrSwapCreated) error {
	swapID := hex.EncodeToString(ev.SwapId[:])
	log.Printf("maker: offer %s taken by %s for %s (%s XMR), swap %s",
		short(o.OfferID), ev.Taker.Hex(), fmtAsset(ev.Value, m.cfg.Asset), FmtXMR(ev.XmrPiconero.Uint64()), short(swapID))

	height, err := m.cfg.Wallet.GetHeight()
	if err != nil {
		return fmt.Errorf("xmr height: %w", err)
	}
	restore := height
	if restore > 20 {
		restore -= 20
	}

	myPub, err := ParsePub(o.SpendPub)
	if err != nil {
		return err
	}
	myView, err := ParseViewPriv(o.ViewPriv)
	if err != nil {
		return err
	}
	theirPub, err := mcrypto.NewPublicKeyFromBytes(ev.TakerSpendPub[:])
	if err != nil {
		return fmt.Errorf("taker spend pub: %w", err)
	}
	theirView, err := mcrypto.NewPrivateViewKeyFromBytes(ev.TakerViewPriv[:])
	if err != nil {
		return fmt.Errorf("taker view key: %w", err)
	}
	addr, _ := SharedAddress(m.cfg.Env, myPub, theirPub, myView, theirView)

	err = m.cfg.Store.Update(func(st *State) {
		st.Offers[o.OfferID].Active = false
		st.Offers[o.OfferID].SwapID = swapID
		st.Swaps[swapID] = &SwapRec{
			SwapID:        swapID,
			OfferID:       o.OfferID,
			Role:          "maker",
			State:         StateXMRLockPending,
			Asset:         o.Asset,
			ValueWei:      ev.Value.String(),
			XmrPiconero:   ev.XmrPiconero.Uint64(),
			MySpendPriv:   o.SpendPriv,
			MyViewPriv:    o.ViewPriv,
			TheirSpendPub: hex.EncodeToString(ev.TakerSpendPub[:]),
			TheirViewPriv: hex.EncodeToString(ev.TakerViewPriv[:]),
			SharedAddress: addr.String(),
			RestoreHeight: restore,
			Timeout1:      ev.Timeout1,
			Timeout2:      ev.Timeout2,
			CreatedBlock:  ev.Raw.BlockNumber,
			EthTxHash:     ev.Raw.TxHash.Hex(),
		}
	})
	if err != nil {
		return err
	}
	m.startSwap(ctx, swapID)
	return nil
}

// postOffer generates fresh one-time keys and posts a single-use offer.
func (m *Maker) postOffer(ctx context.Context) error {
	if !m.cfg.ManualXMR {
		bal, err := m.cfg.Wallet.GetBalance(0)
		if err != nil {
			return fmt.Errorf("xmr balance: %w", err)
		}
		need := new(big.Int).Mul(m.cfg.MaxWei, m.cfg.XmrPerAsset)
		need.Div(need, big.NewInt(1e18))
		if new(big.Int).SetUint64(bal.UnlockedBalance).Cmp(need) < 0 {
			log.Printf("maker: not posting: unlocked XMR %s < %s needed for max offer",
				FmtXMR(bal.UnlockedBalance), FmtXMR(need.Uint64()))
			return nil
		}
	}

	kp, err := mcrypto.GenerateKeys()
	if err != nil {
		return err
	}
	spendPriv, viewPriv, spendPub := HexKeys(kp)
	pending := "pending:" + spendPub

	// keys hit disk before the transaction leaves this machine
	if err := m.cfg.Store.Update(func(st *State) {
		st.Offers[pending] = &OfferRec{OfferID: pending, SpendPriv: spendPriv, ViewPriv: viewPriv, SpendPub: spendPub, Asset: m.cfg.Asset.Hex()}
	}); err != nil {
		return err
	}

	expiry := uint64(time.Now().Add(m.cfg.OfferTTL).Unix())
	receipt, err := m.cfg.Chain.Send(ctx, "postOffer", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.cfg.Chain.Contract.PostOffer(opts, m.cfg.Asset, m.cfg.MinWei, m.cfg.MaxWei, m.cfg.XmrPerAsset,
			expiry, uint64(m.cfg.Timeout1.Seconds()), uint64(m.cfg.Timeout2.Seconds()),
			Bytes32(kp.PublicKeyPair().SpendKey().Bytes()), Bytes32(kp.ViewKey().Bytes()), m.cfg.Payout)
	})
	if err != nil {
		return err
	}

	var offerID string
	for _, l := range receipt.Logs {
		ev, err := m.cfg.Chain.Contract.ParseOfferPosted(*l)
		if err == nil {
			offerID = hex.EncodeToString(ev.OfferId[:])
			break
		}
	}
	if offerID == "" {
		return fmt.Errorf("postOffer receipt has no OfferPosted event")
	}
	log.Printf("maker: posted offer %s: %s-%s %s at %s XMR per unit, expires %s",
		short(offerID), fmtAsset(m.cfg.MinWei, m.cfg.Asset), fmtAsset(m.cfg.MaxWei, m.cfg.Asset), assetName(m.cfg.Asset),
		FmtXMR(m.cfg.XmrPerAsset.Uint64()), time.Unix(int64(expiry), 0).Format(time.RFC3339))

	return m.cfg.Store.Update(func(st *State) {
		rec := st.Offers[pending]
		delete(st.Offers, pending)
		rec.OfferID = offerID
		rec.PostedBlock = receipt.BlockNumber.Uint64()
		rec.TxHash = receipt.TxHash.Hex()
		rec.Active = true
		st.Offers[offerID] = rec
	})
}

func (m *Maker) startSwap(ctx context.Context, swapID string) {
	m.mu.Lock()
	if m.running[swapID] {
		m.mu.Unlock()
		return
	}
	m.running[swapID] = true
	m.mu.Unlock()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer func() {
			m.mu.Lock()
			delete(m.running, swapID)
			m.mu.Unlock()
		}()
		if err := m.runSwap(ctx, swapID); err != nil && ctx.Err() == nil {
			log.Printf("maker: swap %s stopped with error: %s", short(swapID), err)
			_ = m.cfg.Store.Update(func(st *State) { st.Swaps[swapID].LastError = err.Error() })
		}
	}()
}

func (m *Maker) swap(swapID string) *SwapRec {
	var cp SwapRec
	m.cfg.Store.View(func(st *State) { cp = *st.Swaps[swapID] })
	return &cp
}

func (m *Maker) setState(swapID, state string, mut func(*SwapRec)) error {
	log.Printf("maker: swap %s -> %s", short(swapID), state)
	return m.cfg.Store.Update(func(st *State) {
		s := st.Swaps[swapID]
		s.State = state
		s.LastError = ""
		if mut != nil {
			mut(s)
		}
	})
}

// runSwap drives one swap through its states until it is claimed or recovered.
func (m *Maker) runSwap(ctx context.Context, swapID string) error {
	for {
		s := m.swap(swapID)
		switch s.State {
		case StateXMRLockPending:
			if err := m.lockXMR(ctx, s); err != nil {
				return err
			}
		case StateXMRLocked:
			if err := m.awaitClaim(ctx, s); err != nil {
				return err
			}
		case StateRecovering:
			if err := m.recover(ctx, s); err != nil {
				return err
			}
		case StateClaimed, StateRecovered:
			return nil
		default:
			return fmt.Errorf("unknown state %q", s.State)
		}
	}
}

func (m *Maker) lockXMR(ctx context.Context, s *SwapRec) error {
	id := hashFromHex(s.SwapID)
	onchain, err := m.cfg.Chain.Contract.Swaps(&bind.CallOpts{Context: ctx}, id)
	if err != nil {
		return err
	}
	_, now, err := m.cfg.Chain.Head(ctx)
	if err != nil {
		return err
	}
	if onchain.Stage != StagePending || now+uint64((30*time.Minute).Seconds()) > s.Timeout1 {
		// Too late to lock safely: the taker needs ~20 min of confirmations before t1.
		log.Printf("maker: swap %s: not locking XMR (stage=%d, %ds to t1); waiting for taker refund", short(s.SwapID), onchain.Stage, int64(s.Timeout1)-int64(now))
		return m.setState(s.SwapID, StateXMRLocked, nil) // awaitClaim handles refund/recovery
	}

	addr, err := mcrypto.NewAddress(s.SharedAddress, m.cfg.Env)
	if err != nil {
		return err
	}

	if m.cfg.ManualXMR {
		fmt.Printf("\n=== ACTION NEEDED ===\nSend exactly %s XMR to:\n%s\nfrom your own wallet, before %s.\n=====================\n\n",
			FmtXMR(s.XmrPiconero), addr, time.Unix(int64(s.Timeout1), 0).Add(-25*time.Minute).Format(time.RFC3339))
		myView, err := ParseViewPriv(s.MyViewPriv)
		if err != nil {
			return err
		}
		theirView, err := ParseViewPriv(s.TheirViewPriv)
		if err != nil {
			return err
		}
		viewAB := mcrypto.SumPrivateViewKeys(myView, theirView)
		ok, err := WaitForXMR(ctx, m.cfg.Wallet, "maker-watch-"+short(s.SwapID), viewAB, addr, s.RestoreHeight,
			s.XmrPiconero, time.Unix(int64(s.Timeout1), 0).Add(-5*time.Minute))
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("XMR was not paid in time; taker will refund, nothing lost")
		}
		return m.setState(s.SwapID, StateXMRLocked, func(r *SwapRec) { r.XmrLockTxID = "manual" })
	}

	log.Printf("maker: swap %s: sending %s XMR to shared address %s", short(s.SwapID), FmtXMR(s.XmrPiconero), addr)
	tr, err := m.cfg.Wallet.Transfer(ctx, addr, 0, coins.NewPiconeroAmount(s.XmrPiconero), monero.MinSpendConfirmations)
	if err != nil {
		return fmt.Errorf("xmr transfer: %w", err)
	}
	return m.setState(s.SwapID, StateXMRLocked, func(r *SwapRec) { r.XmrLockTxID = tr.TxID })
}

// awaitClaim claims once the taker sets ready (or t1 passes), or recovers if the taker refunds.
func (m *Maker) awaitClaim(ctx context.Context, s *SwapRec) error {
	id := hashFromHex(s.SwapID)
	for {
		onchain, err := m.cfg.Chain.Contract.Swaps(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			log.Printf("maker: read swap failed (will retry): %s", err)
		} else {
			head, now, err := m.cfg.Chain.Head(ctx)
			if err == nil {
				if onchain.Stage == StageCompleted {
					// Either we claimed (would already be recorded) or the taker refunded.
					secret, err := m.findSecret(ctx, "Refunded", s, head)
					if err != nil {
						return err
					}
					if secret != "" {
						return m.setState(s.SwapID, StateRecovering, func(r *SwapRec) { r.TheirSecret = secret })
					}
					return fmt.Errorf("swap completed on-chain but no Refunded/Claimed event found")
				}
				if now < s.Timeout2 && (onchain.Stage == StageReady || now >= s.Timeout1) {
					mySecret, err := ParseSpendPriv(s.MySpendPriv)
					if err != nil {
						return err
					}
					receipt, err := m.cfg.Chain.Send(ctx, "claim", func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return m.cfg.Chain.Contract.Claim(opts, id, Bytes32(mySecret.Bytes()))
					})
					if err != nil {
						log.Printf("maker: claim failed (will re-check): %s", err)
					} else {
						return m.setState(s.SwapID, StateClaimed, func(r *SwapRec) { r.EthTxHash = receipt.TxHash.Hex() })
					}
				} else if now >= s.Timeout2 {
					log.Printf("maker: swap %s past t2 without claim; waiting for taker refund to recover XMR", short(s.SwapID))
				}
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.cfg.Poll):
		}
	}
}

// findSecret returns the secret revealed by a Refunded (or Claimed) event for this swap.
func (m *Maker) findSecret(ctx context.Context, event string, s *SwapRec, head uint64) (string, error) {
	logs, err := m.cfg.Chain.Logs(ctx, s.CreatedBlock, head,
		[][]ethcommon.Hash{{m.cfg.Chain.EventID(event)}, {hashFromHex(s.SwapID)}})
	if err != nil {
		return "", err
	}
	if len(logs) == 0 {
		return "", nil
	}
	switch event {
	case "Refunded":
		ev, err := m.cfg.Chain.Contract.ParseRefunded(logs[0])
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(ev.Secret[:]), nil
	default:
		ev, err := m.cfg.Chain.Contract.ParseClaimed(logs[0])
		if err != nil {
			return "", err
		}
		return hex.EncodeToString(ev.Secret[:]), nil
	}
}

// recover sweeps our XMR back after the taker refunded (their secret is public now).
func (m *Maker) recover(ctx context.Context, s *SwapRec) error {
	if s.XmrLockTxID == "" {
		return m.setState(s.SwapID, StateRecovered, nil) // we never locked
	}
	mySpend, err := ParseSpendPriv(s.MySpendPriv)
	if err != nil {
		return err
	}
	theirSpend, err := ParseSpendPriv(s.TheirSecret)
	if err != nil {
		return err
	}
	myView, err := ParseViewPriv(s.MyViewPriv)
	if err != nil {
		return err
	}
	theirView, err := ParseViewPriv(s.TheirViewPriv)
	if err != nil {
		return err
	}
	spendAB := mcrypto.SumPrivateSpendKeys(mySpend, theirSpend)
	viewAB := mcrypto.SumPrivateViewKeys(myView, theirView)
	log.Printf("maker: swap %s: taker refunded; sweeping XMR back to %s", short(s.SwapID), m.cfg.Wallet.PrimaryAddress())
	ids, err := SweepShared(ctx, m.cfg.Wallet, "maker-recover-"+short(s.SwapID), spendAB, viewAB, s.RestoreHeight, m.cfg.Wallet.PrimaryAddress())
	if err != nil {
		return err
	}
	return m.setState(s.SwapID, StateRecovered, func(r *SwapRec) { r.SweepTxIDs = ids })
}

// ---- small helpers shared with the taker

type xmrSwapCreated = xmrSwapCreatedEvent

func hashFromHex(h string) (out [32]byte) {
	b, _ := hex.DecodeString(strings.TrimPrefix(h, "0x"))
	copy(out[:], b)
	return out
}

func short(id string) string {
	if len(id) > 10 {
		return id[:10]
	}
	return id
}

func assetName(a ethcommon.Address) string {
	if a == (ethcommon.Address{}) {
		return "ETH"
	}
	return a.Hex()
}

func fmtAsset(v *big.Int, a ethcommon.Address) string {
	if a == (ethcommon.Address{}) {
		return FmtETH(v)
	}
	return v.String()
}
