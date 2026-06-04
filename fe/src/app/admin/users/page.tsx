import type { Metadata } from 'next'

import { AdminUsersPage } from '@/screens/admin-users-page'
import { PermissionGuard } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Users',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="users:read">
      <AdminUsersPage />
    </PermissionGuard>
  )
}
