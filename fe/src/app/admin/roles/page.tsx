import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys, RolesTable } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Roles',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.roles.view}>
      <RolesTable />
    </PermissionGuard>
  )
}
