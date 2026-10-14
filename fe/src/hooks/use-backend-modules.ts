'use client'

import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'

import { backendModuleNames } from '@/config/modules'
import { fetchAppInfo } from '@/services/app-info-api'

/**
 * Enabled BE modules as a Set. Undefined while loading; if the endpoint fails
 * (e.g. an older BE without /app/info) every known module is assumed enabled.
 */
export function useBackendModules(): ReadonlySet<string> | undefined {
  const query = useQuery({
    queryKey: ['app', 'info'],
    queryFn: fetchAppInfo,
    staleTime: Infinity,
    retry: 1,
  })

  return useMemo(() => {
    if (query.data) return new Set(query.data.modules)
    if (query.isError) return new Set(backendModuleNames)
    return undefined
  }, [query.data, query.isError])
}
