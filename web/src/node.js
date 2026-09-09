// Talks to a Monero node (directly, or through the CORS relay) and scans blocks for
// payments to a one-time address. Works in the browser and in Node 18+.
import { scanTx } from './monero.js';

async function post(url, body) {
  const r = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) });
  if (!r.ok) throw new Error(`${url}: HTTP ${r.status}`);
  return r.json();
}
async function rpc(base, method, params = {}) {
  const j = await post(`${base.replace(/\/$/, '')}/json_rpc`, { jsonrpc: '2.0', id: '0', method, params });
  if (j.error) throw new Error(`${method}: ${j.error.message}`);
  return j.result;
}

export async function getHeight(base) {
  const r = await rpc(base, 'get_block_count');
  return r.count - 1; // count = height + 1
}

/** Transaction hashes in one block (excluding the coinbase, which never pays us). */
export async function getBlockTxHashes(base, height) {
  const r = await rpc(base, 'get_block', { height });
  return r.tx_hashes || [];
}

/** Decoded transactions (monerod `as_json`) with their block heights. */
export async function getTransactions(base, hashes) {
  if (hashes.length === 0) return [];
  const out = [];
  for (let i = 0; i < hashes.length; i += 50) {
    const r = await post(`${base.replace(/\/$/, '')}/get_transactions`, { txs_hashes: hashes.slice(i, i + 50), decode_as_json: true });
    for (const t of r.txs || []) out.push({ hash: t.tx_hash, height: t.block_height, inPool: t.in_pool, tx: JSON.parse(t.as_json) });
  }
  return out;
}

/** Mempool transactions, decoded. */
export async function getPoolTransactions(base) {
  const r = await post(`${base.replace(/\/$/, '')}/get_transaction_pool`, {});
  return (r.transactions || []).map((t) => ({ hash: t.id_hash, height: null, inPool: true, tx: JSON.parse(t.tx_json) }));
}

/**
 * Scans blocks [from, to] (and the mempool) for outputs paying `spendPub` under `viewPriv`.
 * Calls onProgress(height) as it goes. Returns [{hash, height, inPool, amount, verified}].
 */
export async function scanRange(base, from, to, viewPriv, spendPub, onProgress = () => {}) {
  const hits = [];
  for (let h = from; h <= to; h++) {
    const hashes = await getBlockTxHashes(base, h);
    for (const t of await getTransactions(base, hashes)) {
      for (const o of scanTx(t.tx, viewPriv, spendPub)) hits.push({ hash: t.hash, height: t.height, inPool: false, amount: o.amount, verified: o.verified });
    }
    onProgress(h);
  }
  for (const t of await getPoolTransactions(base)) {
    for (const o of scanTx(t.tx, viewPriv, spendPub)) hits.push({ hash: t.hash, height: null, inPool: true, amount: o.amount, verified: o.verified });
  }
  return hits;
}

/** Sums verified, confirmed hits and reports how many confirmations the youngest has. */
export function summarize(hits, tipHeight, needed = 10) {
  let confirmed = 0n, pending = 0n, minConf = Infinity;
  for (const h of hits) {
    if (!h.verified) continue;
    if (h.inPool || h.height === null) { pending += h.amount; minConf = 0; continue; }
    const conf = tipHeight - h.height + 1;
    if (conf >= needed) confirmed += h.amount; else { pending += h.amount; minConf = Math.min(minConf, conf); }
  }
  return { confirmed, pending, minConf: minConf === Infinity ? null : minConf };
}
