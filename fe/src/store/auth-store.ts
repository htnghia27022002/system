'use client'

import { create } from 'zustand'

import { authTokenService } from '@/services/auth-token-service'
import { env } from '@/config/env'
import type { AuthResponse } from '@/features/auth'
import { authApi } from '@/features/auth/services/auth-api'
import type { AuthUser } from '@/types/auth'

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
  hasHydrated: boolean
  sessionSynced: boolean
  setAccessToken: (token: string | null) => void
  signIn: (response: AuthResponse) => void
  hydrateFromStorage: () => void
  syncSession: () => Promise<void>
  signOut: () => void
}

function buildState(
  accessToken: string | null,
  user: AuthUser | null = null,
  sessionSynced = false,
) {
  const resolvedUser = user ?? userFromAccessToken(accessToken)
  return {
    accessToken,
    user: resolvedUser,
    isAdmin: resolvedUser?.role === 'admin',
    sessionSynced,
  }
}

export const useAuthStore = create<AuthState>((set, get) => ({
  ...buildState(null),
  hasHydrated: false,
  sessionSynced: false,
  setAccessToken: (token) =>
    set({ ...buildState(token, null, false), hasHydrated: true, sessionSynced: false }),
  signIn: (response) => {
    authTokenService.setTokens({
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
    })
    set({
      ...buildState(response.accessToken, response.user, true),
      hasHydrated: true,
      sessionSynced: true,
    })
  },
  hydrateFromStorage: () => {
    const accessToken = authTokenService.getAccessToken()
    set({
      ...buildState(accessToken, null, false),
      hasHydrated: true,
      sessionSynced: !accessToken,
    })
  },
  syncSession: async () => {
    const accessToken = get().accessToken ?? authTokenService.getAccessToken()
    if (!accessToken) {
      set({ sessionSynced: true })
      return
    }

    if (env.VITE_USE_MOCK_API) {
      set({
        ...buildState(accessToken, get().user, true),
        hasHydrated: true,
        sessionSynced: true,
      })
      return
    }

    try {
      const user = await authApi.me()
      set({
        ...buildState(accessToken, user, true),
        hasHydrated: true,
        sessionSynced: true,
      })
    } catch {
      authTokenService.clearTokens()
      set({ ...buildState(null, null, true), hasHydrated: true, sessionSynced: true })
    }
  },
  signOut: () => {
    authTokenService.clearTokens()
    set({ ...buildState(null, null, true), hasHydrated: true, sessionSynced: true })
  },
}))

export function getPostLoginPath(): string {
  return '/admin'
}

export function selectPermissions(state: AuthState): string[] {
  return state.user?.permissions ?? []
}
