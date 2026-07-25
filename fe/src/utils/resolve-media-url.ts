import { env } from '@/config/env'

/**
 * Resolve avatar/media paths from the API to an absolute URL for <img>/Avatar.
 * BE stores paths like `/api/media/avatars/<uuid>.<ext>` (site-origin relative).
 */
export function resolveMediaUrl(
  path: string | null | undefined,
): string | undefined {
  if (!path) return undefined
  const trimmed = path.trim()
  if (!trimmed) return undefined
  if (/^https?:\/\//i.test(trimmed)) return trimmed

  const siteOrigin = env.VITE_SITE_URL.replace(/\/$/, '')
  if (trimmed.startsWith('/')) {
    return `${siteOrigin}${trimmed}`
  }

  const apiBase = env.VITE_API_BASE_URL.replace(/\/$/, '')
  return `${apiBase}/${trimmed}`
}
