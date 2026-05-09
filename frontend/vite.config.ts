import { defineConfig } from 'vitest/config'
import { loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

/**
 * Headwind MDM API proxy
 * ----------------------
 * - ROOT deployment → `http://localhost:8080/rest/...` — set `VITE_BACKEND_CONTEXT=` (empty) in `.env.development`.
 * - `launcher.war` only → `/launcher/rest/...` — set `VITE_BACKEND_CONTEXT=/launcher`.
 * - `mvn tomcat7:run` in `backend/server` → often port **9090**, ROOT — set `TOMCAT_PORT=9090` and `VITE_BACKEND_CONTEXT=`.
 *
 * Env: `TOMCAT_PORT`, `VITE_BACKEND_CONTEXT` (empty = no URL prefix before `/rest`).
 */
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const tomcatPort = env.TOMCAT_PORT || '8080'
  const rawCtx = env.VITE_BACKEND_CONTEXT
  const backendContext =
    rawCtx === undefined ? '' : rawCtx.replace(/\/$/, '')

  const restProxy = {
    '/rest': {
      target: `http://localhost:${tomcatPort}`,
      changeOrigin: true,
      secure: false,
      cookieDomainRewrite: '',
      cookiePathRewrite: '/',
      rewrite: (reqPath: string) => (backendContext ? `${backendContext}${reqPath}` : reqPath),
    },
  }

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      // Default Vite may bind only to IPv6 [::1]; Windows users opening http://127.0.0.1:5173 then see ERR_CONNECTION_REFUSED.
      host: true,
      proxy: restProxy,
    },
    // `vite preview` has no proxy unless configured — without this, `/rest/*` hits the static server → 404.
    preview: {
      host: true,
      proxy: restProxy,
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: ['./src/test/setup.ts'],
      coverage: {
        provider: 'v8',
        reporter: ['text', 'lcov'],
      },
    },
  }
})
