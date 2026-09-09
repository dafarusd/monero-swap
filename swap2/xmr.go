package swap2

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/athanorlabs/atomic-swap/common"
	mcrypto "github.com/athanorlabs/atomic-swap/crypto/monero"
	"github.com/athanorlabs/atomic-swap/monero"
)

// Bytes32 copies a 32-byte slice into an array.
func Bytes32(b []byte) (out [32]byte) {
	copy(out[:], b)
	return out
}

// HexKeys returns hex encodings of a one-time key pair: private spend, private view, public spend.
func HexKeys(kp *mcrypto.PrivateKeyPair) (spendPriv, viewPriv, spendPub string) {
	return kp.SpendKey().Hex(), kp.ViewKey().Hex(), kp.PublicKeyPair().SpendKey().Hex()
}

// SharedAddress derives the one-time Monero address controlled jointly by both parties:
// spend key = aSpend + bSpend, view key = aView + bView.
func SharedAddress(
	env common.Environment,
	aSpendPub, bSpendPub *mcrypto.PublicKey,
	aView, bView *mcrypto.PrivateViewKey,
) (*mcrypto.Address, *mcrypto.PrivateViewKey) {
	spendPub := mcrypto.SumPublicKeys(aSpendPub, bSpendPub)
	viewPriv := mcrypto.SumPrivateViewKeys(aView, bView)
	addr := mcrypto.NewPublicKeyPair(spendPub, viewPriv.Public()).Address(env)
	return addr, viewPriv
}

// ParseSpendPriv / ParseViewPriv / ParsePub decode hex keys.
func ParseSpendPriv(h string) (*mcrypto.PrivateSpendKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return mcrypto.NewPrivateSpendKey(b)
}

func ParseViewPriv(h string) (*mcrypto.PrivateViewKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return mcrypto.NewPrivateViewKeyFromBytes(b)
}

func ParsePub(h string) (*mcrypto.PublicKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, err
	}
	return mcrypto.NewPublicKeyFromBytes(b)
}

// WaitForXMR opens a view-only wallet on the shared address and waits until at least
// `expected` piconero are unlocked there, or the deadline passes.
func WaitForXMR(
	ctx context.Context,
	wc monero.WalletClient,
	prefix string,
	viewAB *mcrypto.PrivateViewKey,
	addr *mcrypto.Address,
	restoreHeight uint64,
	expected uint64,
	deadline time.Time,
) (bool, error) {
	var view monero.WalletClient
	err := retry(ctx, 5, 15*time.Second, "open view-only wallet", func() error {
		var e error
		view, e = monero.CreateViewOnlyWalletFromKeys(wc.CreateWalletConf(prefix), viewAB, addr, restoreHeight)
		return e
	})
	if err != nil {
		return false, err
	}
	defer view.CloseAndRemoveWallet()

	for {
		bal, err := view.GetBalance(0)
		if err != nil {
			log.Printf("xmr: balance check failed (will retry): %s", err)
		} else {
			log.Printf("xmr: shared address %s… balance=%d unlocked=%d blocks-to-unlock=%d (need %d)",
				addr.String()[:12], bal.Balance, bal.UnlockedBalance, bal.BlocksToUnlock, expected)
			if bal.UnlockedBalance >= expected {
				return true, nil
			}
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(30 * time.Second):
		}
	}
}

// SweepShared spends everything at the shared address (spend key = a + b) to `to`.
func SweepShared(
	ctx context.Context,
	wc monero.WalletClient,
	prefix string,
	spendAB *mcrypto.PrivateSpendKey,
	viewAB *mcrypto.PrivateViewKey,
	restoreHeight uint64,
	to *mcrypto.Address,
) ([]string, error) {
	kp := mcrypto.NewPrivateKeyPair(spendAB, viewAB)
	var spend monero.WalletClient
	err := retry(ctx, 5, 15*time.Second, "open spend wallet", func() error {
		var e error
		spend, e = monero.CreateSpendWalletFromKeys(wc.CreateWalletConf(prefix), kp, restoreHeight)
		return e
	})
	if err != nil {
		return nil, err
	}
	defer spend.CloseAndRemoveWallet()

	transfers, err := spend.SweepAll(ctx, to, 0, monero.SweepToSelfConfirmations)
	if err != nil {
		return nil, fmt.Errorf("sweep: %w", err)
	}
	ids := make([]string, 0, len(transfers))
	for _, t := range transfers {
		ids = append(ids, t.TxID)
	}
	return ids, nil
}

// retry runs fn up to `attempts` times, waiting `wait` between failures.
func retry(ctx context.Context, attempts int, wait time.Duration, what string, fn func() error) error {
	var err error
	for i := 1; i <= attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		log.Printf("%s failed (attempt %d of %d): %s", what, i, attempts, err)
		if i == attempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return fmt.Errorf("%s: %w", what, err)
}
