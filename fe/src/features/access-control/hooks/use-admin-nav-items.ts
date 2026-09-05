'use client'

import {
  LayoutGridIcon,
  Link2Icon,
  ShieldCheckIcon,
  UsersIcon,
  WrenchIcon,
} from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { PermissionKeys } from '../permission-keys'
import { usePermissions } from './use-permissions'
import type { NavItem } from '@/types/navigation'

export const ADMIN_HOME_HREF = '/admin'

/**
 * Permission-filtered admin nav tree for the sidebar.
 */
export function useAdminNavItems(): NavItem[] {
  const { t } = useTranslation('admin')
  const { hasPermission } = usePermissions()

  return useMemo(() => {
    const accessControlItems: NavItem[] = [
      {
        title: t('nav.users'),
        href: '/admin/users',
        icon: UsersIcon,
        permission: PermissionKeys.users.view,
      },
      {
        title: t('nav.roles'),
        href: '/admin/roles',
        icon: ShieldCheckIcon,
        permission: PermissionKeys.roles.view,
      },
    ].filter((item) => !item.permission || hasPermission(item.permission))

    const toolsItems: NavItem[] = [
      {
        title: t('nav.webhooks'),
        href: '/admin/tools/webhooks',
        icon: Link2Icon,
        permission: PermissionKeys.webhooks.view,
      },
    ].filter((item) => !item.permission || hasPermission(item.permission))

    const items: NavItem[] = [
      {
        title: t('nav.dashboard'),
        href: ADMIN_HOME_HREF,
        icon: LayoutGridIcon,
        permission: PermissionKeys.dashboard.view,
      },
    ]

    if (toolsItems.length > 0) {
      items.push({
        title: t('nav.tools'),
        href: toolsItems[0].href,
        icon: WrenchIcon,
        items: toolsItems,
      })
    }

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
}
