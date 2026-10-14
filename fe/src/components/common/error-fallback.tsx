'use client'

import { AlertTriangle } from 'lucide-react'
import Link from 'next/link'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type ErrorFallbackProps = {
  error: Error & { digest?: string }
  /** Re-renders the failed segment (Next.js error boundary `reset`). */
  reset: () => void
  /** Where the secondary "go back" link points. */
  homeHref?: string
  className?: string
}

/**
 * Shared body for App Router `error.tsx` boundaries. Shows a generic message —
 * never `error.message`, which may carry server details — plus the digest so a
 * user report can be matched to server logs.
 */
export function ErrorFallback({
  error,
  reset,
  homeHref = '/',
  className,
}: ErrorFallbackProps) {
  const { t } = useTranslation('common')

  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <div
      role="alert"
      className={cn(
        'flex flex-col items-center justify-center gap-4 p-6 text-center',
        className,
      )}
    >
      <AlertTriangle className="size-10 text-destructive" aria-hidden />
      <div className="space-y-1">
        <h1 className="text-xl font-semibold tracking-tight">
          {t('errorBoundary.title')}
        </h1>
        <p className="text-sm text-muted-foreground">
          {t('errorBoundary.description')}
        </p>
        {error.digest ? (
          <p className="font-mono text-xs text-muted-foreground">
            {t('errorBoundary.reference', { digest: error.digest })}
          </p>
        ) : null}
      </div>
      <div className="flex flex-col gap-2 sm:flex-row">
        <Button onClick={reset}>{t('errorBoundary.retry')}</Button>
        <Button variant="outline" asChild>
          <Link href={homeHref}>{t('errorBoundary.home')}</Link>
        </Button>
      </div>
    </div>
  )
}
