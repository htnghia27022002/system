'use client'

import { useEffect, type ReactNode } from 'react'
import { useRouter } from 'next/navigation'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

type GuestGuardProps = {
  children: ReactNode
}

export function GuestGuard({ children }: GuestGuardProps) {
  const hasHydrated = useAuthStore((state) => state.hasHydrated)
  const user = useAuthStore((state) => state.user)
  const router = useRouter()

  useEffect(() => {
    if (hasHydrated && user) {
      router.replace(getPostLoginPath(user.role))
    }
  }, [hasHydrated, user, router])

  if (hasHydrated && user) return null

  return <>{children}</>
}

/** @deprecated Use GuestGuard in layouts. Kept for auth/index.ts export compatibility. */
export function GuestRoute() {
  return null
}
