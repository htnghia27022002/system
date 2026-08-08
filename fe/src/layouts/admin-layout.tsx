'use client'

import type { ReactNode } from 'react'
import { usePathname } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'

import { AdminMobileNav } from '@/components/common/admin-mobile-nav'
import { AppContent } from '@/components/common/app-content'
import { AppShell } from '@/components/common/app-shell'
import { AppSidebar } from '@/components/common/app-sidebar'
import { AppSidebarHeader } from '@/components/common/app-sidebar-header'
import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { TooltipProvider } from '@/components/ui/tooltip'
import type { BreadcrumbItem } from '@/types/navigation'

const BREADCRUMB_MAP: Record<
  string,
  { label: string; parent?: { label: string; href?: string } }
> = {
  '/admin': { label: 'nav.dashboard' },
  '/admin/search': { label: 'nav.search' },
  '/admin/users': {
    label: 'nav.users',
    parent: { label: 'nav.accessControl', href: '/admin/users' },
  },
  '/admin/roles': {
    label: 'nav.roles',
    parent: { label: 'nav.accessControl', href: '/admin/users' },
  },
  '/admin/tools/webhooks': {
    label: 'nav.webhooks',
    parent: { label: 'nav.tools' },
  },
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
      trail.push({
        title: t(entry.parent.label),
        href: entry.parent.href,
      })
    }

    trail.push({ title: t(entry.label) })
    return trail
  }, [pathname, t])

  return (
    <TooltipProvider delayDuration={0}>
      <AppShell variant="sidebar">
        <AppSidebar />
        <AdminMobileNav />
        <AppContent
          variant="sidebar"
          className="overflow-x-hidden pb-[calc(3rem+env(safe-area-inset-bottom,0px))] md:pb-0"
        >
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
