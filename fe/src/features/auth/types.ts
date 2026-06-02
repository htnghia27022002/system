import type { AuthTokens, AuthUser, UserRole } from '@/types/auth'

export type { UserRole }

export type LoginRequest = {
  email: string
  password: string
}

export type RegisterRequest = {
  name: string
  email: string
  password: string
}

export type AuthResponse = AuthTokens & {
  user: AuthUser
}
