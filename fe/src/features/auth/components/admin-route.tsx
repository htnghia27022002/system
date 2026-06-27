'use client'

import type { ReactNode } from 'react'

type AdminGuardProps = {
  children: ReactNode
}

/** @deprecated Admin area access is gated by ProtectedGuard + per-page PermissionGuard. */
export function AdminGuard({ children }: AdminGuardProps) {
  return <>{children}</>
}

/** @deprecated Use ProtectedGuard in layouts. Kept for auth/index.ts export compatibility. */
export function AdminRoute() {
  return null
}
