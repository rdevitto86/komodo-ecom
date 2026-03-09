import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  plugins: [
    tailwindcss(), 
    sveltekit()
  ],
  assetsInclude: ['**/*.glb', '**/*.gltf'],
  server: {
    port: Number(process.env.PORT) || 7001,
    hmr: {
      overlay: true
    }
  }
});
