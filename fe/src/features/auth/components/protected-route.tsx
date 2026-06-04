'use client'

import { useEffect, type ReactNode } from 'react'
import { usePathname, useRouter } from 'next/navigation'

import { useAuthStore } from '@/store/auth-store'

type ProtectedGuardProps = {
  children: ReactNode
}

export function ProtectedGuard({ children }: ProtectedGuardProps) {
  const accessToken = useAuthStore((state) => state.accessToken)
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (!accessToken) {
      router.replace(`/login?from=${encodeURIComponent(pathname)}`)
    }
  }, [accessToken, router, pathname])

  if (!accessToken) return null

  return <>{children}</>
}

/** @deprecated Use ProtectedGuard in layouts. Kept for auth/index.ts export compatibility. */
export function ProtectedRoute() {
  return null
}
