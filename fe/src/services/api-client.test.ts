import { describe, expect, it } from 'vitest'

import { shouldRefreshOn401 } from './api-client'

describe('shouldRefreshOn401', () => {
  it('never refreshes for public auth endpoints', () => {
    for (const url of [
      '/auth/login',
      '/auth/register',
      '/auth/refresh',
      '/auth/logout',
      '/auth/oauth/google/callback',
      '/auth/password/reset',
      '/auth/email/verify',
      '/auth/invite/accept',
      'http://localhost:8080/api/auth/login',
    ]) {
      expect(shouldRefreshOn401(url)).toBe(false)
    }
  })

  it('refreshes for protected endpoints', () => {
    for (const url of ['/auth/me', '/auth/sessions', '/auth/email/verification', '/admin/users', '/webhooks/inbox']) {
      expect(shouldRefreshOn401(url)).toBe(true)
    }
  })
})
