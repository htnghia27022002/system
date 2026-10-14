import type { Metadata } from 'next'

import { PermissionGuard, PermissionKeys } from '@/features/access-control'
import { SettingsPage } from '@/features/system'

export const metadata: Metadata = {
  title: 'Settings',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <PermissionGuard permission={PermissionKeys.settings.view}>
      <SettingsPage />
    </PermissionGuard>
  )
}
