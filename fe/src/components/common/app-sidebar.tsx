import { Link } from 'react-router-dom'
import { LayoutGridIcon, ShieldIcon, UsersIcon } from 'lucide-react'
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
import { usePermissions } from '@/features/access-control'
import type { NavItem } from '@/types/navigation'

export function AppSidebar() {
  const { t } = useTranslation('admin')
  const { hasPermission } = usePermissions()

  const mainNavItems: NavItem[] = useMemo(
    () =>
      [
        {
          title: t('nav.dashboard'),
          href: '/admin',
          icon: LayoutGridIcon,
          permission: 'dashboard:read',
        },
        {
          title: t('nav.users'),
          href: '/admin/users',
          icon: UsersIcon,
          permission: 'users:read',
        },
        {
          title: t('nav.roles'),
          href: '/admin/roles',
          icon: ShieldIcon,
          permission: 'roles:read',
        },
      ].filter(
        (item) => !item.permission || hasPermission(item.permission),
      ),
    [hasPermission, t],
  )

  return (
    <Sidebar collapsible="icon" variant="inset">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to="/admin">
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
