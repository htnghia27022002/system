'use client'

import { useEffect, type ReactNode } from 'react'
import { useRouter } from 'next/navigation'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

type GuestGuardProps = {
  children: ReactNode
}

export function GuestGuard({ children }: GuestGuardProps) {
  const user = useAuthStore((state) => state.user)
  const router = useRouter()

  useEffect(() => {
    if (user) {
      router.replace(getPostLoginPath(user.role))
    }
  }, [user, router])

  if (user) return null

  return <>{children}</>
}

/** @deprecated Use GuestGuard in layouts. Kept for auth/index.ts export compatibility. */
export function GuestRoute() {
  return null
}
