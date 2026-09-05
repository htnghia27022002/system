import type { SocialLink } from '@/types/auth'

export type { SocialLink }

export type UpdateProfileInput = {
  name: string
  phone?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
}

export type ChangePasswordInput = {
  currentPassword: string
  newPassword: string
}

export type AccountSession = {
  id: string
  createdAt: string
  expiresAt: string
  lastUsedAt: string
  ipAddress: string | null
  userAgent: string | null
  current: boolean
}

export type AccountSessionList = {
  items: AccountSession[]
}
