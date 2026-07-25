import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys } from '@/features/access-control'
import { WebhooksToolPage } from '@/features/tools'

export const metadata: Metadata = {
  title: 'Webhooks',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.webhooks.view}>
      <WebhooksToolPage />
    </PermissionGuard>
  )
}
