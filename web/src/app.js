// The buyer page. Vanilla JS. State lives in localStorage; keys never leave this browser
// except the one-time public spend key and private view key posted to the contract.
import { ethers } from 'ethers';
import * as m from './monero.js';
import * as n from './node.js';
import * as e from './eth.js';

const $ = (id) => document.getElementById(id);
const LS_SWAPS = 'xmrswap.swaps';
const LS_SETTINGS = 'xmrswap.settings';
const CONFS = 10;

const state = {
  wallet: null,          // {provider, signer, chainId, account}
  chain: null,           // CHAINS entry
  contract: null,        // ethers.Contract (signer when connected, else read-only)
  readProvider: null,
  offers: [],
  swaps: load(LS_SWAPS, {}),
  settings: load(LS_SETTINGS, { node: 'http://127.0.0.1:8089', rpc: { 11155111: 'https://ethereum-sepolia-rpc.publicnode.com', 8453: 'https://mainnet.base.org' }, chainId: 11155111 }),
  timers: {},
};

function load(k, d) { try { return JSON.parse(localStorage.getItem(k)) ?? d; } catch { return d; } }
function save() { localStorage.setItem(LS_SWAPS, JSON.stringify(state.swaps)); localStorage.setItem(LS_SETTINGS, JSON.stringify(state.settings)); }
const fmtETH = (wei) => ethers.formatEther(wei);
const when = (ts) => new Date(ts * 1000).toLocaleString();
const short = (s) => s.slice(0, 10) + '…';
function say(msg, kind = '') { const el = $('msg'); el.textContent = msg; el.className = 'msg ' + kind; }

// ---------------- setup

async function init() {
  $('nodeUrl').value = state.settings.node;
  $('chainSel').value = String(state.settings.chainId);
  $('connectBtn').onclick = connectWallet;
  $('refreshBtn').onclick = refreshOffers;
  $('nodeUrl').onchange = () => { state.settings.node = $('nodeUrl').value.trim(); save(); };
  $('chainSel').onchange = () => { state.settings.chainId = Number($('chainSel').value); save(); setupChain(); refreshOffers(); };
  setupChain();
  renderSwaps();
  await reconnectSilently();
  await refreshOffers();
  for (const id of Object.keys(state.swaps)) resumeSwap(id);
}

// If the wallet already approved this page, pick it up again without a prompt (after a reload).
async function reconnectSilently() {
  if (!globalThis.ethereum) return;
  try {
    const accounts = await globalThis.ethereum.request({ method: 'eth_accounts' });
    if (accounts && accounts.length) await connectWallet();
  } catch { /* fine, the button still works */ }
}

function setupChain() {
  const chainId = state.wallet ? state.wallet.chainId : state.settings.chainId;
  state.chain = e.CHAINS[chainId];
  if (!state.chain || !state.chain.contract) { say(`No contract configured for chain ${chainId}.`, 'bad'); state.contract = null; return; }
  state.readProvider = new ethers.JsonRpcProvider(state.settings.rpc[chainId] || '', chainId);
  const runner = state.wallet ? state.wallet.signer : state.readProvider;
  state.contract = e.contractAt(state.chain.contract, runner);
  $('chainName').textContent = state.chain.name;
  $('contractLink').textContent = state.chain.contract;
  $('contractLink').href = `${state.chain.explorer}/address/${state.chain.contract}#code`;
}

async function connectWallet() {
  try {
    state.wallet = await e.connect();
    $('account').textContent = state.wallet.account;
    $('connectBtn').textContent = 'Connected';
    if (!e.CHAINS[state.wallet.chainId]) { say(`Your wallet is on chain ${state.wallet.chainId}. Switch it to ${e.CHAINS[state.settings.chainId].name}.`, 'bad'); return; }
    state.settings.chainId = state.wallet.chainId; save();
    $('chainSel').value = String(state.wallet.chainId);
    setupChain();
    say('Wallet connected.');
    await refreshOffers();
  } catch (err) { say(err.message, 'bad'); }
}

// ---------------- offers

