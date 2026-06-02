export type UserRole = 'admin' | 'user'

export type AuthTokens = {
  accessToken: string
  refreshToken: string
}

export type AuthUser = {
  id: string
  email: string
  name: string
  role: UserRole
}

export type JwtPayload = {
  sub: string
  exp: number
  iat?: number
  email?: string
  name?: string
  role?: UserRole
}
