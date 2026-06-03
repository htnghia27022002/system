import { create } from 'zustand'

import { authTokenService } from '@/services/auth-token-service'
import type { AuthResponse } from '@/features/auth'
import type { AuthUser, UserRole } from '@/types/auth'

function userFromAccessToken(token: string | null): AuthUser | null {
  if (!token) {
    return null
  }

  try {
    const payload = authTokenService.decodeAccessToken(token)
    if (!payload.email || !payload.role) {
      return null
    }

    return {
      id: payload.sub,
      email: payload.email,
      name: payload.name ?? payload.email,
      role: payload.role,
      roleId: payload.roleId ?? '',
      permissions: payload.permissions ?? [],
    }
  } catch {
    return null
  }
}

type AuthState = {
  accessToken: string | null
  user: AuthUser | null
  isAdmin: boolean
  setAccessToken: (token: string | null) => void
  signIn: (response: AuthResponse) => void
  hydrateFromStorage: () => void
  signOut: () => void
}

function buildState(accessToken: string | null, user: AuthUser | null = null) {
  const resolvedUser = user ?? userFromAccessToken(accessToken)
  return {
    accessToken,
    user: resolvedUser,
    isAdmin: resolvedUser?.role === 'admin',
  }
}

export const useAuthStore = create<AuthState>((set) => ({
  ...buildState(authTokenService.getAccessToken()),
  setAccessToken: (token) => set(buildState(token)),
  signIn: (response) => {
    authTokenService.setTokens({
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
    })
    set(buildState(response.accessToken, response.user))
  },
  hydrateFromStorage: () =>
    set(buildState(authTokenService.getAccessToken())),
  signOut: () => {
    authTokenService.clearTokens()
    set(buildState(null))
  },
}))

export function getPostLoginPath(role: UserRole): string {
  return role === 'admin' ? '/admin' : '/'
}

export function selectPermissions(state: AuthState): string[] {
  return state.user?.permissions ?? []
}
