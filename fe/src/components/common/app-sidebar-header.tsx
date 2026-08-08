'use client'

import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { Separator } from '@/components/ui/separator'
import { SidebarTrigger, useSidebar } from '@/components/ui/sidebar'
import type { BreadcrumbItem } from '@/types/navigation'

import { Breadcrumbs } from './breadcrumbs'

type AppSidebarHeaderProps = {
  breadcrumbs?: BreadcrumbItem[]
  actions?: ReactNode
}

export function AppSidebarHeader({ breadcrumbs = [], actions }: AppSidebarHeaderProps) {
  const { t } = useTranslation('admin')
  const { isMobile } = useSidebar()

  return (
    <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 max-md:min-h-[calc(4rem+var(--safe-top))] max-md:pt-[var(--safe-top)] max-md:px-[calc(1rem+var(--safe-left))] max-md:pr-[calc(1rem+var(--safe-right))]">
      <div className="flex min-w-0 flex-1 items-center gap-2">
        {!isMobile ? (
          <>
            <SidebarTrigger className="-ml-1" aria-label={t('shell.openMenu')} />
          </>
        ) : null}
        <Breadcrumbs breadcrumbs={breadcrumbs} />
      </div>
      {actions ? (
        <div className="ml-auto flex shrink-0 items-center justify-end">{actions}</div>
      ) : null}
    </header>
  )
}
