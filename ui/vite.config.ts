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
    port: 3000,
    hmr: {
      overlay: true
    }
  }
});
