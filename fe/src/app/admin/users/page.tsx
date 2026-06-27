import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys, UsersTable } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Users',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.users.view}>
      <UsersTable />
    </PermissionGuard>
  )
}
