// ── Adapter ────────────────────────────────────────────────────────────────
// Demo  → adapter-static (S3 + CloudFront). Use: bun run build:demo
// Live  → adapter-node  (EC2 / ECS).        Swap the import below + bun run build
// ──────────────────────────────────────────────────────────────────────────
import adapter from '@sveltejs/adapter-static';
// import adapter from '@sveltejs/adapter-node';

import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    // adapter-static: SPA fallback for CloudFront 404 → index.html routing
    // adapter-node: remove the fallback option when switching
    adapter: adapter({
      fallback: 'index.html',
    }),
    alias: {
      '$lib': 'src/lib',
      '$components': 'src/lib/components',
      '$styles': 'src/lib/styles',
      '$routes': 'src/routes',
      '$static': 'static',
    },
    csrf: {
      checkOrigin: true,
    },
  },
  compilerOptions: {
    runes: true,
  },
};

export default config;
