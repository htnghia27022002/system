import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { authTokenService } from '@/services/auth-token-service'
import { mockAuthApi, MockAuthError } from '@/services/mock/auth.mock'

import type { AuthResponse, LoginRequest, RegisterRequest } from '../types'

export { MockAuthError }

type OAuthProvidersResponse = {
  providers: string[]
}

export const authApi = {
  async login(payload: LoginRequest): Promise<AuthResponse> {
    if (env.VITE_USE_MOCK_API) {
      return mockAuthApi.login(payload)
    }

    const { data } = await apiClient.post<AuthResponse>('/auth/login', payload)
    return data
  },

  async register(payload: RegisterRequest): Promise<AuthResponse> {
    if (env.VITE_USE_MOCK_API) {
      return mockAuthApi.register(payload)
    }

    const { data } = await apiClient.post<AuthResponse>(
      '/auth/register',
      payload,
    )
    return data
  },

  async logout(): Promise<void> {
    if (env.VITE_USE_MOCK_API) {
      return
    }

    const refreshToken = authTokenService.getRefreshToken()
    if (!refreshToken) {
      return
    }

    await apiClient.post<void>('/auth/logout', { refreshToken })
  },

  async getOAuthProviders(): Promise<string[]> {
    if (env.VITE_USE_MOCK_API) {
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
    if (env.VITE_USE_MOCK_API) {
      throw new Error('OAuth is not available in mock mode')
    }

    const { data } = await apiClient.post<AuthResponse>(
      `/auth/oauth/${provider}/callback`,
      { code, redirectUri },
    )
    return data
  },
}
