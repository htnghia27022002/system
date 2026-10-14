import { beforeEach, describe, expect, it } from 'vitest'

import { authTokenService } from './auth-token-service'

describe('authTokenService', () => {
  beforeEach(() => window.localStorage.clear())

  it('does not persist a refresh token the BE kept in an httpOnly cookie', () => {
    window.localStorage.setItem('refresh_token', 'stale')
    authTokenService.setTokens({ accessToken: 'a', sessionId: 's' })
    expect(authTokenService.getAccessToken()).toBe('a')
    expect(authTokenService.getRefreshToken()).toBeNull()
  })

  it('stores a refresh token returned in JSON (cookie mode off)', () => {
    authTokenService.setTokens({ accessToken: 'a', refreshToken: 'r' })
    expect(authTokenService.getRefreshToken()).toBe('r')
    authTokenService.clearTokens()
    expect(authTokenService.getRefreshToken()).toBeNull()
  })
})
