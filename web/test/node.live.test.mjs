// Live test against a public stagenet node: finds the real lock from the last swap.
import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import * as m from '../src/monero.js';
import * as n from '../src/node.js';

const locktx = JSON.parse(fs.readFileSync(new URL('./fixtures/locktx.json', import.meta.url)));
const NODE = process.env.XMR_NODE || 'http://xmr-lux.boldsuck.org:38081';

test('live: scanning the lock block finds 0.02 XMR with a real confirmation count', async () => {
  const viewAB = m.scalarAdd(m.unhex(locktx.view_priv_a), m.unhex(locktx.view_priv_b));
  const { spendPub } = m.decodeAddress(locktx.address);
  const tip = await n.getHeight(NODE);
  const hits = await n.scanRange(NODE, locktx.block_height, locktx.block_height, viewAB, spendPub);
  assert.equal(hits.length, 1);
  assert.equal(hits[0].hash, locktx.txid);
  assert.equal(hits[0].amount, BigInt(locktx.expected_piconero));
  const s = n.summarize(hits, tip);
  assert.equal(s.confirmed, BigInt(locktx.expected_piconero));
  assert.ok(tip >= locktx.block_height + 9);
});
