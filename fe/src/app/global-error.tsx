'use client'

import { useEffect } from 'react'

import { Button } from '@/components/ui/button'

import '@/styles/index.css'

/**
 * Last-resort boundary: replaces the root layout when it throws, so providers
 * (i18n, theme) are unavailable — copy is static English on purpose.
 */
export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  useEffect(() => {
    console.error(error)
  }, [error])

  return (
    <html lang="en">
      <body className="bg-background text-foreground">
        <div
          role="alert"
          className="flex min-h-[100dvh] flex-col items-center justify-center gap-4 p-6 text-center"
        >
          <h1 className="text-xl font-semibold tracking-tight">
            Something went wrong
          </h1>
          <p className="text-sm text-muted-foreground">
            The application failed to load. Please try again.
          </p>
          {error.digest ? (
            <p className="font-mono text-xs text-muted-foreground">
              Reference: {error.digest}
            </p>
          ) : null}
          <Button onClick={reset}>Try again</Button>
        </div>
      </body>
    </html>
  )
}