async function refreshOffers() {
  if (!state.contract) return;
  $('offers').innerHTML = '<p>Loading offers…</p>';
  try {
    state.offers = await e.listOffers(e.contractAt(state.chain.contract, state.readProvider), 10000);
  } catch (err) { $('offers').innerHTML = `<p class="bad">Could not read offers: ${err.message}</p>`; return; }
  if (state.offers.length === 0) { $('offers').innerHTML = '<p>No open offers right now.</p>'; return; }
  $('offers').innerHTML = '';
  for (const o of state.offers) {
    const card = document.createElement('div');
    card.className = 'card';
    const perEth = m.fmtXMR(o.xmrPerAsset);
    card.innerHTML = `
      <div class="row"><b>${perEth} XMR per ETH</b><span class="dim">offer ${short(o.id)}</span></div>
      <div>Pay between ${fmtETH(o.min)} and ${fmtETH(o.max)} ETH. Open until ${when(o.expiry)}.</div>
      <div class="dim">Seller has ${Math.round(o.t1 / 60)} min to lock the Monero after you pay; if it never shows, your ETH comes back.</div>
      <div class="row">
        <label>ETH to pay <input type="number" step="0.001" min="${fmtETH(o.min)}" max="${fmtETH(o.max)}" value="${fmtETH(o.min)}"></label>
        <span class="get"></span>
        <button>Take this offer</button>
      </div>`;
    const input = card.querySelector('input'), get = card.querySelector('.get'), btn = card.querySelector('button');
    const upd = () => { try { get.textContent = `you get ${m.fmtXMR(e.xmrFor(ethers.parseEther(input.value || '0'), o.xmrPerAsset))} XMR`; } catch { get.textContent = ''; } };
    input.oninput = upd; upd();
    btn.onclick = () => takeOffer(o, input.value);
    $('offers').appendChild(card);
  }
}

// ---------------- take

async function takeOffer(o, ethStr) {
  if (!state.wallet) { say('Connect your wallet first.', 'bad'); return; }
  if (!state.settings.node) { say('Set a Monero node relay URL first (top of page).', 'bad'); return; }
  let amount;
  try { amount = ethers.parseEther(ethStr); } catch { say('Bad amount.', 'bad'); return; }
  if (amount < o.min || amount > o.max) { say('Amount is outside the offer range.', 'bad'); return; }

  // Height first, so the scan later starts before the lock.
  let restoreHeight;
  try { restoreHeight = Math.max(0, (await n.getHeight(state.settings.node)) - 20); }
  catch (err) { say(`Cannot reach the Monero node relay at ${state.settings.node}: ${err.message}`, 'bad'); return; }

  const k = m.generateKeys();
  const netId = m.NET[state.chain.xmrNet];
  const shared = m.sharedAddress(netId, m.unhex(o.makerSpendPub), k.spendPub, m.unhex(o.makerViewPriv), k.viewPriv);
  const rec = {
    state: 'pending_take', chainId: state.wallet.chainId, offerId: o.id, xmrNet: state.chain.xmrNet,
    mySpendPriv: m.hex(k.spendPriv), myViewPriv: m.hex(k.viewPriv), mySpendPub: m.hex(k.spendPub),
    theirSpendPub: o.makerSpendPub.slice(2), theirViewPriv: o.makerViewPriv.slice(2),
    sharedAddress: shared.address, restoreHeight, amountWei: amount.toString(), xmrPiconero: e.xmrFor(amount, o.xmrPerAsset).toString(),
    created: Date.now(), lastScanned: restoreHeight - 1, hits: [],
  };
  // Keys are saved before anything is sent. Losing them after paying means losing the Monero.
  const tempId = 'pending:' + rec.mySpendPub;
  state.swaps[tempId] = rec; save();
  offerBackup(tempId, rec);

  try {
    say('Confirm the transaction in your wallet…');
    const r = await e.takeOffer(state.contract, o.id, amount, k.spendPub, k.viewPriv);
    delete state.swaps[tempId];
    Object.assign(rec, { state: 'eth_locked', swapId: r.swapId, timeout1: r.timeout1, timeout2: r.timeout2, createdBlock: r.block, takeTx: r.txHash, xmrPiconero: r.xmrPiconero.toString() });
    state.swaps[r.swapId] = rec; save();
    say('ETH locked. Now watching for the Monero.');
    renderSwaps();
    resumeSwap(r.swapId);
  } catch (err) {
    delete state.swaps[tempId]; save();
    say(`Not taken: ${err.shortMessage || err.message}`, 'bad');
  }
}

