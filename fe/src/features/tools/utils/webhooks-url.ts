const UUID_RE =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

/** True when value looks like a UUID v4-shaped public webhook id. */
export function isWebhookPublicUuid(value: string): boolean {
  return UUID_RE.test(value.trim())
}

/** Product capture path shown to owners: `/tools/webhooks/{uuid}`. */
export function buildWebhookPublicPath(publicUuid: string): string {
  return `/tools/webhooks/${publicUuid.trim()}`
}

/**
 * Absolute public capture URL for clipboard / display.
 * Accepts either `publicPath` (`/tools/webhooks/...`) or a raw UUID.
 */
export function buildWebhookPublicUrl(
  origin: string,
  publicPathOrUuid: string,
): string {
  const trimmed = publicPathOrUuid.trim()
  const path = trimmed.startsWith('/')
    ? trimmed
    : buildWebhookPublicPath(trimmed)
  const base = origin.replace(/\/$/, '')
  return `${base}${path}`
}

/** Normalize method filter for list API; empty string means no filter. */
export function normalizeMethodFilter(method: string | undefined): string {
  if (!method) return ''
  const m = method.trim().toUpperCase()
  if (!m || m === 'ALL') return ''
  return m
}

/**
 * Derive BE capture destination base from `NEXT_PUBLIC_API_BASE_URL`
 * (already includes `/api`). Used by Next rewrite config and tests.
 */
export function captureApiDestination(
  apiBaseUrl: string,
  uuidParam = ':uuid',
): string {
  const base = apiBaseUrl.replace(/\/$/, '')
  return `${base}/webhooks/capture/${uuidParam}`
}
