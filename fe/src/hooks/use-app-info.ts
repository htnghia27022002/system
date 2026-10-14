'use client'

import { useQuery } from '@tanstack/react-query'

import { fetchAppInfo } from '@/services/app-info-api'

/** GET /api/app/info, cached for the session (invalidated after settings change). */
export function useAppInfo() {
  return useQuery({
    queryKey: ['app', 'info'],
    queryFn: fetchAppInfo,
    staleTime: Infinity,
    retry: 1,
  })
}

/**
 * Whether self sign-up is open (`auth.registrationEnabled`). Optimistically true
 * while loading or when the BE does not report it.
 */
export function useRegistrationEnabled(): boolean {
  const { data } = useAppInfo()
  return data?.settings['auth.registrationEnabled'] !== false
}
