// Ethereum side of the buyer page: read offers, take one, set ready, watch the claim, refund.
// Every function takes an ethers Contract so it works with MetaMask in the browser and a
// plain provider in tests.
import { ethers } from 'ethers';
import ABI from './XmrSwap.abi.json' with { type: 'json' };

export const CHAINS = {
  11155111: { name: 'Sepolia (test)', explorer: 'https://sepolia.etherscan.io', xmrNet: 'stagenet', contract: '0xB96bDd5834F455C1A6edA15e5bAF25eFd506d61E' },
  84532: { name: 'Base Sepolia (test)', explorer: 'https://sepolia.basescan.org', xmrNet: 'stagenet', contract: '' },
  8453: { name: 'Base', explorer: 'https://basescan.org', xmrNet: 'mainnet', contract: '' },
};
const LOG_CHUNK = 2000;

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

/** Open ETH offers posted in the last `blocksBack` blocks. */
export async function listOffers(c, blocksBack = 10000) {
  const provider = c.runner.provider || c.runner;
  const tip = await provider.getBlockNumber();
  const now = (await provider.getBlock('latest')).timestamp;
  const from = Math.max(0, tip - blocksBack);
  const logs = await chunkedLogs(c, c.filters.OfferPosted(), from, tip);
  const offers = [];
  for (const l of logs) {
    const id = l.args.offerId;
    const o = await c.offers(id);
    if (!o.active || Number(o.expiry) <= now || o.asset !== ethers.ZeroAddress) continue;
    offers.push({
      id, maker: o.maker, payout: o.payout, min: o.minAmount, max: o.maxAmount, xmrPerAsset: o.xmrPerAsset,
      expiry: Number(o.expiry), t1: Number(o.timeout1Duration), t2: Number(o.timeout2Duration),
      makerSpendPub: o.makerSpendPub, makerViewPriv: o.makerViewPriv, block: l.blockNumber,
    });
  }
  return offers.sort((a, b) => Number(b.xmrPerAsset - a.xmrPerAsset)); // best price first
}

export function xmrFor(amountWei, xmrPerAsset) { return (BigInt(amountWei) * BigInt(xmrPerAsset)) / 10n ** 18n; }

/** Locks ETH against an offer. Returns the SwapCreated details. */
export async function takeOffer(c, offerId, amountWei, spendPub, viewPriv) {
  const tx = await c.takeOffer(offerId, amountWei, b32(spendPub), b32(viewPriv), { value: amountWei });
  const receipt = await tx.wait();
  for (const l of receipt.logs) {
    let p;
    try { p = c.interface.parseLog(l); } catch { continue; }
    if (p && p.name === 'SwapCreated') {
      return {
        swapId: p.args.swapId, xmrPiconero: p.args.xmrPiconero, timeout1: Number(p.args.timeout1), timeout2: Number(p.args.timeout2),
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

export async function chainTime(c) {
  const provider = c.runner.provider || c.runner;
  return (await provider.getBlock('latest')).timestamp;
}
