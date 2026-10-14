'use client'

import { ErrorFallback } from '@/components/common/error-fallback'
import { PageContainer } from '@/components/common/page-layout'

/** Keeps the admin sidebar/header usable when a page inside it throws. */
export default function AdminError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <PageContainer>
      <ErrorFallback
        error={error}
        reset={reset}
        homeHref="/admin"
        className="flex-1"
      />
    </PageContainer>
  )
}
