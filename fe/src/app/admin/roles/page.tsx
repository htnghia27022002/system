import type { Metadata } from 'next'

import { PermissionGuard, RolesTable } from '@/features/access-control'

export const metadata: Metadata = {
  title: 'Roles',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission="roles:read">
      <RolesTable />
    </PermissionGuard>
  )
}
