# monero-swap

Trade Monero for ETH with a stranger and trust nobody. The contract holds the ETH, the Monero sits at an address both of you control, and the only way the seller gets paid is by handing you the key.

Built on [Athanor](https://github.com/AthanorLabs/atomic-swap) (ChainSafe, 2023), which proved the cryptography on mainnet and then went quiet. This is the part they didn't finish: no peer network to die, no daemon holding your savings, a web page for the buyer, and a fee so someone has a reason to keep it alive. What's theirs and what changed is spelled out in [Credit](#credit) below.

**Status: v2 is live on Base and there's a seller you can run. Nobody is running one right now.** The contract is deployed and verified on Base mainnet, the web page talks to it, and `monero-swap-v2` is the seller. I'm not hosting one at the moment, so the board is empty until somebody does — that can be you.

Five swaps and one refund ran end to end on test networks, all against v1. No real-money swap has completed on either version, and the v2 seller has not yet carried a swap from start to finish — it builds and its logic is the v1 seller's, which did run real swaps, but that is not the same as proof. Nothing is audited. Start small.

## How a swap works

1. A seller runs a small program that posts an offer on the contract: how much Monero, at what price, min and max.
2. A buyer opens the web page, picks an offer, and pays ETH into the contract from their own wallet. The page makes a one-time Monero key for this swap and keeps it in the browser.
3. The seller's program sees the payment and sends the Monero to a fresh address that needs *both* one-time keys to spend. Neither side can move it alone.
4. The buyer's page watches that address. After 10 confirmations it flips the contract to *ready*.
5. The seller collects the ETH. Collecting reveals the seller's one-time key on-chain — the contract checks it against the Monero key before paying.
6. The buyer's page reads that key, adds it to its own, and now holds the full spend key for the Monero.

If the seller never sends, the buyer takes their ETH back — any time before the first deadline, and again after the second one. Between the two, the seller can still collect. If the buyer vanishes, the seller collects anyway. If the buyer refunds, that reveals *their* key and the seller takes the Monero back. Nobody can end up with nothing.

Both deadlines are at least 24 hours in v2, up from an hour in v1. Base runs one sequencer, and forcing a transaction in from L1 takes about twelve hours — a one-hour window can close while you're being censored, before you can do anything about it.

The fee — 0.15% of the ETH side — comes out of the seller's payment at step 5 and splits in half: half to me, half to a dev fund so other people working on this can be paid out of it. Both are multisig Safes, fixed when the contract was deployed, and neither can be redirected. The fee isn't sent anywhere during a swap; it's credited, and whoever it belongs to withdraws later. That way a fee address that can't receive ETH can never block someone else's claim. Refunds pay no fee.

A seller also posts a bond to list an offer — 5% of the largest amount the offer accepts. It's held while the offer is live and returned in full when the swap settles, whichever way it settles. It's never paid to anyone else: this contract can't see Monero, so it can't tell a seller who didn't deliver from a buyer who walked away, and a bond it can't award fairly is one it has no business seizing. It's there to make spamming the board cost something.

## What's in here

| Path | What | License |
|---|---|---|
| `contracts-v2/` | `XmrSwapV2.sol`, the live contract, plus `XmrSwap.sol` (v1). No owner, no upgrade. 83 Foundry tests, including two stateful invariant suites. | MIT |
| `web/` | The buyer page. Vanilla JS, ethers, noble ed25519. Runs in any browser with MetaMask. | MIT |
| `relay/` | A 40-line relay so the page can read a Monero node from a browser. Cloudflare Worker or plain Node. | MIT |
| `swap2v2/`, `cmd/monero-swap-v2/` | The seller for the live v2 contract. Go. | LGPL-3.0 |
| `xmrswap2/` | Generated Go bindings for `XmrSwapV2`. Rebuild with `make bindings-v2`. | LGPL-3.0 |
| `swap2/`, `cmd/monero-swap/` | The v1 seller (and a command-line buyer for testing). Go. | LGPL-3.0 |
| everything else | Athanor's original code, kept for its Monero wallet library. | LGPL-3.0 |

## Buy Monero

You need MetaMask (or any browser wallet) with ETH on Base, and a Monero wallet to receive into.

Open **[the page](https://dafarusd.github.io/monero-swap/)**, connect the wallet, pick an offer, pay. Keep the tab open until step 5 above. When the seller collects, the page shows you a private spend key and view key. Restore them as a wallet *from keys* in Feather, Cake or the Monero GUI and send the coins wherever you like. Nobody else can spend them.

Download the key backup the page offers. If you clear the site's storage before the swap finishes, the Monero is gone.

## Sell Monero

You run one program on a machine that stays on. It holds a little ETH for gas and the bond, and either holds Monero to sell or asks you to pay each swap by hand from your own wallet (`--manual-xmr`). A Raspberry Pi 5 is enough.

```bash
git clone https://github.com/dafarusd/monero-swap.git
cd monero-swap
./scripts/install-monero-linux.sh          # Monero wallet tools into ./monero-bin
make build-swap-v2                         # any recent Go; no version pin needed
./bin/monero-swap-v2 --eth-rpc https://mainnet.base.org --contract 0xC2b2e8D385309d6552657c0b80434ca616DE12fC \
  --monerod-host node.monerodevs.org --monerod-port 18089 addresses
```

That prints an ETH address and a Monero address, and tells you what the contract is already holding for you. Read the terms before funding anything — they're fixed in the contract and nobody can change them:

```bash
./bin/monero-swap-v2 --eth-rpc https://mainnet.base.org --contract 0xC2b2e8D385309d6552657c0b80434ca616DE12fC terms
```

**Fund the ETH address with gas *and* bond.** This is the one thing that's different from v1: posting an offer locks 5% of the offer's ceiling as a bond. A `--max 0.1` offer locks 0.05 × 0.1 = 0.005 ETH. You always get it back, whichever way the swap ends — the contract never seizes it — but it comes back as a *credit* inside the contract rather than a transfer, so something has to collect it. `maker run` does that on its own; `withdraw` does it by hand.

Then leave it running:

```bash
./bin/monero-swap-v2 --eth-rpc https://mainnet.base.org --contract 0xC2b2e8D385309d6552657c0b80434ca616DE12fC \
  --monerod-host node.monerodevs.org --monerod-port 18089 \
  maker run --payout YOUR_ETH_ADDRESS --min 0.01 --max 0.1 --price 15.5
```

`--price` is how much Monero the buyer gets per 1 ETH. `--payout` is where your ETH goes — any wallet, it never touches the gas key. It keeps one offer live, serves whoever takes it, collects your bond, and picks up where it left off after a restart.

Timeouts default to 24 hours each and the contract rejects anything shorter, so a swap you start is a swap you should expect to babysit for a day. That's deliberate: Base runs one sequencer, and forcing a transaction in from L1 takes about twelve hours.

To see the board — read straight from contract storage, no log scanning:

```bash
./bin/monero-swap-v2 --eth-rpc https://mainnet.base.org --contract 0xC2b2e8D385309d6552657c0b80434ca616DE12fC offers
```

### Selling against v1

The older `monero-swap` binary drives the **v1** contract only. v2 changed the calls it makes — `takeOffer` lost an argument, offers are read from storage instead of logs, and tokens are gone — so pointing it at the v2 address won't work. It needs Go 1.21, because it pulls in Athanor's peer-to-peer stack.

You run one program on a machine that stays on. It holds a little ETH for gas (a few dollars) and either holds Monero to sell or asks you to pay each swap by hand from your own wallet (`--manual-xmr`). A Raspberry Pi 5 is enough: the install script fetches the arm64 Monero tools and the program cross-compiles with `GOARCH=arm64`. Mine ran on one, with its Monero and Base traffic sent through Tor so the machine's IP stayed off the node it talked to.

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

To sell for a token instead (v1 only — v2 is ETH only), add `--asset` with the token's address and quote `--min`, `--max` and `--price` in that token:

```bash
  --asset 0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913 --min 20 --max 500 --price 0.0066   # USDC on Base
```

Buyers of a token offer approve the contract once, then pay. The fee comes out in the same token.

To see who else is selling, list the open offers — every one carries its seller's address:

```bash
./bin/monero-swap --env mainnet --eth-rpc https://mainnet.base.org --contract 0x67fe8681563F37f2A8BBed84C85784a678FeC693   --monerod-host node.monerodevs.org --monerod-port 18089 offers
```

## Deployments

| Chain | Contract | Version | Fee |
|---|---|---|---|
| **Base** (mainnet) | `0xC2b2e8D385309d6552657c0b80434ca616DE12fC` | **v2 — what the page uses** | 0.15%, split 50/50 |
| Base (mainnet) | `0x67fe8681563F37f2A8BBed84C85784a678FeC693` | v1 — superseded | 0.15% to one wallet |
| Sepolia (test) | `0xB96bDd5834F455C1A6edA15e5bAF25eFd506d61E` | v1 | 0.15% to a throwaway wallet |
| Base Sepolia (test) | `0x97f8A483cFa8680F67aC24D83bbe4Fbc4f250755` | v1 | 0.15% to a throwaway wallet |

v2's two fee addresses, both Safe 1.4.1 multisigs on Base, both fixed at deploy:

- founder — `0xa3B3dcA37aE4deB5B20369E25DA8756e829e2287`
- dev fund — `0x61C9cc608Edf3Ba392B8c654171823984FA32240`

Neither version has an owner and neither can be changed. v2's source is verified on [Basescan](https://basescan.org/address/0xc2b2e8d385309d6552657c0b80434ca616de12fc#code), so what you read there is what runs. v1 is still deployed and still works — it just isn't what the page talks to any more. The test contracts are v1 and pair with Monero stagenet.

## Run a relay

Browsers can't call public Monero nodes directly, so the page reads the chain through a relay that forwards a short list of read-only calls. It holds nothing and sees only block data. The page comes prefilled with mine — `monero-relay.dafarusd.workers.dev` — and you can swap in your own:

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

- **New contract.** `contracts-v2/XmrSwap.sol` replaces their `SwapCreator.sol`. The offer board moved on-chain, so there's no peer-to-peer network and no bootnodes to die — theirs were all dead by the time I tried. A fee, fixed at deploy, is how this one pays for itself. The claim checks the revealed secret against the Monero key directly with an on-chain ed25519 multiply, which is cheap on Base; that removes the secp256k1 side and the cross-curve DLEq proof, the heaviest part of their client. One-time keys can't be reused across swaps. A seller's payout goes to any wallet, never the gas key. `XmrSwapV2.sol` is the one that's live now: the offer board is enumerable straight from storage, so finding an offer never depends on how long a node keeps its logs; the fee splits between two multisigs and is credited rather than sent, so no fee address can block a claim; timeouts start at 24 hours; sellers post a bond that's always returned; offers can price off a Chainlink-style feed instead of a fixed number; and tokens are gone.
- **New seller and buyer programs.** `swap2/` and `cmd/monero-swap/` are new. Their `swapd` is still in the tree but nothing here runs it. A failed call retries instead of ending the swap — a node outage killed their buyer's watcher mid-swap in my first run, and only their restart recovery saved it.
- **A browser buyer.** `web/` is new. Their UI was unmaintained; this one needs only MetaMask.
- **A relay.** `relay/` is new, because browsers can't call Monero nodes.
- **Kept their license** for the Go code (LGPL-3.0). The contract, page and relay are MIT. The ed25519 library is Jan Vornberger's, MIT, turned into a Solidity library by hbs in 2025 for MoneroSwap. The header of `contracts-v2/src/lib/Ed25519.sol` says so.

## This is not the only one

[MoneroSwap](https://codeberg.org/moneroswap/moneroswap) by hbs is an active EVM-to-Monero atomic swap. The Monero community funded it through a [CCS](https://ccs.getmonero.org/proposals/hbs-evm-atomic-swaps.html) for 135 XMR, paid out in early 2026, and it has been shown at EthCC and MoneroKon. It swaps a chain's native currency and it's deployed at the same address on Gnosis, Lens, Stable and Base. Worth being precise about, because I got it wrong publicly once: for a manual swap **both** sides are browser-only. The Monero seller uses a browser wallet plus a normal phone wallet — Cake, Monerujo — to scan a payment QR, and runs no daemon, no `monero-wallet-rpc` and no CLI. Only their market-making bot needs that. Mine needs a program on the selling side. The ed25519 library in this contract is hbs's work.

So this repo is a separate implementation, not a revival of anything abandoned. What actually differs, and it's a short list: this contract has no owner at all, where theirs has one who can change the coverage ratio, the delays and the oracle; and it funds itself with a 0.15% fee fixed in the contract, half of it going to a dev-fund multisig, rather than a grant that ends. The buyer's page also does its own Monero scanning in the browser instead of handing you a QR for a phone wallet.

Base is not on that list. They're deployed on Base too, at the same address as everywhere else, and Gnosis gas is ten wei — their swaps cost less than mine do. v1 swapped ERC-20 tokens as well as ETH; v2 dropped them, because a token issuer can freeze a contract's balance and strand every swap sitting inside it, which is the argument their FAQ made before I started. If you would rather use the community-funded one, use theirs.

If you learned something here, the people to thank are noot, dimalinux and the ChainSafe team.

## Limits, honestly

- Not audited. v2 is 525 lines, all 83 tests pass, and twelve invariants hold across 128,000 fuzzed calls a run — but nobody outside this repo has read it. Tests only prove the things I thought to check.
- The v2 seller has never completed a swap. It compiles, and every step is ported from the v1 seller that did run real swaps on test networks, but nobody has run it against a real buyer. Treat the first one as a test and keep it small.
- Nobody is running a v2 seller right now, so the board is empty. That's not a bug in the contract or the page — both work; there's just nobody selling.
- Selling on v2 ties up capital. The bond is 5% of your offer's ceiling, locked while the offer is live. You always get it back, but it's idle while it sits there.
- The buyer page can't send Monero yet. It hands you the keys; your own wallet does the sending.
- The seller's program has to be online. That's not a bug, it's Monero: only a private key can move it, and a key has to live somewhere.
- Your IP leaks to the Monero node and the Base RPC unless you route through Tor. Mine did — install `tor` and `torsocks`, run the wallet under `torsocks`, and start the seller with `HTTPS_PROXY=socks5://127.0.0.1:9050`. That hides the machine, not the swap: the ETH side is public on-chain either way.
- Offers are single-use. Each swap gets fresh keys, and the contract refuses reused ones.
- Base only, for now. Same contract works on any Ethereum-style chain; each one splits the sellers.
- One person built this, and every swap was replayed on test networks first. The test contracts are listed above and every swap is on those chains for anyone to read. Don't take my word for it.

---

Built by Dafarus — local-first software and hardware you own.

Follow the work on X: [@Dafarusd](https://x.com/Dafarusd)

My companies:
- Steel Valley Burners — [Facebook](https://www.facebook.com/steelvalleyburners)
- Keephaven — [keephaven.co](https://keephaven.co) · [X](https://x.com/Keephaven) · [Facebook](https://www.facebook.com/profile.php?id=61592155452190)

More work: [gate](https://github.com/dafarusd/gate) · [Sentinel](https://github.com/dafarusd/sentinel-public) · [Agent Ultra](https://github.com/dafarusd/Ultra-Agent-Release) · [EveryVoice](https://github.com/dafarusd/everyvoice) · [Mind Meld](https://github.com/dafarusd/mindmeld) · [monero-swap](https://github.com/dafarusd/monero-swap)
