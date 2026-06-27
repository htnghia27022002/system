'use client'

import { useQuery } from '@tanstack/react-query'

import { env } from '@/config/env'

import { authApi } from '../services/auth-api'

export function useOAuthProviders() {
  return useQuery({
    queryKey: ['auth', 'oauth-providers'],
    queryFn: () => authApi.getOAuthProviders(),
    enabled: !env.VITE_USE_MOCK_API,
    staleTime: 60_000,
  })
}

export function useIsGoogleOAuthEnabled(): boolean {
  const { data } = useOAuthProviders()
  return (data ?? []).includes('google')
}
