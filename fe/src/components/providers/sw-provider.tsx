'use client'

import type { ReactNode } from 'react'
import { useEffect } from 'react'
import { toast } from 'sonner'

import { platform } from '@/utils/platform'

type SwProviderProps = {
  children: ReactNode
}

export function SwProvider({ children }: SwProviderProps) {
  useEffect(() => {
    if (!platform.supports.sw || process.env.NODE_ENV !== 'production') return

    let cancelled = false

    async function register() {
      try {
        const registration = await navigator.serviceWorker.register('/sw.js', {
          scope: '/',
        })

        if (cancelled) return

        registration.addEventListener('updatefound', () => {
          const installing = registration.installing
          if (!installing) return

          installing.addEventListener('statechange', () => {
            if (
              installing.state === 'installed' &&
              navigator.serviceWorker.controller
            ) {
              toast('Update available', {
                description: 'A new version is ready.',
                action: {
                  label: 'Reload',
                  onClick: () => window.location.reload(),
                },
                duration: Infinity,
              })
            }
          })
        })
      } catch {
        // SW registration fails on insecure origins or unsupported browsers.
      }
    }

    void register()

    return () => {
      cancelled = true
    }
  }, [])

  return children
}
