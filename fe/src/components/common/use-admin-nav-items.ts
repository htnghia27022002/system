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

import { usePermissions, PermissionKeys } from '@/features/access-control'
import type { NavItem } from '@/types/navigation'

export const ADMIN_HOME_HREF = '/admin'

/**
 * Flatten nested admin nav into leaf destinations (for mobile pins / menu grid).
 * Parent group nodes with children are skipped; home/dashboard is a leaf.
 */
export function flattenAdminNavLeaves(items: NavItem[]): NavItem[] {
  const leaves: NavItem[] = []
  for (const item of items) {
    if (item.items && item.items.length > 0) {
      leaves.push(...item.items)
      continue
    }
    leaves.push(item)
  }
  return leaves
}

/**
 * Permission-filtered admin nav tree shared by desktop sidebar and mobile menu.
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

/** Leaf destinations the user can open / pin (excludes empty groups). */
export function useAdminNavLeaves(): NavItem[] {
  const items = useAdminNavItems()
  return useMemo(() => flattenAdminNavLeaves(items), [items])
}
