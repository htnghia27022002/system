import { jwtDecode } from 'jwt-decode'

import type { AuthTokens, JwtPayload } from '@/types/auth'

const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

function safeLocalStorage(): Storage | null {
  if (typeof window === 'undefined') return null
  return window.localStorage
}

export const authTokenService = {
  getAccessToken(): string | null {
    return safeLocalStorage()?.getItem(ACCESS_TOKEN_KEY) ?? null
  },

  getRefreshToken(): string | null {
    return safeLocalStorage()?.getItem(REFRESH_TOKEN_KEY) ?? null
  },

  setTokens(tokens: AuthTokens): void {
    const storage = safeLocalStorage()
    if (!storage) return
    storage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
    storage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
  },

  clearTokens(): void {
    const storage = safeLocalStorage()
    if (!storage) return
    storage.removeItem(ACCESS_TOKEN_KEY)
    storage.removeItem(REFRESH_TOKEN_KEY)
  },

  decodeAccessToken(token: string): JwtPayload {
    return jwtDecode<JwtPayload>(token)
  },

  isAccessTokenExpired(token: string, skewSeconds = 30): boolean {
    const payload = this.decodeAccessToken(token)
    const now = Math.floor(Date.now() / 1000)
    return payload.exp <= now + skewSeconds
  },
}
