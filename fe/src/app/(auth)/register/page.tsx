import type { Metadata } from 'next'

import { RegisterPage } from '@/screens/register-page'

export const metadata: Metadata = {
  title: 'Create Account',
  description: 'Create a new account.',
  robots: { index: false, follow: false },
}

export default function Page() {
  return <RegisterPage />
}
