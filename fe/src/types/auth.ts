export type UserRole = 'admin' | 'user'

export type AuthTokens = {
  accessToken: string
  refreshToken: string
  /** Refresh-token row id (UUID). Never put this in URLs or query strings. */
  sessionId?: string | null
}

export type SocialLink = {
  label?: string
  url: string
}

export type AuthUser = {
  id: string
  email: string
  name: string
  role: UserRole
  roleId: string
  permissions: string[]
  superAdmin?: boolean
  phone?: string
  avatarUrl?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
  hasPassword?: boolean
}

export type JwtPayload = {
  sub: string
  exp: number
  iat?: number
  email?: string
  name?: string
  role?: UserRole
  roleId?: string
  permissions?: string[]
  superAdmin?: boolean
}
