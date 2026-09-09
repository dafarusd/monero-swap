// Same relay as worker.js, as a plain Node server for local testing:
//   UPSTREAM=http://xmr-lux.boldsuck.org:38081 node relay/local.mjs   (listens on 8089)
import http from 'node:http';
const UPSTREAM = (process.env.UPSTREAM || 'http://xmr-lux.boldsuck.org:38081').replace(/\/$/, '');
const PORT = Number(process.env.PORT || 8089);
const ALLOWED_PATHS = ['/json_rpc', '/get_transactions', '/get_transaction_pool'];
const ALLOWED_METHODS = ['get_block_count', 'get_block', 'get_info', 'get_block_header_by_height'];
const CORS = { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Headers': 'Content-Type', 'Access-Control-Allow-Methods': 'POST, OPTIONS' };

http.createServer(async (req, res) => {
  if (req.method === 'OPTIONS') { res.writeHead(204, CORS); return res.end(); }
  const path = new URL(req.url, 'http://x').pathname;
  if (req.method !== 'POST' || !ALLOWED_PATHS.includes(path)) { res.writeHead(403, CORS); return res.end('not allowed'); }
  let body = '';
  for await (const chunk of req) body += chunk;
  if (path === '/json_rpc') {
    let method;
    try { method = JSON.parse(body).method; } catch { res.writeHead(400, CORS); return res.end('bad json'); }
    if (!ALLOWED_METHODS.includes(method)) { res.writeHead(403, CORS); return res.end('method not allowed'); }
  }
  try {
    const r = await fetch(UPSTREAM + path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body });
    res.writeHead(r.status, { ...CORS, 'Content-Type': 'application/json' });
    res.end(await r.text());
  } catch (e) {
    res.writeHead(502, CORS); res.end(String(e));
  }
}).listen(PORT, '127.0.0.1', () => console.log(`monero relay on http://127.0.0.1:${PORT} -> ${UPSTREAM}`));
