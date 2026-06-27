'use client'

import { useEffect, type ReactNode } from 'react'
import { useRouter } from 'next/navigation'

import { useAuthStore } from '@/store/auth-store'

type AdminGuardProps = {
  children: ReactNode
}

export function AdminGuard({ children }: AdminGuardProps) {
  const hasHydrated = useAuthStore((state) => state.hasHydrated)
  const isAdmin = useAuthStore((state) => state.isAdmin)
  const router = useRouter()

  useEffect(() => {
    if (hasHydrated && !isAdmin) {
      router.replace('/')
    }
  }, [hasHydrated, isAdmin, router])

  if (!hasHydrated || !isAdmin) return null

  return <>{children}</>
}

/** @deprecated Use AdminGuard in layouts. Kept for auth/index.ts export compatibility. */
export function AdminRoute() {
  return null
}
