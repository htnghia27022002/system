'use client'

import { useCallback } from 'react'
import { usePathname } from 'next/navigation'

export function useCurrentPath() {
  const pathname = usePathname()

  const isCurrentPath = useCallback(
    (path: string) => pathname === path,
    [pathname],
  )

  return { pathname, isCurrentPath }
}
