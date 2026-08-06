'use client'

import Link from 'next/link'

import { AdminAppLogo } from '@/components/common/admin-app-logo'
import { NavMain } from '@/components/common/nav-main'
import { NavUser } from '@/components/common/nav-user'
import { useAdminNavItems } from '@/components/common/use-admin-nav-items'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'

export function AppSidebar() {
  const { isMobile } = useSidebar()
  const mainNavItems = useAdminNavItems()

  // Mobile uses AdminMobileNav (dedicated touch drawer) instead of the
  // cramped desktop sidebar-in-sheet pattern from shadcn Sidebar.
  if (isMobile) {
    return null
  }

  return (
    <Sidebar collapsible="icon" variant="inset">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link href="/admin">
                <AdminAppLogo />
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <NavMain items={mainNavItems} />
      </SidebarContent>

      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
    </Sidebar>
  )
}
