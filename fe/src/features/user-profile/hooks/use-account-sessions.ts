'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { authTokenService } from '@/services/auth-token-service'

import { sessionsApi } from '../services/sessions-api'
import type { AccountSession } from '../types'

export const accountSessionsQueryKey = ['auth', 'sessions'] as const

function apiErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as
      | { message?: string; error?: string }
      | undefined
    return data?.message || data?.error || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function markCurrentSessions(
  items: AccountSession[],
  storedSessionId: string | null,
): AccountSession[] {
  return items.map((item) => ({
    ...item,
    current: item.current || Boolean(storedSessionId && item.id === storedSessionId),
  }))
}

export function sortAccountSessions(items: AccountSession[]): AccountSession[] {
  return [...items].sort((a, b) => {
    if (a.current !== b.current) return a.current ? -1 : 1
    return new Date(b.lastUsedAt).getTime() - new Date(a.lastUsedAt).getTime()
  })
}

export function useAccountSessions() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  const query = useQuery({
    queryKey: accountSessionsQueryKey,
    queryFn: () => sessionsApi.listSessions(),
  })

  const storedSessionId = authTokenService.getSessionId()
  const items = sortAccountSessions(
    markCurrentSessions(query.data?.items ?? [], storedSessionId),
  )
  const hasOtherSessions = items.some((item) => !item.current)

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: accountSessionsQueryKey })
  }

  const revokeSession = useMutation({
    mutationFn: (id: string) => sessionsApi.revokeSession(id),
    onSuccess: () => {
      toast.success(t('profile.sessions.toasts.revoked'))
      invalidate()
    },
    onError: (error) => {
      const notFound = axios.isAxiosError(error) && error.response?.status === 404
      toast.error(
        notFound
          ? t('profile.sessions.errors.revokeNotFound')
          : apiErrorMessage(error, t('profile.sessions.errors.revokeFailed')),
      )
      invalidate()
    },
  })

  const revokeOthers = useMutation({
    mutationFn: () => {
      const sessionId = authTokenService.getSessionId()
      if (!sessionId) {
        throw new Error(t('profile.sessions.errors.revokeOthersFailed'))
      }
      return sessionsApi.revokeOthers(sessionId)
    },
    onSuccess: (data) => {
      toast.success(t('profile.sessions.toasts.revokedOthers'))
      queryClient.setQueryData(accountSessionsQueryKey, data)
      invalidate()
    },
    onError: (error) => {
      toast.error(
        apiErrorMessage(error, t('profile.sessions.errors.revokeOthersFailed')),
      )
      invalidate()
    },
  })

  return { query, items, hasOtherSessions, revokeSession, revokeOthers }
}
