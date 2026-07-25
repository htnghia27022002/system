import { apiClient } from '@/services/api-client'

import type {
  ClearAllWebhookRequestsResponse,
  ListWebhookRequestsParams,
  SoftDeleteWebhookRequestResponse,
  WebhookInbox,
  WebhookRequestDetail,
  WebhookRequestListResponse,
} from '../types/webhooks'
import { normalizeMethodFilter } from '../utils/webhooks-url'

export const webhooksApi = {
  /** Get-or-create the signed-in user's inbox. */
  async getInbox(): Promise<WebhookInbox> {
    const { data } = await apiClient.get<WebhookInbox>('/webhooks/inbox')
    return data
  },

  /** Issue a new public UUID; prior history remains on the inbox. */
  async regenerateInbox(): Promise<WebhookInbox> {
    const { data } = await apiClient.post<WebhookInbox>(
      '/webhooks/inbox/regenerate',
    )
    return data
  },

  /** List active (non-soft-deleted) requests, newest first. */
  async listRequests(
    params: ListWebhookRequestsParams = {},
  ): Promise<WebhookRequestListResponse> {
    const read =
      params.read && params.read !== 'all' ? params.read : undefined
    const { data } = await apiClient.get<WebhookRequestListResponse>(
      '/webhooks/inbox/requests',
      {
        params: {
          method: normalizeMethodFilter(params.method) || undefined,
          q: params.q?.trim() || undefined,
          read,
          page: params.page ?? 1,
          limit: params.limit ?? 20,
        },
      },
    )
    return data
  },

  async getRequest(id: string): Promise<WebhookRequestDetail> {
    const { data } = await apiClient.get<WebhookRequestDetail>(
      `/webhooks/inbox/requests/${id}`,
    )
    return data
  },

  async setRequestRead(
    id: string,
    isRead: boolean,
  ): Promise<WebhookRequestDetail> {
    const { data } = await apiClient.patch<WebhookRequestDetail>(
      `/webhooks/inbox/requests/${id}/read`,
      { isRead },
    )
    return data
  },

  async softDeleteRequest(
    id: string,
  ): Promise<SoftDeleteWebhookRequestResponse> {
    const { data } = await apiClient.delete<SoftDeleteWebhookRequestResponse>(
      `/webhooks/inbox/requests/${id}`,
    )
    return data
  },

  /** Soft-delete all active inbox requests in one call (not N+1 by id). */
  async clearAllRequests(): Promise<ClearAllWebhookRequestsResponse> {
    const { data } = await apiClient.delete<ClearAllWebhookRequestsResponse>(
      '/webhooks/inbox/requests',
    )
    return data
  },
}
