'use client'

import { SidebarTrigger } from '@/components/ui/sidebar'
import type { BreadcrumbItem } from '@/types/navigation'

import { Breadcrumbs } from './breadcrumbs'

type AppSidebarHeaderProps = {
  breadcrumbs?: BreadcrumbItem[]
}

export function AppSidebarHeader({ breadcrumbs = [] }: AppSidebarHeaderProps) {
  return (
    <header className="flex h-16 shrink-0 items-center gap-2 border-b border-sidebar-border/50 px-6 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 md:px-4">
      <div className="flex items-center gap-2">
        <SidebarTrigger className="-ml-1" />
        <Breadcrumbs breadcrumbs={breadcrumbs} />
      </div>
    </header>
  )
}
