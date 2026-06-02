import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { mockAuthApi, MockAuthError } from '@/services/mock/auth.mock'

import type { AuthResponse, LoginRequest, RegisterRequest } from '../types'

export { MockAuthError }

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

    await apiClient.post<void>('/auth/logout')
  },
}
