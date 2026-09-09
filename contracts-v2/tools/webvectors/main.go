// webvectors prints JSON test vectors for the browser Monero code: key pairs with
// addresses on both networks, and a shared (summed) address for two of them.
package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/athanorlabs/atomic-swap/common"
	mcrypto "github.com/athanorlabs/atomic-swap/crypto/monero"
)

type pair struct {
	SpendPriv string `json:"spend_priv"`
	ViewPriv  string `json:"view_priv"`
	SpendPub  string `json:"spend_pub"`
	ViewPub   string `json:"view_pub"`
	Stagenet  string `json:"address_stagenet"`
	Mainnet   string `json:"address_mainnet"`
}

func mk() (pair, *mcrypto.PrivateKeyPair) {
	kp, err := mcrypto.GenerateKeys()
	if err != nil {
		panic(err)
	}
	pk := kp.PublicKeyPair()
	return pair{
		SpendPriv: kp.SpendKey().Hex(),
		ViewPriv:  kp.ViewKey().Hex(),
		SpendPub:  pk.SpendKey().Hex(),
		ViewPub:   pk.ViewKey().Hex(),
		Stagenet:  pk.Address(common.Stagenet).String(),
		Mainnet:   pk.Address(common.Mainnet).String(),
	}, kp
}

func main() {
	out := map[string]any{}
	var pairs []pair
	var kps []*mcrypto.PrivateKeyPair
	for i := 0; i < 4; i++ {
		p, kp := mk()
		pairs = append(pairs, p)
		kps = append(kps, kp)
	}
	out["pairs"] = pairs
	// shared address of pair 0 + pair 1
	spend := mcrypto.SumPublicKeys(kps[0].PublicKeyPair().SpendKey(), kps[1].PublicKeyPair().SpendKey())
	view := mcrypto.SumPrivateViewKeys(kps[0].ViewKey(), kps[1].ViewKey())
	spendPriv := mcrypto.SumPrivateSpendKeys(kps[0].SpendKey(), kps[1].SpendKey())
	shared := mcrypto.NewPublicKeyPair(spend, view.Public())
	out["shared01"] = map[string]string{
		"spend_pub":        hex.EncodeToString(spend.Bytes()),
		"view_priv":        view.Hex(),
		"spend_priv":       spendPriv.Hex(),
		"address_stagenet": shared.Address(common.Stagenet).String(),
		"address_mainnet":  shared.Address(common.Mainnet).String(),
	}
	// view key must be derivable from spend key the Monero way
	v, err := kps[2].SpendKey().View()
	if err != nil {
		panic(err)
	}
	out["derived_view_2"] = v.Hex()
	b, _ := json.MarshalIndent(out, "", "  ")
	if err := os.WriteFile("web/test/fixtures/vectors.json", b, 0o644); err != nil {
		panic(err)
	}
	fmt.Println("wrote web/test/fixtures/vectors.json")
}
