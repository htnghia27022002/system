'use client'

import type { ReactNode } from 'react'
import { useMemo } from 'react'
import { usePathname } from 'next/navigation'
import { useTranslation } from 'react-i18next'

import { AppSidebar } from '@/components/common/app-sidebar'
import { Breadcrumbs } from '@/components/common/breadcrumbs'
import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { Separator } from '@/components/ui/separator'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'
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
    <SidebarProvider defaultOpen>
      <AppSidebar />
      <SidebarInset className="overflow-x-hidden">
        <div className="relative">
          <NavLoadingBar />
          <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
            <SidebarTrigger className="-ml-1" />
            <Separator
              orientation="vertical"
              className="mr-2 data-vertical:h-4 data-vertical:self-auto"
            />
            <Breadcrumbs breadcrumbs={breadcrumbs} />
          </header>
        </div>
        {children}
      </SidebarInset>
    </SidebarProvider>
  )
}
