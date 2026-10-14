import type { NextConfig } from 'next'
import withSerwistInit from '@serwist/next'

/**
 * UUID-shaped public capture paths proxy to BE so local/split FE (no nginx)
 * still accepts ingest at `/tools/webhooks/{uuid}` without a Next owner page.
 */
const WEBHOOK_UUID_SOURCE =
  '/tools/webhooks/:uuid([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})'

function webhookCaptureDestination(): string {
  // Prefer server-only base for Docker FE→BE rewrites (container DNS).
  // Browser-facing NEXT_PUBLIC_* often points at localhost/system.local and
  // is unreachable from inside the FE container.
  const apiBase = (
    process.env.WEBHOOK_CAPTURE_API_BASE ??
    process.env.NEXT_PUBLIC_API_BASE_URL ??
    'http://localhost:8080/api'
  ).replace(/\/$/, '')
  return `${apiBase}/webhooks/capture/:uuid`
}

const withSerwist = withSerwistInit({
  swSrc: 'src/app/sw.ts',
  swDest: 'public/sw.js',
  disable: process.env.NODE_ENV === 'development',
})

/**
 * Baseline headers for every page when the FE runs without nginx (standalone
 * deploy). No CSP here: a useful one needs per-deploy API/asset origins — add it
 * at the edge or with a nonce middleware when the project knows them.
 */
const SECURITY_HEADERS = [
  { key: 'X-Content-Type-Options', value: 'nosniff' },
  { key: 'X-Frame-Options', value: 'DENY' },
  { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
  {
    key: 'Permissions-Policy',
    value: 'camera=(), microphone=(), geolocation=(), payment=()',
  },
]

const nextConfig: NextConfig = {
  poweredByHeader: false,
  async headers() {
    return [{ source: '/:path*', headers: SECURITY_HEADERS }]
  },
  async redirects() {
    return [
      {
        source: '/admin/users/roles',
        destination: '/admin/roles',
        permanent: true,
      },
      {
        source: '/admin/permissions',
        destination: '/admin/roles',
        permanent: true,
      },
    ]
  },
  async rewrites() {
    return [
      {
        source: WEBHOOK_UUID_SOURCE,
        destination: webhookCaptureDestination(),
      },
    ]
  },
}

export default withSerwist(nextConfig)
