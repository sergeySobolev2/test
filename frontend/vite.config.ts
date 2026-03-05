import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'
import basicSsl from '@vitejs/plugin-basic-ssl'
import fs from 'node:fs'
import path from 'node:path'

const useHttps = process.env.VITE_HTTPS === 'true'

function resolveHttpsConfig() {
  if (!useHttps) return undefined
  const certDir = path.resolve(__dirname, 'certs')
  const certFile = path.join(certDir, 'localhost.pem')
  const keyFile = path.join(certDir, 'localhost-key.pem')
  try {
    if (fs.existsSync(certFile) && fs.existsSync(keyFile)) {
      return {
        cert: fs.readFileSync(certFile),
        key: fs.readFileSync(keyFile),
      }
    }
  } catch {}
  return undefined
}

const backendProxy = {
  target: 'http://localhost:8082',
  changeOrigin: true,
}


const httpsConfig = resolveHttpsConfig()
const useBasicSsl = useHttps && !httpsConfig

const basePath = process.env.VITE_BASE || '/'

// Единый порт 5196 для dev + preview, и прокси в обоих режимах
export default defineConfig({
  base: basePath,
  plugins: [
    react(),
    // Если нет пользовательских сертификатов — используем basicSsl для self-signed
    useBasicSsl ? basicSsl() : undefined,
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: [
        'placeholder.svg',
        'placeholder-partition.svg',
        'placeholder-symptom.svg',
        'partition-cart-icon.svg',
        'draft-icon.svg',
      ],
      manifest: {
        name: 'Partition Soundproofing',
        short_name: 'Partition',
        start_url: '.',
        display: 'standalone',
        background_color: '#ffffff',
        theme_color: '#0d6efd',
        icons: [
          { src: 'icons/icon-192.png', sizes: '192x192', type: 'image/png' },
          { src: 'icons/icon-512.png', sizes: '512x512', type: 'image/png' }
        ]
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,svg,png,ico}'],
        runtimeCaching: [
          {
            urlPattern: /\/api\/partitions.*/,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-partitions',
              networkTimeoutSeconds: 5,
              cacheableResponse: { statuses: [0, 200] }
            }
          }
        ]
      }
    })
  ],
  server: {
    port: 5196,
    strictPort: true,
    https: httpsConfig,
    proxy: {
      '/api': backendProxy,
      '/swagger': backendProxy
    }
  },
  preview: {
    port: 5196,
    https: httpsConfig,
    proxy: {
      '/api': backendProxy,
      '/swagger': backendProxy
    }
  }
})
