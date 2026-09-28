import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// Backend dev (HTTP_ADDR di .env): 127.0.0.1:8080.
// API dan SPA satu origin; dev proxy /api + /health ke Go (ARCHITECTURE §7).
const BACKEND = 'http://127.0.0.1:8080';

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: {
      '/api': { target: BACKEND, changeOrigin: false },
      '/health': { target: BACKEND, changeOrigin: false },
    },
  },
});
