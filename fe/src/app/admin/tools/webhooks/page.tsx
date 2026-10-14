import type { Metadata } from 'next'

import { PermissionGuard } from '@/features/access-control'
import { WebhooksPermissions, WebhooksToolPage } from '@/features/tools'

export const metadata: Metadata = {
  title: 'Webhooks',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={WebhooksPermissions.view}>
      <WebhooksToolPage />
    </PermissionGuard>
  )
}
