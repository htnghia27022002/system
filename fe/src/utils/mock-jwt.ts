import type { UserRole } from '@/types/auth'

type MockJwtPayload = {
  sub: string
  email: string
  name: string
  role: UserRole
  exp: number
  iat: number
}

function base64UrlEncode(value: string): string {
  return btoa(value).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

export function createMockJwt(payload: MockJwtPayload): string {
  const header = base64UrlEncode(JSON.stringify({ alg: 'none', typ: 'JWT' }))
  const body = base64UrlEncode(JSON.stringify(payload))
  return `${header}.${body}.mock-signature`
}

export function createMockAuthTokens(
  userId: string,
  email: string,
  name: string,
  role: UserRole,
) {
  const now = Math.floor(Date.now() / 1000)
  const base = { sub: userId, email, name, role }
  const accessPayload: MockJwtPayload = {
    ...base,
    iat: now,
    exp: now + 60 * 60,
  }
  const refreshPayload: MockJwtPayload = {
    ...base,
    iat: now,
    exp: now + 60 * 60 * 24 * 7,
  }

  return {
    accessToken: createMockJwt(accessPayload),
    refreshToken: createMockJwt(refreshPayload),
  }
}
