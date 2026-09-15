// Live read-only test against the deployed XmrSwapV2 on Base. Reads only — it never sends a
// transaction and needs no key. Run with: npm run test:live
import test from 'node:test';
import assert from 'node:assert/strict';
import { ethers } from 'ethers';
import * as e from '../src/eth.js';

const CHAIN = 8453;
const RPC = process.env.ETH_RPC || 'https://mainnet.base.org';

test('live: the deployed contract matches the terms the page promises', async () => {
  const provider = new ethers.JsonRpcProvider(RPC, CHAIN);
  const c = e.contractAt(e.CHAINS[CHAIN].contract, provider);

  // The economics are immutable, so the page can state them as fact. If any of these drift,
  // the page is describing a contract that is not the one it is talking to.
  assert.equal(await c.feeBps(), 15n, 'fee is 0.15%');
  assert.equal(await c.founderShareBps(), 5000n, 'fee splits 50/50');
  assert.equal(await c.bondBps(), 500n, 'bond is 5% of an offer ceiling');
  assert.equal(await c.MIN_TIMEOUT(), 86400n, 'timeouts are at least 24h');
  assert.equal(await c.bondFor(10n ** 18n), 5n * 10n ** 16n, 'bondFor(1 ETH) is 0.05 ETH');

  // No owner and no upgrade path: there is no admin entry point to call at all.
  assert.equal(typeof c.owner, 'undefined', 'contract exposes no owner()');
});

test('live: offers read straight from contract storage', async () => {
  const provider = new ethers.JsonRpcProvider(RPC, CHAIN);
  const c = e.contractAt(e.CHAINS[CHAIN].contract, provider);

  const count = Number(await c.offerCount());
  const offers = await e.listOffers(c);
  console.log(`offers on the board: ${count}, takeable right now: ${offers.length}`);

  // listOffers filters out expired and unpriceable offers, so it can only ever be a subset.
  assert.ok(offers.length <= count, 'listOffers returned more than the board holds');

  for (const o of offers) {
    assert.equal(o.makerSpendPub.length, 66);
    assert.equal(o.makerViewPriv.length, 66);
    assert.ok(o.max >= o.min, 'offer ceiling below its floor');
    assert.equal(o.asset.symbol, 'ETH', 'v2 is ETH only');
    assert.ok(o.xmrPerAsset > 0n, 'offer resolved to a zero price');
    assert.equal(e.xmrFor(10n ** 16n, o.xmrPerAsset), o.xmrPerAsset / 100n);
  }
});
