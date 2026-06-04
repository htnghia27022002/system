import type { Metadata } from 'next'

import { AdminRolesPage } from '@/screens/admin-roles-page'
import { PermissionGuard } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Roles',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="roles:read">
      <AdminRolesPage />
    </PermissionGuard>
  )
}
