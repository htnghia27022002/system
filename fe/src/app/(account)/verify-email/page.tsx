import { Suspense } from 'react'
import type { Metadata } from 'next'

import { Spinner } from '@/components/ui/spinner'
import { VerifyEmailPage } from '@/features/auth'

export const metadata: Metadata = {
  title: 'Verify email',
  robots: { index: false, follow: false },
}

export default function Page() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-svh items-center justify-center">
          <Spinner className="size-8" />
        </div>
      }
    >
      <VerifyEmailPage />
    </Suspense>
  )
}
