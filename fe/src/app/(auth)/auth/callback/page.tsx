import { Suspense } from 'react'
import type { Metadata } from 'next'

import { OAuthCallbackPage } from '@/features/auth/components/oauth-callback-page'
import { Spinner } from '@/components/ui/spinner'

export const metadata: Metadata = {
  title: 'Signing in',
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
      <OAuthCallbackPage />
    </Suspense>
  )
}
