// Monero key math and output scanning for the buyer page. Pure functions, no I/O.
// Byte order: every 32-byte key is Monero's little-endian encoding, same as the
// Go code and the contract.
import * as ed from '@noble/ed25519';
import { keccak_256 } from '@noble/hashes/sha3';
import { sha512 } from '@noble/hashes/sha512';

ed.etc.sha512Sync = (...m) => sha512(ed.etc.concatBytes(...m));

const L = ed.CURVE.n; // group order
const Point = ed.ExtendedPoint;
// Monero's second generator H (pedersen commitments)
const H_HEX = '8b655970153799af2aeadc9ff1add0ea6c7251d54154cfa92c173a0dd39c1f94';

export const NET = { mainnet: 18, stagenet: 24, testnet: 53 };

// ---- bytes / hex helpers
export const hex = (b) => Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
export const unhex = (s) => {
  s = s.replace(/^0x/, '');
  if (s.length % 2) throw new Error('odd hex');
  return Uint8Array.from(s.match(/../g).map((h) => parseInt(h, 16)));
};
export const concat = (...arrs) => {
  const out = new Uint8Array(arrs.reduce((n, a) => n + a.length, 0));
  let o = 0;
  for (const a of arrs) { out.set(a, o); o += a.length; }
  return out;
};
const te = new TextEncoder();

export function bigFromLE(b) {
  let n = 0n;
  for (let i = b.length - 1; i >= 0; i--) n = (n << 8n) | BigInt(b[i]);
  return n;
}
export function bigToLE(n, len = 32) {
  const out = new Uint8Array(len);
  for (let i = 0; i < len; i++) { out[i] = Number(n & 0xffn); n >>= 8n; }
  return out;
}
const mod = (a, m) => ((a % m) + m) % m;

// ---- scalars and points
export function scReduce32(b) { return bigToLE(mod(bigFromLE(b), L)); }
export function scalarAdd(a, b) { return bigToLE(mod(bigFromLE(a) + bigFromLE(b), L)); }
export function randomScalar() {
  const r = new Uint8Array(32);
  crypto.getRandomValues(r);
  return scReduce32(r);
}
export function pubFromPriv(priv) {
  const s = bigFromLE(priv);
  if (s === 0n || s >= L) throw new Error('scalar out of range');
  return Point.BASE.multiply(s).toRawBytes();
}
export function pointAdd(a, b) { return Point.fromHex(hex(a)).add(Point.fromHex(hex(b))).toRawBytes(); }
export function viewFromSpend(spendPriv) { return scReduce32(keccak_256(spendPriv)); }

/** Fresh one-time key pair the Monero way: random spend, view = H(spend) mod l. */
export function generateKeys() {
  const spendPriv = randomScalar();
  const viewPriv = viewFromSpend(spendPriv);
  return { spendPriv, viewPriv, spendPub: pubFromPriv(spendPriv), viewPub: pubFromPriv(viewPriv) };
}

// ---- Monero base58 (8-byte blocks -> 11 chars)
const ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';
const BLOCK_SIZES = [0, 2, 3, 5, 6, 7, 9, 10, 11];
function encodeBlock(bytes) {
  let n = 0n;
  for (const b of bytes) n = (n << 8n) | BigInt(b);
  const size = BLOCK_SIZES[bytes.length];
  let out = '';
  for (let i = 0; i < size; i++) { out = ALPHABET[Number(n % 58n)] + out; n /= 58n; }
  return out;
}
function decodeBlock(str, byteLen) {
  let n = 0n;
  for (const c of str) {
    const v = ALPHABET.indexOf(c);
    if (v < 0) throw new Error('bad base58 char');
    n = n * 58n + BigInt(v);
  }
  const out = new Uint8Array(byteLen);
  for (let i = byteLen - 1; i >= 0; i--) { out[i] = Number(n & 0xffn); n >>= 8n; }
  if (n !== 0n) throw new Error('base58 overflow');
  return out;
}
export function base58Encode(bytes) {
  let out = '';
  for (let i = 0; i < bytes.length; i += 8) out += encodeBlock(bytes.subarray(i, Math.min(i + 8, bytes.length)));
  return out;
}
export function base58Decode(str) {
  const full = Math.floor(str.length / 11);
  const rem = str.length % 11;
  const remBytes = BLOCK_SIZES.indexOf(rem);
  if (rem && remBytes < 0) throw new Error('bad base58 length');
  const out = new Uint8Array(full * 8 + (rem ? remBytes : 0));
  for (let i = 0; i < full; i++) out.set(decodeBlock(str.slice(i * 11, i * 11 + 11), 8), i * 8);
  if (rem) out.set(decodeBlock(str.slice(full * 11), remBytes), full * 8);
  return out;
}

// ---- addresses
export function address(net, spendPub, viewPub) {
  const data = concat(Uint8Array.of(net), spendPub, viewPub);
  const sum = keccak_256(data).subarray(0, 4);
  return base58Encode(concat(data, sum));
}
export function decodeAddress(str) {
  const raw = base58Decode(str);
  if (raw.length !== 69) throw new Error('not a standard address');
  const data = raw.subarray(0, 65);
  const sum = keccak_256(data).subarray(0, 4);
  if (hex(sum) !== hex(raw.subarray(65))) throw new Error('bad address checksum');
  return { net: raw[0], spendPub: raw.subarray(1, 33), viewPub: raw.subarray(33, 65) };
}
export function netName(net) { return Object.keys(NET).find((k) => NET[k] === net) || `net${net}`; }

