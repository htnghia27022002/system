import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { authApi } from '@/features/auth/services/auth-api'
import type { AuthUser } from '@/types/auth'

import type { ChangePasswordInput, UpdateProfileInput } from '../types'

async function mockUpdateProfile(
  input: UpdateProfileInput,
): Promise<AuthUser> {
  const user = await authApi.me()
  return {
    ...user,
    name: input.name,
    phone: input.phone ?? '',
    general: input.general ?? '',
    birthday: input.birthday || null,
    address: input.address ?? '',
    socialLinks: input.socialLinks ?? [],
  }
}

export const profileApi = {
  async getProfile(): Promise<AuthUser> {
    return authApi.me()
  },

  async updateProfile(input: UpdateProfileInput): Promise<AuthUser> {
    if (env.VITE_USE_MOCK_API) {
      return mockUpdateProfile(input)
    }
    const { data } = await apiClient.patch<AuthUser>('/auth/profile', input)
    return data
  },

  async uploadAvatar(file: File): Promise<AuthUser> {
    if (env.VITE_USE_MOCK_API) {
      const user = await authApi.me()
      const objectUrl = URL.createObjectURL(file)
      return { ...user, avatarUrl: objectUrl }
    }
    const formData = new FormData()
    formData.append('file', file)
    const { data } = await apiClient.post<AuthUser>(
      '/auth/profile/avatar',
      formData,
      {
        headers: { 'Content-Type': 'multipart/form-data' },
      },
    )
    return data
  },

  async changePassword(input: ChangePasswordInput): Promise<void> {
    if (env.VITE_USE_MOCK_API) {
      return
    }
    await apiClient.post('/auth/change-password', input)
  },
}
