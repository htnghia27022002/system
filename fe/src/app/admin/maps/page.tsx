import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys } from '@/features/access-control'
import { MapsPage } from '@/features/maps'

export const metadata: Metadata = {
  title: 'Maps',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.maps.view}>
      <MapsPage />
    </PermissionGuard>
  )
}
