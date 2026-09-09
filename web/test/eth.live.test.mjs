// Live read-only test against Sepolia: the deployed contract lists the seller's open offer.
import test from 'node:test';
import assert from 'node:assert/strict';
import { ethers } from 'ethers';
import * as e from '../src/eth.js';

const RPC = process.env.ETH_RPC || 'https://ethereum-sepolia-rpc.publicnode.com';

test('live: lists the open test offer on Sepolia with the maker keys readable', async () => {
  const provider = new ethers.JsonRpcProvider(RPC);
  const c = e.contractAt(e.CHAINS[11155111].contract, provider);
  const offers = await e.listOffers(c, 3000);
  assert.ok(offers.length >= 1, 'at least one open offer');
  const o = offers[0];
  assert.equal(o.makerSpendPub.length, 66);
  assert.equal(o.makerViewPriv.length, 66);
  assert.ok(o.max >= o.min);
  assert.equal(e.xmrFor(10n ** 16n, o.xmrPerAsset), o.xmrPerAsset / 100n); // 0.01 ETH
  const s = await e.getSwap(c, '0x38f35a5295034bb9c687d7f483a0769180b0e17fac9375b1bdff027b8ae5ccdd');
  assert.equal(s.stage, 3, 'last swap is COMPLETED');
  const claimed = await e.findClaimed(c, '0x38f35a5295034bb9c687d7f483a0769180b0e17fac9375b1bdff027b8ae5ccdd', 11664160);
  assert.equal(claimed.secret.length, 64);
});
