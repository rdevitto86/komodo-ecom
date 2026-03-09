import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      out: 'build',
      precompress: true,
    }),
    paths: {
      base: '',
      assets: '',
    },
    alias: {
      $lib: 'src/lib',
      $components: 'src/lib/components',
      $api: 'src/routes/api',
    },
  }
};
