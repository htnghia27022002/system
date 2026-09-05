'use client'

import { useQuery } from '@tanstack/react-query'

import { useAuthStore } from '@/store/auth-store'

import { profileApi } from '../services/profile-api'

export const profileQueryKey = ['auth', 'profile'] as const

export function useProfile() {
  const storeUser = useAuthStore((s) => s.user)

  return useQuery({
    queryKey: profileQueryKey,
    queryFn: () => profileApi.getProfile(),
    placeholderData: storeUser ?? undefined,
  })
}