function offerBackup(id, rec) {
  const blob = new Blob([JSON.stringify({ id, ...rec }, null, 2)], { type: 'application/json' });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = `xmrswap-keys-${id.slice(0, 12)}.json`;
  a.textContent = 'Download a backup of this swap’s keys';
  a.className = 'backup';
  $('backups').prepend(a);
}

// ---------------- running swaps

function renderSwaps() {
  const box = $('swaps');
  box.innerHTML = '';
  const ids = Object.keys(state.swaps).filter((k) => !k.startsWith('pending:'));
  if (ids.length === 0) { box.innerHTML = '<p class="dim">No swaps yet.</p>'; return; }
  for (const id of ids.reverse()) {
    const s = state.swaps[id];
    const card = document.createElement('div');
    card.className = 'card'; card.id = 'swap-' + id;
    box.appendChild(card);
    renderSwap(id);
  }
}

function renderSwap(id) {
  const s = state.swaps[id];
  const card = $('swap-' + id);
  if (!card) return;
  const lines = [`<div class="row"><b>${fmtETH(s.amountWei)} ETH → ${m.fmtXMR(s.xmrPiconero)} XMR</b><span class="dim">swap ${short(id)}</span></div>`];
  const steps = { eth_locked: 1, xmr_seen: 2, ready: 3, claimed: 4, done: 5, refunded: 5 };
  const names = ['Paid ETH into the contract', 'Monero arrived and confirmed', 'Contract set to ready', 'Seller collected the ETH and revealed its secret', s.state === 'refunded' ? 'ETH refunded to you' : 'Monero is yours'];
  const cur = steps[s.state] || 0;
  lines.push('<ol class="steps">' + names.map((t, i) => `<li class="${i < cur ? 'done' : i === cur ? 'now' : ''}">${t}</li>`).join('') + '</ol>');
  if (s.note) lines.push(`<div class="note">${s.note}</div>`);
  if (s.state === 'eth_locked') lines.push(`<div class="dim">Watching Monero address ${s.sharedAddress}<br>Seller's deadline ${when(s.timeout1)}. If nothing arrives by ${when(s.timeout1 - 600)}, you can take your ETH back.</div>`);
  if (s.state === 'claimed' || s.state === 'done') {
    lines.push(`<div class="keys"><b>Your Monero: ${m.fmtXMR(s.xmrPiconero)} XMR at</b><br><code>${s.sharedAddress}</code><br>
      Restore this as a wallet <i>from keys</i> in Feather, Cake, or the Monero GUI, then send it wherever you like. Nobody else can spend it.<br>
      private spend key <code>${s.finalSpendPriv}</code><br>private view key <code>${s.finalViewPriv}</code><br>restore height <code>${s.restoreHeight}</code></div>`);
    if (s.state === 'claimed') lines.push(`<button data-act="done">I have moved the Monero</button>`);
  }
  if (s.state === 'eth_locked' || s.state === 'ready') lines.push(`<button data-act="refund" class="warn">Take my ETH back</button>`);
  card.innerHTML = lines.join('');
  card.querySelectorAll('button').forEach((b) => { b.onclick = () => action(id, b.dataset.act); });
}

function setSwap(id, patch) { Object.assign(state.swaps[id], patch); save(); renderSwap(id); }

async function action(id, act) {
  const s = state.swaps[id];
  if (act === 'done') { setSwap(id, { state: 'done', note: '' }); return; }
  if (act === 'refund') {
    if (!state.wallet) { say('Connect your wallet first.', 'bad'); return; }
    if (!confirm('Take the ETH back? Only do this if the Monero never arrived, or the seller never collected. It cannot be undone.')) return;
    try {
      say('Confirm the refund in your wallet…');
      const tx = await e.refund(state.contract, id, m.unhex(s.mySpendPriv));
      setSwap(id, { state: 'refunded', note: `Refund transaction ${tx}` });
      say('ETH refunded.');
    } catch (err) { say(`Refund not accepted: ${err.shortMessage || err.message}`, 'bad'); }
  }
}

