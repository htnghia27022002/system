import type { PaginatedResponse } from '@/features/access-control/types'

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export function normalizePaginatedResponse<T>(
  payload: unknown,
  fallbackPageSize = 10,
): PaginatedResponse<T> {
  if (Array.isArray(payload)) {
    return {
      items: payload as T[],
      total: payload.length,
      page: 1,
      pageSize: payload.length || fallbackPageSize,
    }
  }

  if (!isRecord(payload)) {
    return {
      items: [],
      total: 0,
      page: 1,
      pageSize: fallbackPageSize,
    }
  }

  const items = payload.items ?? payload.data
  if (!Array.isArray(items)) {
    return {
      items: [],
      total: typeof payload.total === 'number' ? payload.total : 0,
      page: typeof payload.page === 'number' ? payload.page : 1,
      pageSize:
        typeof payload.pageSize === 'number'
          ? payload.pageSize
          : fallbackPageSize,
    }
  }

  return {
    items: items as T[],
    total: typeof payload.total === 'number' ? payload.total : items.length,
    page: typeof payload.page === 'number' ? payload.page : 1,
    pageSize:
      typeof payload.pageSize === 'number'
        ? payload.pageSize
        : fallbackPageSize,
  }
}
