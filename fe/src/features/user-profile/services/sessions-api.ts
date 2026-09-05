import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { authTokenService } from '@/services/auth-token-service'

import type { AccountSession, AccountSessionList } from '../types'

export const SESSION_ID_HEADER = 'X-Session-Id'

type RawSessionListItem = {
  id: string
  createdAt: string
  expiresAt: string
  lastUsedAt: string
  ipAddress?: string | null
  userAgent?: string | null
  current?: boolean
}

export function mapSessionListItem(raw: RawSessionListItem): AccountSession {
  return {
    id: raw.id,
    createdAt: raw.createdAt,
    expiresAt: raw.expiresAt,
    lastUsedAt: raw.lastUsedAt,
    ipAddress: raw.ipAddress ?? null,
    userAgent: raw.userAgent ?? null,
    current: Boolean(raw.current),
  }
}

export function mapSessionList(raw: { items?: RawSessionListItem[] | null }): AccountSessionList {
  return {
    items: (raw.items ?? []).map(mapSessionListItem),
  }
}

function sessionIdHeaders(): Record<string, string> | undefined {
  const sessionId = authTokenService.getSessionId()
  if (!sessionId) return undefined
  return { [SESSION_ID_HEADER]: sessionId }
}

function mockListSessions(): AccountSessionList {
  const now = new Date()
  const sessionId = authTokenService.getSessionId() ?? 'mock-session-id'
  return mapSessionList({
    items: [
      {
        id: sessionId,
        createdAt: now.toISOString(),
        expiresAt: new Date(now.getTime() + 14 * 24 * 60 * 60 * 1000).toISOString(),
        lastUsedAt: now.toISOString(),
        ipAddress: '127.0.0.1',
        userAgent:
          typeof navigator !== 'undefined' ? navigator.userAgent : null,
        current: true,
      },
    ],
  })
}

export const sessionsApi = {
  async listSessions(): Promise<AccountSessionList> {
    if (env.VITE_USE_MOCK_API) {
      return mockListSessions()
    }

    const { data } = await apiClient.get<{ items?: RawSessionListItem[] }>(
      '/auth/sessions',
      { headers: sessionIdHeaders() },
    )
    return mapSessionList(data)
  },

  async revokeSession(id: string): Promise<void> {
    if (env.VITE_USE_MOCK_API) {
      return
    }

    await apiClient.delete(`/auth/sessions/${id}`, {
      headers: sessionIdHeaders(),
    })
  },

  async revokeOthers(sessionId: string): Promise<AccountSessionList> {
    if (env.VITE_USE_MOCK_API) {
      return mockListSessions()
    }

    const { data } = await apiClient.post<{ items?: RawSessionListItem[] }>(
      '/auth/sessions/revoke-others',
      { sessionId },
      { headers: sessionIdHeaders() },
    )
    return mapSessionList(data)
  },
}
