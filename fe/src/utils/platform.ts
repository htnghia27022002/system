type NavigatorStandalone = Navigator & { standalone?: boolean }

export const platform = {
  isIos:
    typeof navigator !== 'undefined' &&
    /iPad|iPhone|iPod/.test(navigator.userAgent),
  isAndroid:
    typeof navigator !== 'undefined' && /Android/.test(navigator.userAgent),
  isStandalone:
    typeof window !== 'undefined' &&
    (window.matchMedia('(display-mode: standalone)').matches ||
      (navigator as NavigatorStandalone).standalone === true),
  supports: {
    sw: typeof navigator !== 'undefined' && 'serviceWorker' in navigator,
    push: typeof window !== 'undefined' && 'PushManager' in window,
    share: typeof navigator !== 'undefined' && 'share' in navigator,
    vibrate: typeof navigator !== 'undefined' && 'vibrate' in navigator,
    sync: typeof window !== 'undefined' && 'SyncManager' in window,
  },
} as const

/** iOS web push only works in an installed Home Screen PWA (16.4+). */
export function canUsePush(): boolean {
  if (!platform.supports.push) return false
  if (platform.isIos && !platform.isStandalone) return false
  return true
}
