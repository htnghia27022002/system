import type { Metadata } from 'next'

import { ProfilePage } from '@/features/user-profile'

export const metadata: Metadata = {
  title: 'Account settings',
  robots: { index: false, follow: false },
}

export default function Page() {
  return <ProfilePage />
}
