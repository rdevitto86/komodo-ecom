import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import { PORT as PORT_KEY } from './src/lib/config';

export default defineConfig({
  plugins: [
    tailwindcss(), 
    sveltekit()
  ],
  assetsInclude: ['**/*.glb', '**/*.gltf'],
  server: {
    port: Number(process.env[PORT_KEY]) || 7001,
    hmr: {
      overlay: true
    }
  }
});
