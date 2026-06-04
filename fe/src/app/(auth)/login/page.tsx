import type { Metadata } from 'next'

import { LoginPage } from '@/features/auth'

export const metadata: Metadata = {
  title: 'Sign In',
  description: 'Sign in to your account.',
  robots: { index: false, follow: false },
}

export default function Page() {
  return <LoginPage />
}
