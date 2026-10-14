import type { Metadata } from 'next'

import { ForgotPasswordPage } from '@/features/auth'

export const metadata: Metadata = {
  title: 'Forgot password',
  robots: { index: false, follow: false },
}

export default function Page() {
  return <ForgotPasswordPage />
}
