// Ethereum side of the buyer page: read offers, take one, set ready, watch the claim, refund.
// Every function takes an ethers Contract so it works with MetaMask in the browser and a
// plain provider in tests.
//
// This drives XmrSwapV2. Two things changed from v1 and both simplify this file:
// offers are enumerable straight from storage, so finding one no longer means scanning logs
// and an offer stays visible for its whole life; and the contract is ETH only, so the token
// approve dance is gone. `assetInfo` stays as an ETH-shaped stub so the page code above it
// keeps working unchanged.
import { ethers } from 'ethers';
import ABI from './XmrSwapV2.abi.json' with { type: 'json' };

export const CHAINS = {
  8453: { name: 'Base', explorer: 'https://basescan.org', xmrNet: 'mainnet', contract: '0xC2b2e8D385309d6552657c0b80434ca616DE12fC', relay: 'https://monero-relay.dafarusd.workers.dev' },
};
const LOG_CHUNK = 2000;
const OFFER_PAGE = 50;
const ETH = { address: ethers.ZeroAddress, symbol: 'ETH', decimals: 18 };

/** Kept for the page above: v2 is ETH only, so this always describes ETH. */
export async function assetInfo() { return ETH; }
export const fmtAmount = (v, decimals) => ethers.formatUnits(v, decimals);
export const parseAmount = (s, decimals) => ethers.parseUnits(s, decimals);

export async function connect() {
  if (!globalThis.ethereum) throw new Error('No wallet found. Install MetaMask or another browser wallet.');
  const provider = new ethers.BrowserProvider(globalThis.ethereum);
  await provider.send('eth_requestAccounts', []);
  const signer = await provider.getSigner();
  const net = await provider.getNetwork();
  return { provider, signer, chainId: Number(net.chainId), account: await signer.getAddress() };
}

export function contractAt(address, runner) { return new ethers.Contract(address, ABI, runner); }
export const b32 = (bytes) => '0x' + Array.from(bytes, (x) => x.toString(16).padStart(2, '0')).join('');

async function chunkedLogs(c, filter, from, to, concurrency = 5) {
  const ranges = [];
  for (let s = from; s <= to; s += LOG_CHUNK) ranges.push([s, Math.min(s + LOG_CHUNK - 1, to)]);
  const out = [];
  for (let i = 0; i < ranges.length; i += concurrency) {
    const batch = ranges.slice(i, i + concurrency).map(([s, e]) => c.queryFilter(filter, s, e));
    for (const logs of await Promise.all(batch)) out.push(...logs);
  }
  return out.sort((a, b) => a.blockNumber - b.blockNumber || a.index - b.index);
}

/** What one offer costs right now. Fixed-price offers carry their own number; feed-priced ones
 *  resolve on-chain and can refuse (stale feed, or above the seller's ceiling), in which case the
 *  offer is not takeable at the moment and we leave it out. */
async function resolvePrice(c, o) {
  if (o.xmrPerAsset !== 0n) return { price: o.xmrPerAsset, fixed: true };
  try {
    return { price: await c.priceOf(o.id), fixed: false };
  } catch {
    return null;
  }
}

/** Every live offer, read straight from the contract's own list. */
export async function listOffers(c) {
  const provider = c.runner.provider || c.runner;
  const now = (await provider.getBlock('latest')).timestamp;
  const total = Number(await c.offerCount());
  const offers = [];
  for (let offset = 0; offset < total; offset += OFFER_PAGE) {
    const page = await c.listOffers(offset, OFFER_PAGE);
    for (const o of page) {
      if (!o.active || Number(o.expiry) <= now) continue;
      const p = await resolvePrice(c, o);
      if (!p) continue;
      offers.push({
        id: o.id, maker: o.maker, payout: o.payout, min: o.minAmount, max: o.maxAmount,
        xmrPerAsset: p.price, fixedPrice: p.fixed, asset: ETH,
        expiry: Number(o.expiry), t1: Number(o.timeout1Duration), t2: Number(o.timeout2Duration),
        makerSpendPub: o.makerSpendPub, makerViewPriv: o.makerViewPriv, bond: o.bond,
      });
    }
  }
  return offers.sort((a, b) => Number(b.xmrPerAsset - a.xmrPerAsset)); // best price first
}

