'use client'

import type { ReactNode } from 'react'

import { ProtectedGuard } from '@/features/auth'
import { AdminLayout } from '@/layouts/admin-layout'

export default function AdminGroupLayout({ children }: { children: ReactNode }) {
  return (
    <ProtectedGuard>
      <AdminLayout>{children}</AdminLayout>
    </ProtectedGuard>
  )
}
