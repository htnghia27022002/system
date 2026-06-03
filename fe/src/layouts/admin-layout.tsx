import { Outlet, useMatches } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'

import { AppContent } from '@/components/common/app-content'
import { AppShell } from '@/components/common/app-shell'
import { AppSidebar } from '@/components/common/app-sidebar'
import { AppSidebarHeader } from '@/components/common/app-sidebar-header'
import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { TooltipProvider } from '@/components/ui/tooltip'
import type { BreadcrumbItem, RouteBreadcrumbHandle } from '@/types/navigation'

export function AdminLayout() {
  const { t } = useTranslation('admin')
  const matches = useMatches()

  const breadcrumbs = useMemo((): BreadcrumbItem[] => {
    const items: BreadcrumbItem[] = []

    for (const match of matches) {
      const handle = match.handle as RouteBreadcrumbHandle | undefined
      const key = handle?.breadcrumbKey
      if (!key) continue

      if (handle?.parentBreadcrumbKey && handle.parentHref) {
        const parentExists = items.some(
          (item) => item.href === handle.parentHref,
        )
        if (!parentExists) {
          items.push({
            title: t(handle.parentBreadcrumbKey),
            href: handle.parentHref,
          })
        }
      }

      items.push({
        title: t(key),
        href: match.pathname,
      })
    }

    if (items.length === 0) {
      return [{ title: t('nav.dashboard'), href: '/admin' }]
    }

    return items
  }, [matches, t])

  return (
    <TooltipProvider delayDuration={0}>
      <AppShell variant="sidebar">
        <AppSidebar />
        <AppContent variant="sidebar" className="overflow-x-hidden">
          <div className="relative">
            <NavLoadingBar />
            <AppSidebarHeader breadcrumbs={breadcrumbs} />
          </div>
          <Outlet />
        </AppContent>
      </AppShell>
    </TooltipProvider>
  )
}
