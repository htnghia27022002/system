'use client'

import type { ReactNode } from 'react'

import { SidebarProvider } from '@/components/ui/sidebar'

type AppShellProps = {
  children: ReactNode
  variant?: 'header' | 'sidebar'
}

export function AppShell({ children, variant = 'header' }: AppShellProps) {
  if (variant === 'header') {
    return <div className="flex min-h-svh w-full flex-col">{children}</div>
  }

  return <SidebarProvider defaultOpen>{children}</SidebarProvider>
}
