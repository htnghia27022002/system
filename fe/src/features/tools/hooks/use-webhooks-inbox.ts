'use client'

import {
  useMutation,
  useQuery,
  useQueryClient,
  type UseQueryOptions,
} from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { webhooksApi } from '../services/webhooks-api'
import type {
  ClearAllWebhookRequestsResponse,
  ListWebhookRequestsParams,
  SoftDeleteWebhookRequestResponse,
  WebhookInbox,
  WebhookRequestDetail,
  WebhookRequestListResponse,
} from '../types/webhooks'

export const webhooksInboxKey = ['tools', 'webhooks', 'inbox'] as const
export const webhooksRequestsKey = ['tools', 'webhooks', 'requests'] as const

const POLL_MS = 5_000

function isListResponse(
  value: unknown,
): value is WebhookRequestListResponse {
  return (
    typeof value === 'object' &&
    value !== null &&
    Array.isArray((value as WebhookRequestListResponse).items) &&
    typeof (value as WebhookRequestListResponse).hasMore === 'boolean'
  )
}

function applyInboxCounters(
  queryClient: ReturnType<typeof useQueryClient>,
  result: Pick<
    SoftDeleteWebhookRequestResponse,
    'activeCount' | 'lifetimeReceived'
  >,
) {
  queryClient.setQueryData(webhooksInboxKey, (prev: WebhookInbox | undefined) =>
    prev
      ? {
          ...prev,
          activeCount: result.activeCount,
          lifetimeReceived: result.lifetimeReceived,
        }
      : prev,
  )
}

function patchListItemRead(
  queryClient: ReturnType<typeof useQueryClient>,
  id: string,
  isRead: boolean,
) {
  queryClient.setQueriesData(
    {
      predicate: (query) =>
        query.queryKey[0] === webhooksRequestsKey[0] &&
        query.queryKey[1] === webhooksRequestsKey[1] &&
        query.queryKey[2] === webhooksRequestsKey[2] &&
        query.queryKey[3] !== 'detail',
    },
    (prev) => {
      if (!isListResponse(prev)) return prev
      return {
        ...prev,
        items: prev.items.map((item) =>
          item.id === id ? { ...item, isRead } : item,
        ),
      }
    },
  )
}

function applySoftDeleteCounters(
  queryClient: ReturnType<typeof useQueryClient>,
  deletedId: string,
  result: SoftDeleteWebhookRequestResponse,
) {
  applyInboxCounters(queryClient, result)

  queryClient.setQueriesData(
    {
      predicate: (query) =>
        query.queryKey[0] === webhooksRequestsKey[0] &&
        query.queryKey[1] === webhooksRequestsKey[1] &&
        query.queryKey[2] === webhooksRequestsKey[2] &&
        query.queryKey[3] !== 'detail',
    },
    (prev) => {
      if (!isListResponse(prev)) return prev
      const items = prev.items.filter((item) => item.id !== deletedId)
      const removed = items.length < prev.items.length ? 1 : 0
      return {
        ...prev,
        items,
        activeCount: result.activeCount,
        lifetimeReceived: result.lifetimeReceived,
        total: Math.max(0, prev.total - removed),
      }
    },
  )
  queryClient.removeQueries({
    queryKey: [...webhooksRequestsKey, 'detail', deletedId],
  })
}

function applyClearAllCounters(
  queryClient: ReturnType<typeof useQueryClient>,
  result: ClearAllWebhookRequestsResponse,
) {
  applyInboxCounters(queryClient, result)

  queryClient.setQueriesData(
    {
      predicate: (query) =>
        query.queryKey[0] === webhooksRequestsKey[0] &&
        query.queryKey[1] === webhooksRequestsKey[1] &&
        query.queryKey[2] === webhooksRequestsKey[2] &&
        query.queryKey[3] !== 'detail',
    },
    (prev) => {
      if (!isListResponse(prev)) return prev
      return {
        ...prev,
        items: [],
        activeCount: result.activeCount,
        lifetimeReceived: result.lifetimeReceived,
        total: 0,
        hasMore: false,
      }
    },
  )
}

export function useWebhookInbox(
  options?: Partial<UseQueryOptions<WebhookInbox>>,
) {
  return useQuery({
    queryKey: webhooksInboxKey,
    queryFn: () => webhooksApi.getInbox(),
    refetchInterval: POLL_MS,
    refetchIntervalInBackground: false,
    meta: { skipNavLoading: true },
    ...options,
  })
}

export function useWebhookRequests(
  params: ListWebhookRequestsParams,
  options?: Partial<UseQueryOptions<WebhookRequestListResponse>>,
) {
  return useQuery({
    queryKey: [...webhooksRequestsKey, params],
    queryFn: () => webhooksApi.listRequests(params),
    refetchInterval: POLL_MS,
    refetchIntervalInBackground: false,
    meta: { skipNavLoading: true },
    ...options,
  })
}

export function useWebhookRequestDetail(
  id: string | null,
  options?: Partial<UseQueryOptions<WebhookRequestDetail>>,
) {
  const queryClient = useQueryClient()
  return useQuery({
    queryKey: [...webhooksRequestsKey, 'detail', id],
    queryFn: async () => {
      const detail = await webhooksApi.getRequest(id!)
      // GET marks read server-side; mirror into list caches immediately.
      patchListItemRead(queryClient, detail.id, true)
      return detail
    },
    enabled: Boolean(id),
    meta: { skipNavLoading: true },
    ...options,
  })
}

export function useWebhookInboxMutations() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('common')

  const invalidateAll = () => {
    void queryClient.invalidateQueries({ queryKey: webhooksInboxKey })
    void queryClient.invalidateQueries({ queryKey: webhooksRequestsKey })
  }

  const regenerate = useMutation({
    mutationFn: () => webhooksApi.regenerateInbox(),
    onSuccess: (inbox) => {
      queryClient.setQueryData(webhooksInboxKey, inbox)
      toast.success(t('tools.webhooks.toasts.regenerated'))
    },
    onError: () => {
      toast.error(t('tools.webhooks.toasts.regenerateFailed'))
    },
  })

  const softDelete = useMutation({
    mutationFn: (id: string) => webhooksApi.softDeleteRequest(id),
    onSuccess: (result, deletedId) => {
      applySoftDeleteCounters(queryClient, deletedId, result)
      invalidateAll()
      toast.success(t('tools.webhooks.toasts.removed'))
    },
    onError: () => {
      toast.error(t('tools.webhooks.toasts.removeFailed'))
    },
  })

  const clearAll = useMutation({
    mutationFn: () => webhooksApi.clearAllRequests(),
    onSuccess: (result) => {
      applyClearAllCounters(queryClient, result)
      invalidateAll()
      toast.success(t('tools.webhooks.toasts.cleared'))
    },
    onError: () => {
      toast.error(t('tools.webhooks.toasts.clearFailed'))
    },
  })

  const setRead = useMutation({
    mutationFn: ({ id, isRead }: { id: string; isRead: boolean }) =>
      webhooksApi.setRequestRead(id, isRead),
    onSuccess: (detail) => {
      queryClient.setQueryData(
        [...webhooksRequestsKey, 'detail', detail.id],
        detail,
      )
      patchListItemRead(queryClient, detail.id, detail.isRead)
    },
    onError: () => {
      toast.error(t('tools.webhooks.toasts.readFailed'))
    },
  })

  return { regenerate, softDelete, clearAll, setRead, invalidateAll }
}
