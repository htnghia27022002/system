import type { ReactNode } from 'react'

/**
 * Owner inbox redirects to `/admin/tools/webhooks` (page.tsx).
 * No ProtectedGuard here — admin layout owns auth + PermissionGuard.
 * UUID capture paths never hit this layout (rewrite → BE).
 */
export default function WebhooksToolLayout({
  children,
}: {
  children: ReactNode
}) {
  return children
}
