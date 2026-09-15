package swap2v2

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
	"github.com/athanorlabs/atomic-swap/swap2"
	"github.com/athanorlabs/atomic-swap/xmrswap2"
)

// Contract stages.
const (
	StageInvalid   = 0
	StagePending   = 1
	StageReady     = 2
	StageCompleted = 3
)

// lockGrace is how much of timeout1 must remain before it is still safe to send Monero: the
// buyer needs ~20 minutes of confirmations before t1, or they will refund and we lose the fee
// on a swap we already paid for.
const lockGrace = 40 * time.Minute

// MakerConfig configures a maker (Monero seller) against XmrSwapV2.
type MakerConfig struct {
	Env    common.Environment
	Chain  *Chain
	Wallet monero.WalletClient // primary Monero wallet; holds inventory unless ManualXMR
	Store  *swap2.Store
	Payout ethcommon.Address // receives the ETH; never the gas key

	MinWei      *big.Int
	MaxWei      *big.Int
	XmrPerAsset *big.Int // piconero per 1e18 wei
	OfferTTL    time.Duration
	Timeout1    time.Duration
	Timeout2    time.Duration
	ManualXMR   bool // print the shared address and wait for the human to pay
	Poll        time.Duration

	// MinWithdraw is the smallest credited balance worth spending gas to collect.
	// Bonds and fees accrue inside the contract; nothing pushes them out.
	MinWithdraw *big.Int
}

// Maker keeps one live offer on the contract and runs every swap that takes it.
type Maker struct {
	cfg     MakerConfig
	terms   Terms
	bond    *big.Int
	wg      sync.WaitGroup
	mu      sync.Mutex
	running map[string]bool
}

// NewMaker returns a maker; call Run.
func NewMaker(cfg MakerConfig) *Maker {
	if cfg.Poll == 0 {
		cfg.Poll = 20 * time.Second
	}
	if cfg.MinWithdraw == nil {
		cfg.MinWithdraw = big.NewInt(1e14) // 0.0001 ETH; below this the gas is not worth it
	}
	return &Maker{cfg: cfg, running: map[string]bool{}}
}

