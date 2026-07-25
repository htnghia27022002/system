import { env } from '@/config/env'

const OAUTH_CALLBACK_PATH = '/auth/callback'

/**
 * Google OAuth redirect_uri must match the URL the browser is on.
 * Prefer window.location.origin so mis-set NEXT_PUBLIC_SITE_URL (or its
 * localhost:3000 fallback) cannot send users to the wrong host after login.
 */
export function getOAuthCallbackUri(): string {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin.replace(/\/$/, '')}${OAUTH_CALLBACK_PATH}`
  }
  const base = env.VITE_SITE_URL.replace(/\/$/, '')
  return `${base}${OAUTH_CALLBACK_PATH}`
}

export function getGoogleOAuthStartUrl(): string {
  const apiBase = env.VITE_API_BASE_URL.replace(/\/$/, '')
  const redirectUri = encodeURIComponent(getOAuthCallbackUri())
  return `${apiBase}/auth/oauth/google/start?redirect_uri=${redirectUri}`
}

export { OAUTH_CALLBACK_PATH }
