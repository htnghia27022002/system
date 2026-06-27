'use client'

import { useEffect, type ReactNode } from 'react'
import { usePathname, useRouter } from 'next/navigation'

import { useAuthStore } from '@/store/auth-store'

type ProtectedGuardProps = {
  children: ReactNode
}

export function ProtectedGuard({ children }: ProtectedGuardProps) {
  const hasHydrated = useAuthStore((state) => state.hasHydrated)
  const accessToken = useAuthStore((state) => state.accessToken)
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (hasHydrated && !accessToken) {
      router.replace(`/login?from=${encodeURIComponent(pathname)}`)
    }
  }, [accessToken, hasHydrated, router, pathname])

  if (!hasHydrated || !accessToken) return null

  return <>{children}</>
}

/** @deprecated Use ProtectedGuard in layouts. Kept for auth/index.ts export compatibility. */
export function ProtectedRoute() {
  return null
}
