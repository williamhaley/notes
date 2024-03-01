import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    host: '0.0.0.0',
    port: 3000,
    hmr: {
      protocol: 'ws',
      port: 3001,
      clientPort: 3001
    }
  },
  build: {
    target: 'esnext'
  }
});
