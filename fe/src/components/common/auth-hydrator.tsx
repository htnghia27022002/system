'use client'

import { useEffect, type ReactNode } from 'react'

import { useAuthStore } from '@/store/auth-store'

export function AuthHydrator({ children }: { children: ReactNode }) {
  const hydrateFromStorage = useAuthStore((state) => state.hydrateFromStorage)
  const syncSession = useAuthStore((state) => state.syncSession)

  useEffect(() => {
    hydrateFromStorage()
    void syncSession()
  }, [hydrateFromStorage, syncSession])

  return <>{children}</>
}
