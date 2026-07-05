'use client'

import type { ReactNode } from 'react'

import { SidebarTrigger } from '@/components/ui/sidebar'
import type { BreadcrumbItem } from '@/types/navigation'

import { Breadcrumbs } from './breadcrumbs'

type AppSidebarHeaderProps = {
  breadcrumbs?: BreadcrumbItem[]
  actions?: ReactNode
}

export function AppSidebarHeader({ breadcrumbs = [], actions }: AppSidebarHeaderProps) {
  return (
    <header className="flex h-16 shrink-0 items-center gap-3 border-b border-sidebar-border/50 px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 md:px-4">
      <div className="flex min-w-0 items-center gap-2">
        <SidebarTrigger className="-ml-1" />
        <Breadcrumbs breadcrumbs={breadcrumbs} />
      </div>
      {actions ? (
        <div className="ml-auto flex shrink-0 items-center justify-end">
          {actions}
        </div>
      ) : null}
    </header>
  )
}
