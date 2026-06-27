import { env } from '@/config/env'

const OAUTH_CALLBACK_PATH = '/auth/callback'

export function getOAuthCallbackUri(): string {
  const base = env.VITE_SITE_URL.replace(/\/$/, '')
  return `${base}${OAUTH_CALLBACK_PATH}`
}

export function getGoogleOAuthStartUrl(): string {
  const apiBase = env.VITE_API_BASE_URL.replace(/\/$/, '')
  const redirectUri = encodeURIComponent(getOAuthCallbackUri())
  return `${apiBase}/auth/oauth/google/start?redirect_uri=${redirectUri}`
}

export { OAUTH_CALLBACK_PATH }
