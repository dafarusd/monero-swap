import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import * as m from '../src/monero.js';

const vectors = JSON.parse(fs.readFileSync(new URL('./fixtures/vectors.json', import.meta.url)));
const locktx = JSON.parse(fs.readFileSync(new URL('./fixtures/locktx.json', import.meta.url)));

test('key derivation matches the Go key code for every vector', () => {
  for (const p of vectors.pairs) {
    const spendPriv = m.unhex(p.spend_priv);
    assert.equal(m.hex(m.pubFromPriv(spendPriv)), p.spend_pub);
    assert.equal(m.hex(m.viewFromSpend(spendPriv)), p.view_priv);
    assert.equal(m.hex(m.pubFromPriv(m.unhex(p.view_priv))), p.view_pub);
  }
});

test('addresses encode and decode like Go on both networks', () => {
  for (const p of vectors.pairs) {
    const sp = m.unhex(p.spend_pub), vp = m.unhex(p.view_pub);
    assert.equal(m.address(m.NET.stagenet, sp, vp), p.address_stagenet);
    assert.equal(m.address(m.NET.mainnet, sp, vp), p.address_mainnet);
    const d = m.decodeAddress(p.address_stagenet);
    assert.equal(d.net, m.NET.stagenet);
    assert.equal(m.hex(d.spendPub), p.spend_pub);
    assert.equal(m.hex(d.viewPub), p.view_pub);
  }
  assert.throws(() => m.decodeAddress(vectors.pairs[0].address_stagenet.slice(0, -1) + '1'));
});

test('shared address of two parties matches Go', () => {
  const [a, b] = vectors.pairs;
  const s = vectors.shared01;
  const r = m.sharedAddress(m.NET.stagenet, m.unhex(a.spend_pub), m.unhex(b.spend_pub), m.unhex(a.view_priv), m.unhex(b.view_priv));
  assert.equal(m.hex(r.spendPub), s.spend_pub);
  assert.equal(m.hex(r.viewPriv), s.view_priv);
  assert.equal(r.address, s.address_stagenet);
  assert.equal(m.hex(m.scalarAdd(m.unhex(a.spend_priv), m.unhex(b.spend_priv))), s.spend_priv);
  // the summed spend key really controls the summed public key
  assert.equal(m.hex(m.pubFromPriv(m.unhex(s.spend_priv))), s.spend_pub);
});

test('generateKeys produces a consistent, Monero-style pair', () => {
  const k = m.generateKeys();
  assert.equal(m.hex(m.pubFromPriv(k.spendPriv)), m.hex(k.spendPub));
  assert.equal(m.hex(m.viewFromSpend(k.spendPriv)), m.hex(k.viewPriv));
  assert.notEqual(m.hex(m.generateKeys().spendPriv), m.hex(k.spendPriv));
});

test('scanner finds the real 0.02 XMR lock from the last swap and verifies its commitment', () => {
  const viewAB = m.scalarAdd(m.unhex(locktx.view_priv_a), m.unhex(locktx.view_priv_b));
  const { spendPub } = m.decodeAddress(locktx.address);
  const found = m.scanTx(locktx.as_json, viewAB, spendPub);
  assert.equal(found.length, 1, 'exactly one output to the shared address');
  assert.equal(found[0].amount, BigInt(locktx.expected_piconero));
  assert.equal(found[0].verified, true);
  assert.equal(m.fmtXMR(found[0].amount), '0.020000000000');
});

test('scanner finds nothing with the wrong view key', () => {
  const { spendPub } = m.decodeAddress(locktx.address);
  assert.equal(m.scanTx(locktx.as_json, m.unhex(locktx.view_priv_a), spendPub).length, 0);
});