// Run blocks until ctx ends.
func (m *Maker) Run(ctx context.Context) error {
	terms, err := m.cfg.Chain.Terms(ctx)
	if err != nil {
		return fmt.Errorf("read contract terms: %w", err)
	}
	m.terms = terms

	if m.cfg.Timeout1 < terms.MinTimeout || m.cfg.Timeout2 < terms.MinTimeout {
		return fmt.Errorf("this contract requires both timeouts to be at least %s; got t1=%s t2=%s",
			terms.MinTimeout, m.cfg.Timeout1, m.cfg.Timeout2)
	}

	bond, err := m.cfg.Chain.BondFor(ctx, m.cfg.MaxWei)
	if err != nil {
		return fmt.Errorf("read bond: %w", err)
	}
	m.bond = bond

	log.Printf("maker: eth %s, payout %s, contract %s, chain %s",
		m.cfg.Chain.From.Hex(), m.cfg.Payout.Hex(), m.cfg.Chain.Addr.Hex(), m.cfg.Chain.ChainID)
	log.Printf("maker: xmr wallet %s", m.cfg.Wallet.PrimaryAddress())
	log.Printf("maker: contract terms: fee %d bps, bond %d bps, minimum timeout %s",
		terms.FeeBps, terms.BondBps, terms.MinTimeout)
	log.Printf("maker: each offer locks a bond of %s ETH from the gas wallet. It is always returned, "+
		"but it is returned as a credit inside the contract — this program collects it for you.",
		swap2.FmtETH(bond))

	// resume unfinished swaps
	var resume []string
	m.cfg.Store.View(func(st *swap2.State) {
		for id, s := range st.Swaps {
			if s.Role == "maker" && s.State != swap2.StateClaimed && s.State != swap2.StateRecovered {
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
		if err := m.collect(ctx); err != nil {
			log.Printf("maker: withdraw failed (will retry): %s", err)
		}
		select {
		case <-ctx.Done():
			m.wg.Wait()
			return ctx.Err()
		case <-time.After(m.cfg.Poll):
		}
	}
}

// collect withdraws anything the contract is holding for us — returned bonds, and fees if this
// address is a fee recipient. v2 credits rather than transfers, so without this the bonds just
// sit there.
func (m *Maker) collect(ctx context.Context) error {
	owed, err := m.cfg.Chain.Owed(ctx)
	if err != nil {
		return err
	}
	if owed.Sign() == 0 || owed.Cmp(m.cfg.MinWithdraw) < 0 {
		return nil
	}
	log.Printf("maker: collecting %s ETH credited by the contract", swap2.FmtETH(owed))
	_, err = m.cfg.Chain.Send(ctx, "withdraw", nil, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.cfg.Chain.Contract.Withdraw(opts)
	})
	return err
}

// reconcileOffers detects takes of our offers, retires expired ones, and keeps one live.
func (m *Maker) reconcileOffers(ctx context.Context) error {
	var active []*swap2.OfferRec
	m.cfg.Store.View(func(st *swap2.State) {
		for _, o := range st.Offers {
			if o.TxHash != "" && o.SwapID == "" && (o.Active || !o.TakeChecked) {
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
		onchain, err := m.cfg.Chain.Contract.GetOffer(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			return fmt.Errorf("read offer %s: %w", o.OfferID, err)
		}
		if onchain.Active {
			if onchain.Expiry <= now {
				log.Printf("maker: offer %s expired, cancelling (bond comes back as a credit)", short(o.OfferID))
				if _, err := m.cfg.Chain.Send(ctx, "cancelOffer", nil, func(opts *bind.TransactOpts) (*types.Transaction, error) {
					return m.cfg.Chain.Contract.CancelOffer(opts, id)
				}); err != nil {
					return err
				}
				_ = m.cfg.Store.Update(func(st *swap2.State) { st.Offers[o.OfferID].Active = false })
				continue
			}
			live++
			continue
		}
		// not active on chain: taken or cancelled. Look for the take.
		// SwapCreated(swapId indexed, offerId indexed, taker indexed, ...): offerId is topic[2].
		logs, err := m.cfg.Chain.Logs(ctx, o.PostedBlock, head,
			[][]ethcommon.Hash{{m.cfg.Chain.EventID("SwapCreated")}, nil, {id}})
		if err != nil {
			return err
		}
		if len(logs) == 0 {
			log.Printf("maker: offer %s inactive with no take (cancelled?)", short(o.OfferID))
			_ = m.cfg.Store.Update(func(st *swap2.State) {
				st.Offers[o.OfferID].Active = false
				st.Offers[o.OfferID].TakeChecked = true
			})
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
func (m *Maker) recordTake(ctx context.Context, o *swap2.OfferRec, ev *xmrswap2.XmrSwapV2SwapCreated) error {
	swapID := hex.EncodeToString(ev.SwapId[:])
	log.Printf("maker: offer %s taken by %s for %s ETH (%s XMR), swap %s",
		short(o.OfferID), ev.Taker.Hex(), swap2.FmtETH(ev.Value), swap2.FmtXMR(ev.XmrPiconero.Uint64()), short(swapID))

	height, err := m.cfg.Wallet.GetHeight()
	if err != nil {
		return fmt.Errorf("xmr height: %w", err)
	}
	restore := height
	if restore > 20 {
		restore -= 20
	}

	myPub, err := swap2.ParsePub(o.SpendPub)
	if err != nil {
		return err
	}
	myView, err := swap2.ParseViewPriv(o.ViewPriv)
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
	addr, _ := swap2.SharedAddress(m.cfg.Env, myPub, theirPub, myView, theirView)

	err = m.cfg.Store.Update(func(st *swap2.State) {
		st.Offers[o.OfferID].Active = false
		st.Offers[o.OfferID].SwapID = swapID
		st.Swaps[swapID] = &swap2.SwapRec{
			SwapID:        swapID,
			OfferID:       o.OfferID,
			Role:          "maker",
			State:         swap2.StateXMRLockPending,
			Asset:         "ETH",
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

// postOffer generates fresh one-time keys and posts a single-use offer, locking the bond.
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
				swap2.FmtXMR(bal.UnlockedBalance), swap2.FmtXMR(need.Uint64()))
			return nil
		}
	}

	// The bond is real money leaving the gas wallet. Refuse rather than fail mid-transaction.
	ethBal, err := m.cfg.Chain.Balance(ctx)
	if err != nil {
		return fmt.Errorf("eth balance: %w", err)
	}
	needETH := new(big.Int).Add(m.bond, big.NewInt(2e14)) // bond + gas headroom
	if ethBal.Cmp(needETH) < 0 {
		log.Printf("maker: not posting: gas wallet holds %s ETH, needs %s (bond %s + gas). Fund %s.",
			swap2.FmtETH(ethBal), swap2.FmtETH(needETH), swap2.FmtETH(m.bond), m.cfg.Chain.From.Hex())
		return nil
	}

	kp, err := mcrypto.GenerateKeys()
	if err != nil {
		return err
	}
	spendPriv, viewPriv, spendPub := swap2.HexKeys(kp)
	pending := "pending:" + spendPub

	// keys hit disk before the transaction leaves this machine
	if err := m.cfg.Store.Update(func(st *swap2.State) {
		st.Offers[pending] = &swap2.OfferRec{
			OfferID: pending, SpendPriv: spendPriv, ViewPriv: viewPriv, SpendPub: spendPub, Asset: "ETH",
		}
	}); err != nil {
		return err
	}

	expiry := uint64(time.Now().Add(m.cfg.OfferTTL).Unix())
	params := xmrswap2.XmrSwapV2OfferParams{
		MinAmount:        m.cfg.MinWei,
		MaxAmount:        m.cfg.MaxWei,
		XmrPerAsset:      m.cfg.XmrPerAsset,
		Expiry:           expiry,
		Timeout1Duration: uint64(m.cfg.Timeout1.Seconds()),
		Timeout2Duration: uint64(m.cfg.Timeout2.Seconds()),
		MakerSpendPub:    swap2.Bytes32(kp.PublicKeyPair().SpendKey().Bytes()),
		MakerViewPriv:    swap2.Bytes32(kp.ViewKey().Bytes()),
		Payout:           m.cfg.Payout,
		// Fixed price, so every oracle field stays zero. The contract rejects an offer that
		// sets both a price and a feed.
		Oracle:         ethcommon.Address{},
		OracleRatioBps: big.NewInt(0),
		OracleOffset:   big.NewInt(0),
		OracleMaxAge:   0,
		OracleMaxPrice: big.NewInt(0),
	}

	receipt, err := m.cfg.Chain.Send(ctx, "postOffer", m.bond, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return m.cfg.Chain.Contract.PostOffer(opts, params)
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
	log.Printf("maker: posted offer %s: %s to %s ETH, bond %s ETH locked, expires %s",
		short(offerID), swap2.FmtETH(m.cfg.MinWei), swap2.FmtETH(m.cfg.MaxWei), swap2.FmtETH(m.bond),
		time.Unix(int64(expiry), 0).Format(time.RFC3339))

	return m.cfg.Store.Update(func(st *swap2.State) {
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
			_ = m.cfg.Store.Update(func(st *swap2.State) { st.Swaps[swapID].LastError = err.Error() })
		}
	}()
}

func (m *Maker) swap(swapID string) *swap2.SwapRec {
	var cp swap2.SwapRec
	m.cfg.Store.View(func(st *swap2.State) { cp = *st.Swaps[swapID] })
	return &cp
}

func (m *Maker) setState(swapID, state string, mut func(*swap2.SwapRec)) error {
	log.Printf("maker: swap %s -> %s", short(swapID), state)
	return m.cfg.Store.Update(func(st *swap2.State) {
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
	backoff := m.cfg.Poll
	for {
		s := m.swap(swapID)
		var err error
		switch s.State {
		case swap2.StateXMRLockPending:
			err = m.lockXMR(ctx, s)
		case swap2.StateXMRLocked:
			err = m.awaitClaim(ctx, s)
		case swap2.StateRecovering:
			err = m.recover(ctx, s)
		case swap2.StateClaimed, swap2.StateRecovered:
			return nil
		default:
			return fmt.Errorf("unknown state %q", s.State)
		}
		if err == nil {
			backoff = m.cfg.Poll
			continue
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// A failed call is not a failed swap. Record it, wait, try the same state again.
		log.Printf("maker: swap %s in %s: %s (retrying in %s)", short(swapID), s.State, err, backoff)
		_ = m.cfg.Store.Update(func(st *swap2.State) { st.Swaps[swapID].LastError = err.Error() })
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		if backoff < 5*time.Minute {
			backoff *= 2
		}
	}
}

func (m *Maker) lockXMR(ctx context.Context, s *swap2.SwapRec) error {
	id := hashFromHex(s.SwapID)
	onchain, err := m.cfg.Chain.Contract.Swaps(&bind.CallOpts{Context: ctx}, id)
	if err != nil {
		return err
	}
	_, now, err := m.cfg.Chain.Head(ctx)
	if err != nil {
		return err
	}
	if onchain.Stage != StagePending || now+uint64(lockGrace.Seconds()) > s.Timeout1 {
		// Too late to lock safely: the taker needs confirmations before t1.
		log.Printf("maker: swap %s: not locking XMR (stage=%d, %ds to t1); waiting for taker refund",
			short(s.SwapID), onchain.Stage, int64(s.Timeout1)-int64(now))
		return m.setState(s.SwapID, swap2.StateXMRLocked, nil) // awaitClaim handles refund/recovery
	}

	addr, err := mcrypto.NewAddress(s.SharedAddress, m.cfg.Env)
	if err != nil {
		return err
	}

	if m.cfg.ManualXMR {
		fmt.Printf("\n=== ACTION NEEDED ===\nSend exactly %s XMR to:\n%s\nfrom your own wallet, before %s.\n=====================\n\n",
			swap2.FmtXMR(s.XmrPiconero), addr, time.Unix(int64(s.Timeout1), 0).Add(-25*time.Minute).Format(time.RFC3339))
		myView, err := swap2.ParseViewPriv(s.MyViewPriv)
		if err != nil {
			return err
		}
		theirView, err := swap2.ParseViewPriv(s.TheirViewPriv)
		if err != nil {
			return err
		}
		viewAB := mcrypto.SumPrivateViewKeys(myView, theirView)
		ok, err := swap2.WaitForXMR(ctx, m.cfg.Wallet, "maker-watch-"+short(s.SwapID), viewAB, addr, s.RestoreHeight,
			s.XmrPiconero, time.Unix(int64(s.Timeout1), 0).Add(-5*time.Minute))
		if err != nil {
			return err
		}
		if !ok {
			log.Printf("maker: swap %s: XMR was not paid in time; the taker will refund, nothing is lost", short(s.SwapID))
			return m.setState(s.SwapID, swap2.StateXMRLocked, nil)
		}
		return m.setState(s.SwapID, swap2.StateXMRLocked, func(r *swap2.SwapRec) { r.XmrLockTxID = "manual" })
	}

	log.Printf("maker: swap %s: sending %s XMR to shared address %s", short(s.SwapID), swap2.FmtXMR(s.XmrPiconero), addr)
	tr, err := m.cfg.Wallet.Transfer(ctx, addr, 0, coins.NewPiconeroAmount(s.XmrPiconero), monero.MinSpendConfirmations)
	if err != nil {
		return fmt.Errorf("xmr transfer: %w", err)
	}
	return m.setState(s.SwapID, swap2.StateXMRLocked, func(r *swap2.SwapRec) { r.XmrLockTxID = tr.TxID })
}

// awaitClaim claims once the taker sets ready (or t1 passes), or recovers if the taker refunds.
func (m *Maker) awaitClaim(ctx context.Context, s *swap2.SwapRec) error {
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
						return m.setState(s.SwapID, swap2.StateRecovering, func(r *swap2.SwapRec) { r.TheirSecret = secret })
					}
					return fmt.Errorf("swap completed on-chain but no Refunded/Claimed event found")
				}
				if now < s.Timeout2 && (onchain.Stage == StageReady || now >= s.Timeout1) {
					mySecret, err := swap2.ParseSpendPriv(s.MySpendPriv)
					if err != nil {
						return err
					}
					receipt, err := m.cfg.Chain.Send(ctx, "claim", nil, func(opts *bind.TransactOpts) (*types.Transaction, error) {
						return m.cfg.Chain.Contract.Claim(opts, id, swap2.Bytes32(mySecret.Bytes()))
					})
					if err != nil {
						log.Printf("maker: claim failed (will re-check): %s", err)
					} else {
						// The payment went to the payout address; the bond is now a credit.
						if cerr := m.collect(ctx); cerr != nil {
							log.Printf("maker: bond withdraw failed (will retry on the next tick): %s", cerr)
						}
						return m.setState(s.SwapID, swap2.StateClaimed, func(r *swap2.SwapRec) { r.EthTxHash = receipt.TxHash.Hex() })
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
func (m *Maker) findSecret(ctx context.Context, event string, s *swap2.SwapRec, head uint64) (string, error) {
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
func (m *Maker) recover(ctx context.Context, s *swap2.SwapRec) error {
	if s.XmrLockTxID == "" {
		return m.setState(s.SwapID, swap2.StateRecovered, nil) // we never locked
	}
	mySpend, err := swap2.ParseSpendPriv(s.MySpendPriv)
	if err != nil {
		return err
	}
	theirSpend, err := swap2.ParseSpendPriv(s.TheirSecret)
	if err != nil {
		return err
	}
	myView, err := swap2.ParseViewPriv(s.MyViewPriv)
	if err != nil {
		return err
	}
	theirView, err := swap2.ParseViewPriv(s.TheirViewPriv)
	if err != nil {
		return err
	}
	spendAB := mcrypto.SumPrivateSpendKeys(mySpend, theirSpend)
	viewAB := mcrypto.SumPrivateViewKeys(myView, theirView)
	log.Printf("maker: swap %s: taker refunded; sweeping XMR back to %s", short(s.SwapID), m.cfg.Wallet.PrimaryAddress())
	ids, err := swap2.SweepShared(ctx, m.cfg.Wallet, "maker-recover-"+short(s.SwapID), spendAB, viewAB, s.RestoreHeight, m.cfg.Wallet.PrimaryAddress())
	if err != nil {
		return err
	}
	// The bond comes back on a refund too — the contract never seizes it.
	if cerr := m.collect(ctx); cerr != nil {
		log.Printf("maker: bond withdraw failed (will retry on the next tick): %s", cerr)
	}
	return m.setState(s.SwapID, swap2.StateRecovered, func(r *swap2.SwapRec) { r.SweepTxIDs = ids })
}

// ---- small helpers

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