/** The one-time address both swap parties control: spend = A + B, view = a + b. */
export function sharedAddress(net, spendPubA, spendPubB, viewPrivA, viewPrivB) {
  const spendPub = pointAdd(spendPubA, spendPubB);
  const viewPriv = scalarAdd(viewPrivA, viewPrivB);
  return { address: address(net, spendPub, pubFromPriv(viewPriv)), spendPub, viewPriv };
}

// ---- transaction scanning
function varint(n) {
  const out = [];
  while (n >= 0x80) { out.push((n & 0x7f) | 0x80); n >>>= 7; }
  out.push(n);
  return Uint8Array.from(out);
}
function readVarint(b, i) {
  let n = 0, shift = 0;
  for (;;) {
    const x = b[i++];
    n |= (x & 0x7f) << shift;
    if (!(x & 0x80)) break;
    shift += 7;
  }
  return [n, i];
}

/** Transaction public keys from tx_extra (main key, then any additional keys). */
export function parseExtra(extra) {
  const keys = [];
  const additional = [];
  let i = 0;
  while (i < extra.length) {
    const tag = extra[i++];
    if (tag === 0x00) { while (i < extra.length && extra[i] === 0) i++; continue; } // padding
    if (tag === 0x01) { keys.push(extra.subarray(i, i + 32)); i += 32; continue; }
    if (tag === 0x02) { const [n, j] = readVarint(extra, i); i = j + n; continue; } // nonce / payment id
    if (tag === 0x03) { const [, j] = readVarint(extra, i); i = j + 32; continue; } // merge mining
    if (tag === 0x04) { const [n, j] = readVarint(extra, i); i = j; for (let k = 0; k < n; k++) { additional.push(extra.subarray(i, i + 32)); i += 32; } continue; }
    break; // unknown tag: stop parsing
  }
  return { txPub: keys[0] || null, additional };
}

/** D = 8 * a * R, compressed. */
export function keyDerivation(viewPriv, txPub) {
  // Monero clears the cofactor: D = 8 * (a * R). Three doublings is exactly x8.
  return Point.fromHex(hex(txPub)).multiply(bigFromLE(viewPriv)).double().double().double().toRawBytes();
}
export function derivationToScalar(D, index) { return scReduce32(keccak_256(concat(D, varint(index)))); }
export function derivePublicKey(D, index, spendPub) {
  const s = bigFromLE(derivationToScalar(D, index));
  return Point.BASE.multiply(s).add(Point.fromHex(hex(spendPub))).toRawBytes();
}
export function viewTag(D, index) { return keccak_256(concat(te.encode('view_tag'), D, varint(index)))[0]; }

function decryptAmount(sharedSec, encAmountHex) {
  const key = keccak_256(concat(te.encode('amount'), sharedSec)).subarray(0, 8);
  const enc = unhex(encAmountHex);
  const out = new Uint8Array(8);
  for (let i = 0; i < 8; i++) out[i] = enc[i] ^ key[i];
  return bigFromLE(out);
}
function commitmentMatches(sharedSec, amount, outPkHex) {
  const mask = bigFromLE(scReduce32(keccak_256(concat(te.encode('commitment_mask'), sharedSec))));
  const H = Point.fromHex(H_HEX);
  let C = Point.BASE.multiply(mask);
  if (amount > 0n) C = C.add(H.multiply(amount));
  return hex(C.toRawBytes()) === outPkHex.toLowerCase();
}

/**
 * Finds outputs in a decoded transaction (monerod `as_json`) that pay `spendPub` with
 * private view key `viewPriv`. Returns [{index, amount (BigInt piconero), verified}].
 * `verified` is true only when the decrypted amount matches the on-chain commitment.
 */
export function scanTx(tx, viewPriv, spendPub) {
  const extra = Uint8Array.from(tx.extra);
  const { txPub, additional } = parseExtra(extra);
  if (!txPub) return [];
  const found = [];
  const rct = tx.rct_signatures || {};
  tx.vout.forEach((out, i) => {
    const target = out.target.tagged_key || out.target.key && { key: out.target.key };
    if (!target) return;
    const candidates = [txPub];
    if (additional[i]) candidates.push(additional[i]);
    for (const R of candidates) {
      const D = keyDerivation(viewPriv, R);
      if (target.view_tag !== undefined && viewTag(D, i) !== parseInt(target.view_tag, 16)) continue;
      if (hex(derivePublicKey(D, i, spendPub)) !== target.key.toLowerCase()) continue;
      let amount = BigInt(out.amount || 0);
      let verified = true;
      if (rct.ecdhInfo && rct.ecdhInfo[i]) {
        const sharedSec = derivationToScalar(D, i);
        amount = decryptAmount(sharedSec, rct.ecdhInfo[i].amount);
        verified = commitmentMatches(sharedSec, amount, rct.outPk[i]);
      }
      found.push({ index: i, amount, verified });
      break;
    }
  });
  return found;
}

export const fmtXMR = (piconero) => {
  const p = BigInt(piconero);
  return `${p / 1000000000000n}.${(p % 1000000000000n).toString().padStart(12, '0')}`;
};
