import { beforeEach, describe, expect, it, vi } from 'vitest'

const { getSessionId, apiGet, apiDelete, apiPost } = vi.hoisted(() => ({
  getSessionId: vi.fn(() => '8c2a0c3e-1b2a-4d5e-9f10-111213141516'),
  apiGet: vi.fn(),
  apiDelete: vi.fn(),
  apiPost: vi.fn(),
}))

vi.mock('@/config/env', () => ({
  env: { VITE_USE_MOCK_API: false },
}))

vi.mock('@/services/auth-token-service', () => ({
  authTokenService: {
    getSessionId,
    getRefreshToken: vi.fn(() => 'secret-refresh-token-must-not-leak'),
  },
}))

vi.mock('@/services/api-client', () => ({
  apiClient: {
    get: (...args: unknown[]) => apiGet(...args),
    delete: (...args: unknown[]) => apiDelete(...args),
    post: (...args: unknown[]) => apiPost(...args),
  },
}))

import { mapSessionListItem, sessionsApi, SESSION_ID_HEADER } from './sessions-api'

const currentId = '8c2a0c3e-1b2a-4d5e-9f10-111213141516'

describe('mapSessionListItem', () => {
  it('maps camelCase list fields including nullable ip and user-agent', () => {
    expect(
      mapSessionListItem({
        id: currentId,
        createdAt: '2026-09-04T10:00:00Z',
        expiresAt: '2026-09-18T10:00:00Z',
        lastUsedAt: '2026-09-04T12:30:00Z',
        ipAddress: '203.0.113.10',
        userAgent: 'Mozilla/5.0',
        current: true,
      }),
    ).toEqual({
      id: currentId,
      createdAt: '2026-09-04T10:00:00Z',
      expiresAt: '2026-09-18T10:00:00Z',
      lastUsedAt: '2026-09-04T12:30:00Z',
      ipAddress: '203.0.113.10',
      userAgent: 'Mozilla/5.0',
      current: true,
    })

    expect(
      mapSessionListItem({
        id: '9d3b1d4f-2c3b-5e6f-0a21-222324252627',
        createdAt: '2026-09-03T08:00:00Z',
        expiresAt: '2026-09-17T08:00:00Z',
        lastUsedAt: '2026-09-03T09:00:00Z',
        ipAddress: null,
        userAgent: null,
        current: false,
      }),
    ).toMatchObject({
      id: '9d3b1d4f-2c3b-5e6f-0a21-222324252627',
      ipAddress: null,
      userAgent: null,
      current: false,
    })
  })
})

describe('sessionsApi', () => {
  beforeEach(() => {
    apiGet.mockReset()
    apiDelete.mockReset()
    apiPost.mockReset()
    getSessionId.mockReturnValue(currentId)
  })

  it('lists sessions with X-Session-Id and never puts the refresh token in the request', async () => {
    apiGet.mockResolvedValue({
      data: {
        items: [
          {
            id: currentId,
            createdAt: '2026-09-04T10:00:00Z',
            expiresAt: '2026-09-18T10:00:00Z',
            lastUsedAt: '2026-09-04T12:30:00Z',
            ipAddress: '203.0.113.10',
            userAgent: 'Mozilla/5.0',
            current: true,
          },
        ],
      },
    })

    const result = await sessionsApi.listSessions()

    expect(result.items[0]).toMatchObject({
      id: currentId,
      createdAt: '2026-09-04T10:00:00Z',
      expiresAt: '2026-09-18T10:00:00Z',
      lastUsedAt: '2026-09-04T12:30:00Z',
      ipAddress: '203.0.113.10',
      userAgent: 'Mozilla/5.0',
      current: true,
    })
    expect(apiGet).toHaveBeenCalledWith('/auth/sessions', {
      headers: { [SESSION_ID_HEADER]: currentId },
    })
    const [url, config] = apiGet.mock.calls[0] as [
      string,
      { headers?: Record<string, string>; params?: Record<string, unknown> },
    ]
    expect(url).not.toMatch(/refresh/i)
    expect(JSON.stringify(config)).not.toContain('secret-refresh-token-must-not-leak')
    expect(config.params).toBeUndefined()
  })

  it('revokes a session by id without a query string', async () => {
    apiDelete.mockResolvedValue({ data: undefined })
    await sessionsApi.revokeSession('9d3b1d4f-2c3b-5e6f-0a21-222324252627')
    expect(apiDelete).toHaveBeenCalledWith(
      '/auth/sessions/9d3b1d4f-2c3b-5e6f-0a21-222324252627',
      { headers: { [SESSION_ID_HEADER]: currentId } },
    )
  })

  it('revokes others with a JSON sessionId body', async () => {
    apiPost.mockResolvedValue({ data: { items: [] } })
    await sessionsApi.revokeOthers(currentId)
    expect(apiPost).toHaveBeenCalledWith(
      '/auth/sessions/revoke-others',
      { sessionId: currentId },
      { headers: { [SESSION_ID_HEADER]: currentId } },
    )
  })
})
