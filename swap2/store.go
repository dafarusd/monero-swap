// Package swap2 implements the maker and taker sides of XmrSwap: Monero <-> ETH/ERC-20
// atomic swaps with the offer board on-chain and no peer-to-peer layer.
package swap2

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Maker swap states.
const (
	StateXMRLockPending = "xmr_lock_pending" // taker locked ETH; we must send XMR to the shared address
	StateXMRLocked      = "xmr_locked"       // XMR confirmed; waiting for ready / t1 to claim
	StateClaimed        = "claimed"          // ETH collected
	StateRecovering     = "recovering"       // taker refunded; sweeping our XMR back
	StateRecovered      = "recovered"
)

// Taker swap states.
const (
	StateETHLocked = "eth_locked" // waiting for XMR at the shared address
	StateXMRSeen   = "xmr_seen"   // XMR arrived and unlocked; about to set ready
	StateReady     = "ready"      // ready set; waiting for maker's claim (reveals secret)
	StateSwept     = "swept"      // we spent the XMR to the payout address
	StateRefunded  = "refunded"   // we took the ETH back
)

// OfferRec is an offer we posted, with the one-time keys behind it.
type OfferRec struct {
	OfferID     string `json:"offer_id"`
	SpendPriv   string `json:"spend_priv"`
	ViewPriv    string `json:"view_priv"`
	SpendPub    string `json:"spend_pub"`
	Asset       string `json:"asset"`
	PostedBlock uint64 `json:"posted_block"`
	TxHash      string `json:"tx_hash"`
	Active      bool   `json:"active"`
	SwapID      string `json:"swap_id,omitempty"`
}

// SwapRec is one swap in progress or finished, from our side.
type SwapRec struct {
	SwapID        string   `json:"swap_id"`
	OfferID       string   `json:"offer_id"`
	Role          string   `json:"role"` // "maker" or "taker"
	State         string   `json:"state"`
	Asset         string   `json:"asset"`
	ValueWei      string   `json:"value_wei"`
	XmrPiconero   uint64   `json:"xmr_piconero"`
	MySpendPriv   string   `json:"my_spend_priv"`
	MyViewPriv    string   `json:"my_view_priv"`
	TheirSpendPub string   `json:"their_spend_pub"`
	TheirViewPriv string   `json:"their_view_priv"`
	SharedAddress string   `json:"shared_address"`
	RestoreHeight uint64   `json:"restore_height"`
	Timeout1      uint64   `json:"timeout1"`
	Timeout2      uint64   `json:"timeout2"`
	CreatedBlock  uint64   `json:"created_block"`
	XmrLockTxID   string   `json:"xmr_lock_txid,omitempty"`
	TheirSecret   string   `json:"their_secret,omitempty"`
	EthTxHash     string   `json:"eth_tx_hash,omitempty"`
	SweepTxIDs    []string `json:"sweep_txids,omitempty"`
	LastError     string   `json:"last_error,omitempty"`
}

// State is everything persisted to disk.
type State struct {
	Offers map[string]*OfferRec `json:"offers"`
	Swaps  map[string]*SwapRec  `json:"swaps"`
}

// Store is a JSON file written atomically on every change.
type Store struct {
	path string
	mu   sync.Mutex
	st   *State
}

// OpenStore loads or creates the state file.
func OpenStore(path string) (*Store, error) {
	st := &State{Offers: map[string]*OfferRec{}, Swaps: map[string]*SwapRec{}}
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := json.Unmarshal(data, st); err != nil {
			return nil, err
		}
		if st.Offers == nil {
			st.Offers = map[string]*OfferRec{}
		}
		if st.Swaps == nil {
			st.Swaps = map[string]*SwapRec{}
		}
	case errors.Is(err, os.ErrNotExist):
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	s := &Store{path: path, st: st}
	return s, s.save()
}

// Update applies fn under the lock and writes the file.
func (s *Store) Update(fn func(*State)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.st)
	return s.save()
}

// View runs fn under the lock without writing.
func (s *Store) View(fn func(*State)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(s.st)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.st, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
