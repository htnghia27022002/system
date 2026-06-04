import type { Metadata } from 'next'

import { AdminOverviewPage } from '@/screens/admin-overview-page'
import { PermissionGuard } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Dashboard',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="dashboard:read">
      <AdminOverviewPage />
    </PermissionGuard>
  )
}
