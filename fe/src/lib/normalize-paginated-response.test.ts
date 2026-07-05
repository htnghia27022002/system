import { describe, expect, it } from 'vitest'

import { normalizePaginatedResponse } from './normalize-paginated-response'

describe('normalizePaginatedResponse', () => {
  it('passes through paginated items shape', () => {
    const result = normalizePaginatedResponse<{ id: string }>({
      items: [{ id: '1' }],
      total: 1,
      page: 1,
      pageSize: 50,
    })

    expect(result.items).toHaveLength(1)
    expect(result.total).toBe(1)
  })

  it('supports legacy data array wrapper', () => {
    const result = normalizePaginatedResponse<{ id: string }>({
      data: [{ id: '1' }, { id: '2' }],
      total: 2,
      page: 1,
      pageSize: 50,
    })

    expect(result.items).toHaveLength(2)
  })

  it('supports raw array responses', () => {
    const result = normalizePaginatedResponse<{ id: string }>([
      { id: '1' },
      { id: '2' },
    ])

    expect(result.items).toHaveLength(2)
    expect(result.total).toBe(2)
  })
})