function resumeSwap(id) {
  if (id.startsWith('pending:')) return;
  if (state.timers[id]) return;
  const tick = async () => {
    try { await stepSwap(id); } catch (err) { setSwap(id, { note: `Waiting: ${err.shortMessage || err.message}` }); }
    const s = state.swaps[id];
    if (s && (s.state === 'eth_locked' || s.state === 'xmr_seen' || s.state === 'ready')) state.timers[id] = setTimeout(tick, 30000);
    else delete state.timers[id];
  };
  tick();
}

async function stepSwap(id) {
  const s = state.swaps[id];
  const c = state.contract;
  if (s.state === 'eth_locked') {
    const viewAB = m.scalarAdd(m.unhex(s.myViewPriv), m.unhex(s.theirViewPriv));
    const { spendPub } = m.decodeAddress(s.sharedAddress);
    const tip = await n.getHeight(state.settings.node);
    if (tip > s.lastScanned) {
      const hits = await n.scanRange(state.settings.node, s.lastScanned + 1, tip, viewAB, spendPub);
      const merged = [...s.hits.filter((h) => !h.inPool), ...hits].filter((h, i, a) => a.findIndex((x) => x.hash === h.hash) === i);
      setSwap(id, { hits: merged.map((h) => ({ ...h, amount: h.amount.toString() })), lastScanned: tip });
    }
    const sum = n.summarize(s.hits.map((h) => ({ ...h, amount: BigInt(h.amount) })), tip, CONFS);
    const need = BigInt(s.xmrPiconero);
    if (sum.confirmed >= need) { setSwap(id, { state: 'xmr_seen', note: '' }); return stepSwap(id); }
    const now = await e.chainTime(c);
    if (sum.pending + sum.confirmed >= need) setSwap(id, { note: `Monero seen, ${sum.minConf ?? 0} of ${CONFS} confirmations.` });
    else setSwap(id, { note: `No Monero yet (scanned to block ${tip}).` });
    if (now >= s.timeout1 - 600 && now < s.timeout1) setSwap(id, { note: 'The seller is out of time. Take your ETH back.' });
    return;
  }
  if (s.state === 'xmr_seen') {
    if (!state.wallet) { setSwap(id, { note: 'Monero confirmed. Connect your wallet to set the contract to ready.' }); return; }
    const now = await e.chainTime(c);
    if (now >= s.timeout1) { setSwap(id, { note: 'Too late to set ready; the seller may still collect until the second deadline, then you can refund.', state: 'ready' }); return; }
    say('Monero confirmed. Confirm "set ready" in your wallet…');
    const tx = await e.setReady(c, id);
    setSwap(id, { state: 'ready', readyTx: tx, note: '' });
    return;
  }
  if (s.state === 'ready') {
    const claimed = await e.findClaimed(c, id, s.createdBlock);
    if (claimed) {
      const finalSpend = m.scalarAdd(m.unhex(s.mySpendPriv), m.unhex(claimed.secret));
      const finalView = m.scalarAdd(m.unhex(s.myViewPriv), m.unhex(s.theirViewPriv));
      // sanity: the summed key must control the shared address
      const { spendPub } = m.decodeAddress(s.sharedAddress);
      if (m.hex(m.pubFromPriv(finalSpend)) !== m.hex(spendPub)) throw new Error('revealed secret does not match the shared address');
      setSwap(id, { state: 'claimed', finalSpendPriv: m.hex(finalSpend), finalViewPriv: m.hex(finalView), claimTx: claimed.txHash, note: '' });
      say('The seller collected. Your Monero is ready to move.');
      return;
    }
    const now = await e.chainTime(c);
    if (now >= s.timeout2) setSwap(id, { note: 'The seller never collected. You can take your ETH back now.' });
    else setSwap(id, { note: `Waiting for the seller to collect (they have until ${when(s.timeout2)}).` });
  }
}

init();