export function xmrFor(amountWei, xmrPerAsset) { return (BigInt(amountWei) * BigInt(xmrPerAsset)) / 10n ** 18n; }

/** Locks ETH against an offer. Returns the SwapCreated details.
 *  `minXmrPiconero` is the floor the contract enforces: for a fixed-price offer that is exactly the
 *  quote, and for a feed-priced one we allow 1% of drift between quoting and mining rather than
 *  letting the seller's feed move the price out from under the buyer. */
export async function takeOffer(c, offerId, amountWei, spendPub, viewPriv, asset, onStatus = () => {}) {
  const o = await c.getOffer(offerId);
  const p = await resolvePrice(c, o);
  if (!p) throw new Error('That offer cannot be priced right now — its price feed is stale or above the seller’s limit.');
  const quoted = xmrFor(amountWei, p.price);
  const floor = p.fixed ? quoted : (quoted * 99n) / 100n;
  onStatus('Confirm the payment in your wallet…');

  const tx = await c.takeOffer(offerId, b32(spendPub), b32(viewPriv), floor, { value: amountWei });
  const receipt = await tx.wait();
  for (const l of receipt.logs) {
    let parsed;
    try { parsed = c.interface.parseLog(l); } catch { continue; }
    if (parsed && parsed.name === 'SwapCreated') {
      return {
        swapId: parsed.args.swapId, xmrPiconero: parsed.args.xmrPiconero,
        timeout1: Number(parsed.args.timeout1), timeout2: Number(parsed.args.timeout2),
        block: receipt.blockNumber, txHash: receipt.hash,
      };
    }
  }
  throw new Error('takeOffer succeeded but no SwapCreated event was found');
}

export async function getSwap(c, swapId) {
  const s = await c.swaps(swapId);
  return { stage: Number(s.stage), taker: s.taker, maker: s.maker, value: s.value, timeout1: Number(s.timeout1), timeout2: Number(s.timeout2) };
}

export async function setReady(c, swapId) { const tx = await c.setReady(swapId); return (await tx.wait()).hash; }
export async function refund(c, swapId, secret) { const tx = await c.refund(swapId, b32(secret)); return (await tx.wait()).hash; }

/** The maker's revealed secret (hex, no 0x) once it has claimed, else null. */
export async function findClaimed(c, swapId, fromBlock) {
  const provider = c.runner.provider || c.runner;
  const tip = await provider.getBlockNumber();
  const logs = await chunkedLogs(c, c.filters.Claimed(swapId), fromBlock, tip);
  if (logs.length === 0) return null;
  return { secret: logs[0].args.secret.slice(2), txHash: logs[0].transactionHash };
}

/** Finds a swap on-chain from the taker's one-time spend pubkey, for restoring a keys-only backup.
 *  The pubkey is unique per swap and lives in the SwapCreated event, so this works from any wallet. */
export async function findSwapByTakerSpendPub(c, spendPubHex, blocksBack = 500000, onStatus = () => {}) {
  const provider = c.runner.provider || c.runner;
  const tip = await provider.getBlockNumber();
  const from = Math.max(0, tip - blocksBack);
  const want = ('0x' + String(spendPubHex).replace(/^0x/, '')).toLowerCase();
  onStatus(`Scanning the chain for your swap (blocks ${from} to ${tip})…`);
  const logs = await chunkedLogs(c, c.filters.SwapCreated(), from, tip);
  for (const l of logs) {
    if (String(l.args.takerSpendPub).toLowerCase() === want) {
      return {
        swapId: l.args.swapId, xmrPiconero: l.args.xmrPiconero,
        timeout1: Number(l.args.timeout1), timeout2: Number(l.args.timeout2),
        block: l.blockNumber, txHash: l.transactionHash,
      };
    }
  }
  return null;
}

export async function chainTime(c) {
  const provider = c.runner.provider || c.runner;
  return (await provider.getBlock('latest')).timestamp;
}
