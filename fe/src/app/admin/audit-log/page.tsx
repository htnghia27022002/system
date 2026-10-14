import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys } from '@/features/access-control'
import { AuditLogTable } from '@/features/system'

export const metadata: Metadata = {
  title: 'Audit log',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.audit.view}>
      <AuditLogTable />
    </PermissionGuard>
  )
}
