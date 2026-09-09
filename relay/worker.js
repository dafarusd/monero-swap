// Cloudflare Worker: lets a browser page read a Monero node.
// Public monerod nodes do not send the CORS header browsers require, so the page
// cannot call them directly. This relay forwards a short allow-list of read-only
// calls and adds the header. It keeps no state. Anyone can deploy their own:
//   npx wrangler deploy   (set UPSTREAM to the monerod RPC base URL)
const ALLOWED_PATHS = ['/json_rpc', '/get_transactions', '/get_transaction_pool'];
const ALLOWED_METHODS = ['get_block_count', 'get_block', 'get_info', 'get_block_header_by_height'];
const CORS = { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Headers': 'Content-Type', 'Access-Control-Allow-Methods': 'POST, OPTIONS' };

export default {
  async fetch(req, env) {
    if (req.method === 'OPTIONS') return new Response(null, { headers: CORS });
    const path = new URL(req.url).pathname;
    if (req.method !== 'POST' || !ALLOWED_PATHS.includes(path)) return new Response('not allowed', { status: 403, headers: CORS });
    const body = await req.text();
    if (path === '/json_rpc') {
      let method;
      try { method = JSON.parse(body).method; } catch { return new Response('bad json', { status: 400, headers: CORS }); }
      if (!ALLOWED_METHODS.includes(method)) return new Response('method not allowed', { status: 403, headers: CORS });
    }
    const upstream = (env.UPSTREAM || 'http://xmr-lux.boldsuck.org:38081').replace(/\/$/, '');
    const r = await fetch(upstream + path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body });
    return new Response(await r.text(), { status: r.status, headers: { ...CORS, 'Content-Type': 'application/json' } });
  },
};
