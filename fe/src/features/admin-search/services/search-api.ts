import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'

import type { SearchQueryParams, SearchResponse } from '../types'

function toQueryParams(params: SearchQueryParams): Record<string, string | number> {
  const query: Record<string, string | number> = {
    q: params.q,
    page: params.page ?? 1,
    pageSize: params.pageSize ?? 20,
  }

  if (params.types?.length) {
    query.types = params.types.join(',')
  }

  return query
}

const mockHits = (q: string): SearchResponse => ({
  hits: [
    {
      entityType: 'user',
      entityId: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa',
      title: 'Admin User',
      snippet: `admin@example.com ${q}`,
      metadata: { email: 'admin@example.com', status: 'active' },
      updatedAt: new Date().toISOString(),
    },
  ],
  pagination: {
    page: 1,
    pageSize: 20,
    total: 1,
    totalPages: 1,
  },
  degraded: true,
})

export const searchApi = {
  search(params: SearchQueryParams): Promise<SearchResponse> {
    if (env.VITE_USE_MOCK_API) {
      const q = params.q.trim()
      if (!q) {
        return Promise.resolve({
          hits: [],
          pagination: { page: 1, pageSize: 20, total: 0, totalPages: 0 },
        })
      }
      return Promise.resolve(mockHits(q))
    }

    return apiClient
      .get<SearchResponse>('/admin/search', {
        params: toQueryParams(params),
        skipNavLoading: true,
      })
      .then((response) => response.data)
  },
}
