# monero-swap

Trade Monero for ETH or USDC with a stranger and trust nobody. The contract holds the ETH, the Monero sits at an address both of you control, and the only way the seller gets paid is by handing you the key.

Built on [Athanor](https://github.com/AthanorLabs/atomic-swap) (ChainSafe, 2023), which proved the cryptography on mainnet and then went quiet. This is the part they didn't finish: no peer network to die, no daemon holding your savings, a web page for the buyer, and a fee so someone has a reason to keep it alive. What's theirs and what changed is spelled out in [Credit](#credit) below.

**Status: live on Base, with a real offer on the board.** The contract is deployed on Base mainnet and the first seller — me — has Monero up for sale. Five swaps and one refund ran end to end on test networks first; no real-money swap has completed yet. Nothing is audited. Start small.

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

You need MetaMask (or any browser wallet) with ETH on Base, and a Monero wallet to receive into.

Open **[the page](https://dafarusd.github.io/monero-swap/)**, connect the wallet, pick an offer, pay. Keep the tab open until step 5 above. When the seller collects, the page shows you a private spend key and view key. Restore them as a wallet *from keys* in Feather, Cake or the Monero GUI and send the coins wherever you like. Nobody else can spend them.

Download the key backup the page offers. If you clear the site's storage before the swap finishes, the Monero is gone.

## Sell Monero

You run one program on a machine that stays on. It holds a little ETH for gas (a few dollars) and either holds Monero to sell or asks you to pay each swap by hand from your own wallet (`--manual-xmr`). A Raspberry Pi 5 is enough: the install script fetches the arm64 Monero tools and the program cross-compiles with `GOARCH=arm64`. I'm moving mine to one.

```bash
git clone https://github.com/dafarusd/monero-swap.git
cd monero-swap
./scripts/install-monero-linux.sh          # Monero wallet tools into ./monero-bin
make build-swap                            # needs Go 1.21 — newer Go breaks the network layer
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract 0x67fe8681563F37f2A8BBed84C85784a678FeC693 \
  --monerod-host node.monerodevs.org --monerod-port 18089 addresses
```

That prints a gas wallet address and a Monero wallet address. Fund the gas wallet with a little ETH. Then:

```bash
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract 0x67fe8681563F37f2A8BBed84C85784a678FeC693 \
  --monerod-host node.monerodevs.org --monerod-port 18089 \
  maker run --payout YOUR_ETH_ADDRESS --min 0.01 --max 0.1 --price 15.5
```

`--price` is how much Monero the buyer gets per 1 ETH. `--payout` is where your ETH goes — any wallet, it never touches the gas key. Leave it running. It keeps one offer live, serves whoever takes it, and picks up where it left off after a restart.

To sell for a token instead, add `--asset` with the token's address and quote `--min`, `--max` and `--price` in that token:

```bash
  --asset 0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913 --min 20 --max 500 --price 0.0066   # USDC on Base
```

Buyers of a token offer approve the contract once, then pay. The fee comes out in the same token.

To see who else is selling, list the open offers — every one carries its seller's address:

```bash
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract 0x67fe8681563F37f2A8BBed84C85784a678FeC693   --monerod-host node.monerodevs.org --monerod-port 18089 offers
```

## Deployments

| Chain | Contract | Fee |
|---|---|---|
| **Base** (mainnet) | `0x67fe8681563F37f2A8BBed84C85784a678FeC693` | 0.15% to `0x16F57804d30CF991cC006e0713294A6D52778260` |
| Sepolia (test) | `0xB96bDd5834F455C1A6edA15e5bAF25eFd506d61E` | 0.15% to a throwaway wallet |
| Base Sepolia (test) | `0x97f8A483cFa8680F67aC24D83bbe4Fbc4f250755` | 0.15% to a throwaway wallet |

The Base contract has no owner and can't be changed. Source is verified on [Basescan](https://basescan.org/address/0x67fe8681563f37f2a8bbed84c85784a678fec693#code), so what you read there is what runs. Test contracts pair with Monero stagenet.

## Run a relay

Browsers can't call public Monero nodes directly, so the page reads the chain through a relay that forwards a short list of read-only calls. It holds nothing and sees only block data. The page comes prefilled with mine — `monero-relay.dafarusd.workers.dev` for Base, `monero-relay-stagenet.dafarusd.workers.dev` for the test networks — and you can swap in your own:

```bash
cd relay && npx wrangler deploy      # Cloudflare, free tier
# or
UPSTREAM=http://node.monerodevs.org:18089 node relay/local.mjs
```

Then put the relay URL in the box at the top of the page.

## Credit

This repo is a fork of [AthanorLabs/atomic-swap](https://github.com/AthanorLabs/atomic-swap), archived by its owners on 1 September 2026. Their work is why this exists:

- **The protocol.** Buyer locks on Ethereum, seller locks Monero to a shared key, the claim reveals the secret. Joël Gugger's [2020 paper](https://eprint.iacr.org/2020/1126) designed it for Bitcoin; Athanor carried it to Ethereum and ran it on mainnet in June 2023. Nothing about that dance is mine.
- **The Monero wallet library.** `monero/` and `crypto/monero/` are theirs, untouched except for one added constructor. The seller program drives `monero-wallet-rpc` through their code.
- **The first test.** The first swap I ran, on 8 September 2026, was their unchanged code on Sepolia and stagenet. It worked, built with Go 1.21.

What changed, and why:

- **New contract.** `contracts-v2/XmrSwap.sol` replaces their `SwapCreator.sol`. The offer board moved on-chain, so there's no peer-to-peer network and no bootnodes to die — theirs were all dead by the time I tried. A fee, fixed at deploy, pays whoever keeps this running. The claim checks the revealed secret against the Monero key directly with an on-chain ed25519 multiply, which is cheap on Base; that removes the secp256k1 side and the cross-curve DLEq proof, the heaviest part of their client. One-time keys can't be reused across swaps. A seller's payout goes to any wallet, never the gas key.
- **New seller and buyer programs.** `swap2/` and `cmd/monero-swap/` are new. Their `swapd` is still in the tree but nothing here runs it. A failed call retries instead of ending the swap — a node outage killed their buyer's watcher mid-swap in my first run, and only their restart recovery saved it.
- **A browser buyer.** `web/` is new. Their UI was unmaintained; this one needs only MetaMask.
- **A relay.** `relay/` is new, because browsers can't call Monero nodes.
- **Kept their license** for the Go code (LGPL-3.0). The contract, page and relay are MIT. The ed25519 library is Jan Vornberger's, MIT, via Wrapsynth.

If you learned something here, the people to thank are noot, dimalinux and the ChainSafe team.

## Limits, honestly

- Not audited. The contract is under 400 lines and the tests pass, but nobody outside this repo has read it.
- The buyer page can't send Monero yet. It hands you the keys; your own wallet does the sending.
- The seller's program has to be online. That's not a bug, it's Monero: only a private key can move it, and a key has to live somewhere.
- Offers are single-use. Each swap gets fresh keys, and the contract refuses reused ones.
- Base only, for now. Same contract works on any Ethereum-style chain; each one splits the sellers.
- One person built this, and every swap was replayed on test networks first. The test contracts are listed above and every swap is on those chains for anyone to read. Don't take my word for it.

Built by [Dafarus](https://x.com/Dafarusd).
