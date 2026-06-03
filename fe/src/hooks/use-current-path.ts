import { useCallback } from 'react'
import { useLocation } from 'react-router-dom'

export function useCurrentPath() {
  const { pathname } = useLocation()

  const isCurrentPath = useCallback(
    (path: string) => pathname === path,
    [pathname],
  )

  return { pathname, isCurrentPath }
}
