import { jwtDecode } from 'jwt-decode'

import type { AuthTokens, JwtPayload } from '@/types/auth'

const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const SESSION_ID_KEY = 'session_id'

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

  getSessionId(): string | null {
    return safeLocalStorage()?.getItem(SESSION_ID_KEY) ?? null
  },

  setTokens(tokens: AuthTokens): void {
    const storage = safeLocalStorage()
    if (!storage) return
    storage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
    storage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
    if (tokens.sessionId) {
      storage.setItem(SESSION_ID_KEY, tokens.sessionId)
    } else {
      storage.removeItem(SESSION_ID_KEY)
    }
  },

  clearTokens(): void {
    const storage = safeLocalStorage()
    if (!storage) return
    storage.removeItem(ACCESS_TOKEN_KEY)
    storage.removeItem(REFRESH_TOKEN_KEY)
    storage.removeItem(SESSION_ID_KEY)
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
