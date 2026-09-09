# monero-swap

Trade Monero for ETH or USDC with a stranger and trust nobody. The contract holds the ETH, the Monero sits at an address both of you control, and the only way the seller gets paid is by handing you the key.

Built on the [Athanor](https://github.com/AthanorLabs/atomic-swap) protocol (ChainSafe, 2023), which proved the cryptography on mainnet and then went quiet. This is the part they didn't finish: no peer network to die, no daemon holding your savings, a web page for the buyer, and a fee so someone has a reason to keep it alive.

**Status: test networks only.** Three swaps and one refund have run end to end on Sepolia + Monero stagenet. Nothing is audited. Nothing is on Base yet. Don't put real money near it.

## How a swap works

1. A seller runs a small program that posts an offer on the contract: how much Monero, at what price, min and max.
2. A buyer opens the web page, picks an offer, and pays ETH into the contract from their own wallet. The page makes a one-time Monero key for this swap and keeps it in the browser.
3. The seller's program sees the payment and sends the Monero to a fresh address that needs *both* one-time keys to spend. Neither side can move it alone.
4. The buyer's page watches that address. After 10 confirmations it flips the contract to *ready*.
5. The seller collects the ETH. Collecting reveals the seller's one-time key on-chain — the contract checks it against the Monero key before paying.
6. The buyer's page reads that key, adds it to its own, and now holds the full spend key for the Monero.

If the seller never sends, the buyer's ETH refunds after an hour. If the buyer vanishes, the seller collects after the deadline anyway. If the buyer refunds, that reveals *their* key and the seller takes the Monero back. Nobody can end up with nothing.

The fee — 0.15% of the ETH side — comes out of the seller's payment at step 5, to a wallet fixed when the contract was deployed. Refunds pay no fee.

## What's in here

| Path | What | License |
|---|---|---|
| `contracts-v2/` | `XmrSwap.sol`, the Solidity contract. No owner, no upgrade. 36 Foundry tests. | MIT |
| `web/` | The buyer page. Vanilla JS, ethers, noble ed25519. Runs in any browser with MetaMask. | MIT |
| `relay/` | A 40-line relay so the page can read a Monero node from a browser. Cloudflare Worker or plain Node. | MIT |
| `swap2/`, `cmd/monero-swap/` | The seller program (and a command-line buyer for testing). Go. | LGPL-3.0 |
| everything else | Athanor's original code, kept for its Monero wallet library. | LGPL-3.0 |

## Buy Monero

You need MetaMask (or any browser wallet) with ETH on the right network, and a Monero wallet to receive into.

Open the page, connect the wallet, pick an offer, pay. Keep the tab open until step 5 above. When the seller collects, the page shows you a private spend key and view key. Restore them as a wallet *from keys* in Feather, Cake or the Monero GUI and send the coins wherever you like. Nobody else can spend them.

Download the key backup the page offers. If you clear the site's storage before the swap finishes, the Monero is gone.

## Sell Monero

You run one program on a machine that stays on. It holds a little ETH for gas (a few dollars) and either holds Monero to sell or asks you to pay each swap by hand from your own wallet (`--manual-xmr`).

```bash
git clone https://github.com/dafarusd/monero-swap.git
cd monero-swap
./scripts/install-monero-linux.sh          # Monero wallet tools into ./monero-bin
make build-swap                            # needs Go 1.21 — newer Go breaks the network layer
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract CONTRACT_ADDRESS \
  --monerod-host node.monerodevs.org --monerod-port 18089 addresses
```

That prints a gas wallet address and a Monero wallet address. Fund the gas wallet with a little ETH. Then:

```bash
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract CONTRACT_ADDRESS \
  --monerod-host node.monerodevs.org --monerod-port 18089 \
  maker run --payout YOUR_ETH_ADDRESS --min 0.01 --max 0.1 --price 15.5
```

`--price` is how much Monero the buyer gets per 1 ETH. `--payout` is where your ETH goes — any wallet, it never touches the gas key. Leave it running. It keeps one offer live, serves whoever takes it, and picks up where it left off after a restart.

`CONTRACT_ADDRESS` is filled in here once the Base deploy is done.

## Run a relay

Browsers can't call public Monero nodes directly, so the page reads the chain through a relay that forwards a short list of read-only calls. Yours or mine, it doesn't matter — it holds nothing.

```bash
cd relay && npx wrangler deploy      # Cloudflare, free tier
# or
UPSTREAM=http://node.monerodevs.org:18089 node relay/local.mjs
```

Then put the relay URL in the box at the top of the page.

## Limits, honestly

- Not audited. The contract is under 400 lines and the tests pass, but nobody outside this repo has read it.
- The buyer page can't send Monero yet. It hands you the keys; your own wallet does the sending.
- The seller's program has to be online. That's not a bug, it's Monero: only a private key can move it, and a key has to live somewhere.
- Offers are single-use. Each swap gets fresh keys, and the contract refuses reused ones.
- Base only, for now. Same contract works on any Ethereum-style chain; each one splits the sellers.
- One person built this, with AI writing most of the code against a spec and every swap replayed on test networks. The run logs and transaction hashes are in the commit messages. Don't take my word for it.

Built by [Dafarus](https://x.com/Dafarusd).
