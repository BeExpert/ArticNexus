import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      port: 5173,
      proxy: mode === 'development'
        ? {
            '/api': {
              target: 'http://localhost:8080',
              changeOrigin: true
            }
          }
        : undefined
    },
    define: {
      __APP_TITLE__: JSON.stringify(env.VITE_APP_TITLE || 'ArticNexus')
    }
  }
})