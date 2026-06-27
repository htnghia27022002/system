'use client'

import Link from 'next/link'
import { LayoutGridIcon, ShieldCheckIcon } from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { AdminAppLogo } from '@/components/common/admin-app-logo'
import { NavMain } from '@/components/common/nav-main'
import { NavUser } from '@/components/common/nav-user'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { usePermissions, PermissionKeys } from '@/features/access-control'
import type { NavItem } from '@/types/navigation'

export function AppSidebar() {
  const { t } = useTranslation('admin')
  const { hasPermission } = usePermissions()

  const mainNavItems: NavItem[] = useMemo(() => {
    const accessControlItems: NavItem[] = [
      {
        title: t('nav.users'),
        href: '/admin/users',
        permission: PermissionKeys.users.view,
      },
      {
        title: t('nav.roles'),
        href: '/admin/roles',
        permission: PermissionKeys.roles.view,
      },
    ].filter((item) => !item.permission || hasPermission(item.permission))

    const items: NavItem[] = [
      {
        title: t('nav.dashboard'),
        href: '/admin',
        icon: LayoutGridIcon,
        permission: PermissionKeys.dashboard.view,
      },
    ]

    if (accessControlItems.length > 0) {
      items.push({
        title: t('nav.accessControl'),
        href: accessControlItems[0].href,
        icon: ShieldCheckIcon,
        items: accessControlItems,
      })
    }

    return items.filter(
      (item) => !item.permission || hasPermission(item.permission),
    )
  }, [hasPermission, t])

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
