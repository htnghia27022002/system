import { Suspense } from 'react'
import type { Metadata } from 'next'

import { Spinner } from '@/components/ui/spinner'
import { AcceptInvitePage } from '@/features/auth'

export const metadata: Metadata = {
  title: 'Accept invitation',
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
      <AcceptInvitePage />
    </Suspense>
  )
}
