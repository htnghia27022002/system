'use client'

import { useQuery } from '@tanstack/react-query'

import { searchApi } from '../services/search-api'
import type { SearchQueryParams } from '../types'

const searchKey = ['admin', 'search'] as const

export function useAdminSearch(params: SearchQueryParams) {
  const trimmed = params.q.trim()

  return useQuery({
    queryKey: [...searchKey, params],
    queryFn: () => searchApi.search({ ...params, q: trimmed }),
    enabled: trimmed.length > 0,
    meta: { skipNavLoading: true },
  })
}
