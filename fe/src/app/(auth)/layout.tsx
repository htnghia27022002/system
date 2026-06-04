'use client'

import type { ReactNode } from 'react'

import { GuestGuard } from '@/features/auth'

export default function AuthGroupLayout({ children }: { children: ReactNode }) {
  return <GuestGuard>{children}</GuestGuard>
}
