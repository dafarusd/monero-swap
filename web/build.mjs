import { build } from 'esbuild';
import { mkdirSync, copyFileSync } from 'node:fs';
mkdirSync('dist', { recursive: true });
await build({ entryPoints: ['src/app.js'], bundle: true, minify: true, format: 'iife', target: ['es2022'], outfile: 'dist/app.js', logLevel: 'info' });
copyFileSync('index.html', 'dist/index.html');
