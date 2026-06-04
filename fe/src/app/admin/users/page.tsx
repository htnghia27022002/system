import type { Metadata } from 'next'

import { PermissionGuard, UsersTable } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Users',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="users:read">
      <UsersTable />
    </PermissionGuard>
  )
}
