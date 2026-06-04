'use client'

import type { ReactNode } from 'react'

import { AdminGuard, ProtectedGuard } from '@/features/auth'
import { AdminLayout } from '@/layouts/admin-layout'

export default function AdminGroupLayout({ children }: { children: ReactNode }) {
  return (
    <ProtectedGuard>
      <AdminGuard>
        <AdminLayout>{children}</AdminLayout>
      </AdminGuard>
    </ProtectedGuard>
  )
}
