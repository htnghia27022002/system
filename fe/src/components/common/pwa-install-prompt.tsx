'use client'

import { Download, Share, X } from 'lucide-react'
import { useEffect, useState } from 'react'

import { Button } from '@/components/ui/button'
import { platform } from '@/utils/platform'

const DISMISS_KEY = 'pwa-install-dismissed'

interface BeforeInstallPromptEvent extends Event {
  prompt(): Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

export function PwaInstallPrompt() {
  const [deferredPrompt, setDeferredPrompt] =
    useState<BeforeInstallPromptEvent | null>(null)
  const [showIosHint, setShowIosHint] = useState(false)
  const [dismissed, setDismissed] = useState(true)

  useEffect(() => {
    if (platform.isStandalone) return

    try {
      if (localStorage.getItem(DISMISS_KEY) === '1') return
    } catch {
      // Private browsing may block localStorage.
    }

    setDismissed(false)

    if (platform.isIos) {
      setShowIosHint(true)
      return
    }

    const onBeforeInstall = (event: Event) => {
      event.preventDefault()
      setDeferredPrompt(event as BeforeInstallPromptEvent)
    }

    window.addEventListener('beforeinstallprompt', onBeforeInstall)
    return () =>
      window.removeEventListener('beforeinstallprompt', onBeforeInstall)
  }, [])

  function dismiss() {
    setDismissed(true)
    try {
      localStorage.setItem(DISMISS_KEY, '1')
    } catch {
      // Ignore storage errors.
    }
  }

  async function install() {
    if (!deferredPrompt) return
    await deferredPrompt.prompt()
    await deferredPrompt.userChoice
    setDeferredPrompt(null)
    dismiss()
  }

  if (dismissed || platform.isStandalone) return null
  if (!showIosHint && !deferredPrompt) return null

  return (
    <div
      role="status"
      className="fixed inset-x-0 bottom-0 z-50 border-t bg-background/95 p-4 pb-[calc(1rem+env(safe-area-inset-bottom,0px))] backdrop-blur supports-[backdrop-filter]:bg-background/80"
    >
      <div className="mx-auto flex max-w-7xl items-start gap-3">
        <div className="min-w-0 flex-1">
          <p className="text-sm font-medium">Install System App</p>
          <p className="text-muted-foreground mt-1 text-xs">
            {showIosHint ? (
              <>
                Tap <Share className="inline size-3.5 align-text-bottom" /> Share,
                then &quot;Add to Home Screen&quot; for a native-like experience.
              </>
            ) : (
              'Add this app to your home screen for faster access and offline support.'
            )}
          </p>
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {deferredPrompt ? (
            <Button size="sm" onClick={install}>
              <Download className="size-4" />
              Install
            </Button>
          ) : null}
          <Button
            variant="ghost"
            size="icon-sm"
            aria-label="Dismiss install prompt"
            onClick={dismiss}
          >
            <X className="size-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
