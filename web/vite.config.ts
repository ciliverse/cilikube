import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'node:path'

/** True only for that package's own files — not peer-dep suffixes in pnpm paths. */
function isPkg(id: string, name: string): boolean {
  return (
    id.includes(`/node_modules/${name}/`) || id.includes(`/node_modules/.pnpm/${name}@`)
  )
}

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    target: 'es2020',
    cssCodeSplit: true,
    chunkSizeWarningLimit: 700,
    // Keep heavy lazy vendors out of entry modulepreload.
    modulePreload: {
      resolveDependencies(_filename, deps) {
        return deps.filter(
          (d) =>
            !d.includes('vendor-charts') &&
            !d.includes('vendor-xterm') &&
            !d.includes('vendor-yaml') &&
            !d.includes('recharts') &&
            !d.includes('/d3-') &&
            !d.includes('xterm'),
        )
      },
    },
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return
          if (isPkg(id, 'xterm') || isPkg(id, 'xterm-addon-fit')) return 'vendor-xterm'
          // Leave recharts / framer-motion / axios to natural chunking.
          // Isolating axios into vendor-axios broke CJS interop (runtime: "e is not a function") → black screen.
          if (isPkg(id, 'js-yaml')) return 'vendor-yaml'
          if (id.includes('@tanstack/')) return 'vendor-query'
          if (isPkg(id, 'react-router') || isPkg(id, 'react-router-dom') || isPkg(id, 'cookie')) {
            return 'vendor-router'
          }
          if (id.includes('lucide-react')) return 'vendor-icons'
          if (isPkg(id, 'react-dom') || isPkg(id, 'react') || isPkg(id, 'scheduler')) {
            return 'vendor-react'
          }
        },
      },
    },
  },
  server: {
    host: '0.0.0.0',
    port: 8888,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        ws: true,
      },
    },
  },
})
