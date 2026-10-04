import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { authTokenService } from '@/services/auth-token-service'
import { loadAuthMock, MockAuthError } from '@/services/mock'

import type { AuthResponse, LoginRequest, RegisterRequest } from '../types'
import type { AuthUser } from '@/types/auth'

export { MockAuthError }

type OAuthProvidersResponse = {
  providers: string[]
}

export const authApi = {
  async login(payload: LoginRequest): Promise<AuthResponse> {
    if (env.USE_MOCK_API) {
      return (await loadAuthMock()).login(payload)
    }

    const { data } = await apiClient.post<AuthResponse>('/auth/login', payload)
    return data
  },

  async register(payload: RegisterRequest): Promise<AuthResponse> {
    if (env.USE_MOCK_API) {
      return (await loadAuthMock()).register(payload)
    }

    const { data } = await apiClient.post<AuthResponse>(
      '/auth/register',
      payload,
    )
    return data
  },

  async me(): Promise<AuthUser> {
    if (env.USE_MOCK_API) {
      return (await loadAuthMock()).me()
    }

    const { data } = await apiClient.get<AuthUser>('/auth/me')
    return data
  },

  async logout(): Promise<void> {
    if (env.USE_MOCK_API) {
      return
    }

    const refreshToken = authTokenService.getRefreshToken()
    if (!refreshToken) {
      return
    }

    await apiClient.post<void>('/auth/logout', { refreshToken })
  },

  async getOAuthProviders(): Promise<string[]> {
    if (env.USE_MOCK_API) {
      return []
    }

    const { data } = await apiClient.get<OAuthProvidersResponse>(
      '/auth/oauth/providers',
    )
    return data.providers ?? []
  },

  async oauthCallback(
    provider: string,
    code: string,
    redirectUri: string,
  ): Promise<AuthResponse> {
    if (env.USE_MOCK_API) {
      throw new Error('OAuth is not available in mock mode')
    }

    const { data } = await apiClient.post<AuthResponse>(
      `/auth/oauth/${provider}/callback`,
      { code, redirectUri },
    )
    return data
  },
}
