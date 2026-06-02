import { jwtDecode } from 'jwt-decode'

import type { AuthTokens, JwtPayload } from '@/types/auth'

const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

export const authTokenService = {
  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY)
  },

  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_TOKEN_KEY)
  },

  setTokens(tokens: AuthTokens): void {
    localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken)
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken)
  },

  clearTokens(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
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
