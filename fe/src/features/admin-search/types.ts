export type SearchEntityType = 'user' | 'role' | 'permission'

export type SearchQueryParams = {
  q: string
  types?: SearchEntityType[]
  page?: number
  pageSize?: number
}

export type SearchHit = {
  entityType: SearchEntityType
  entityId: string
  title: string
  snippet?: string
  metadata?: Record<string, string>
  updatedAt: string
}

export type SearchPagination = {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

export type SearchResponse = {
  hits: SearchHit[]
  pagination: SearchPagination
  degraded?: boolean
}
