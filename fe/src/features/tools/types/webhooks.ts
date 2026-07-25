/** Body encoding from GET /webhooks/inbox/requests/:id (contracts/endpoints.md). */
export type WebhookBodyEncoding = 'utf-8' | 'base64'

/** Capture outcome for a stored request row. */
export type WebhookCaptureStatus = 'ok' | 'oversized' | 'error'

/** Owner inbox summary from GET /webhooks/inbox (camelCase). */
export type WebhookInbox = {
  id: string
  publicUuid: string
  publicPath: string
  activeCount: number
  lifetimeReceived: number
  createdAt: string
  updatedAt: string
}

/** List row from GET /webhooks/inbox/requests. */
export type WebhookRequestListItem = {
  id: string
  method: string
  url: string
  clientIp: string
  createdAt: string
  snippet: string
  isRead: boolean
}

/** Paginated list response. */
export type WebhookRequestListResponse = {
  items: WebhookRequestListItem[]
  activeCount: number
  lifetimeReceived: number
  page: number
  limit: number
  total: number
  hasMore: boolean
}

/** Detail from GET /webhooks/inbox/requests/:id. */
export type WebhookRequestDetail = {
  id: string
  inboxId: string
  method: string
  url: string
  clientIp: string
  headers: Record<string, unknown>
  query: Record<string, unknown>
  form: Record<string, unknown>
  body: string
  /** `utf-8` (text) or `base64` (binary); when binary, `isBinary` is true. */
  bodyEncoding: WebhookBodyEncoding | string
  isBinary: boolean
  contentType: string
  bodyTruncated: boolean
  captureStatus: WebhookCaptureStatus | string
  isRead: boolean
  createdAt: string
}

export type SoftDeleteWebhookRequestResponse = {
  ok: boolean
  activeCount: number
  lifetimeReceived: number
}

/** Response from DELETE /webhooks/inbox/requests (clear all active). */
export type ClearAllWebhookRequestsResponse = SoftDeleteWebhookRequestResponse

export type WebhookReadFilter = 'all' | 'read' | 'unread'

export type ListWebhookRequestsParams = {
  method?: string
  q?: string
  read?: WebhookReadFilter
  page?: number
  limit?: number
}

/** HTTP methods commonly shown in the method filter. */
export const WEBHOOK_METHOD_FILTERS = [
  'ALL',
  'GET',
  'POST',
  'PUT',
  'PATCH',
  'DELETE',
  'HEAD',
  'OPTIONS',
] as const

export type WebhookMethodFilter = (typeof WEBHOOK_METHOD_FILTERS)[number]

export const WEBHOOK_READ_FILTERS = ['all', 'unread', 'read'] as const
