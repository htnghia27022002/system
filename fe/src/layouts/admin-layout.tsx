'use client'

import type { ReactNode } from 'react'
import { usePathname } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'

import { AppContent } from '@/components/common/app-content'
import { AppShell } from '@/components/common/app-shell'
import { AppSidebar } from '@/components/common/app-sidebar'
import { AppSidebarHeader } from '@/components/common/app-sidebar-header'
import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { TooltipProvider } from '@/components/ui/tooltip'
import type { BreadcrumbItem } from '@/types/navigation'

const BREADCRUMB_MAP: Record<
  string,
  { label: string; parent?: string }
> = {
  '/admin': { label: 'nav.dashboard' },
  '/admin/users': { label: 'nav.users', parent: 'nav.accessControl' },
  '/admin/roles': { label: 'nav.roles', parent: 'nav.accessControl' },
}

type AdminLayoutProps = {
  children: ReactNode
}

export function AdminLayout({ children }: AdminLayoutProps) {
  const { t } = useTranslation('admin')
  const pathname = usePathname()

  const breadcrumbs = useMemo((): BreadcrumbItem[] => {
    const entry = BREADCRUMB_MAP[pathname]
    if (!entry) return [{ title: t('nav.dashboard'), href: '/admin' }]
    if (pathname === '/admin') return [{ title: t(entry.label) }]

    const trail: BreadcrumbItem[] = [
      { title: t('nav.dashboard'), href: '/admin' },
    ]

    if (entry.parent) {
      trail.push({ title: t(entry.parent) })
    }

    trail.push({ title: t(entry.label) })
    return trail
  }, [pathname, t])

  return (
    <TooltipProvider delayDuration={0}>
      <AppShell variant="sidebar">
        <AppSidebar />
        <AppContent variant="sidebar" className="overflow-x-hidden">
          <div className="relative">
            <NavLoadingBar />
            <AppSidebarHeader breadcrumbs={breadcrumbs} />
          </div>
          {children}
        </AppContent>
      </AppShell>
    </TooltipProvider>
  )
}
