import type { Metadata } from 'next'

import { PermissionGuard } from '@/features/access-control'
import { AdminDashboardOverview } from '@/features/admin-dashboard'

export const metadata: Metadata = {
  title: 'Dashboard',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="dashboard:read">
      <AdminDashboardOverview />
    </PermissionGuard>
  )
}
