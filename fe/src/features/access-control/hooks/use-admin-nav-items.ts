'use client'

import {
  HistoryIcon,
  LayoutGridIcon,
  Settings2Icon,
  ShieldCheckIcon,
  UsersIcon,
} from 'lucide-react'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { PermissionKeys } from '../permission-keys'
import { usePermissions } from './use-permissions'
import type { NavItem } from '@/types/navigation'
import type { ModuleNavItem } from '@/types/modules'

export const ADMIN_HOME_HREF = '/admin'

/**
 * Permission-filtered admin nav tree: Dashboard, then module entries (grouped by
 * `group.key`), then Access control, then System. Module entries come from the caller
 * (`src/config/modules.ts`) so this core feature never imports a module.
 */
export function useAdminNavItems(moduleItems: readonly ModuleNavItem[] = []): NavItem[] {
  const { t } = useTranslation('admin')
  const { hasPermission } = usePermissions()

  return useMemo(() => {
    const allowed = <T extends { permission?: string }>(item: T) =>
      !item.permission || hasPermission(item.permission)

    const items: NavItem[] = [
      {
        title: t('nav.dashboard'),
        href: ADMIN_HOME_HREF,
        icon: LayoutGridIcon,
        permission: PermissionKeys.dashboard.view,
      },
    ]

    const groups = new Map<string, NavItem>()
    for (const entry of moduleItems.filter(allowed)) {
      const item: NavItem = {
        title: t(entry.titleKey),
        href: entry.href,
        icon: entry.icon,
        permission: entry.permission,
      }
      if (!entry.group) {
        items.push(item)
        continue
      }
      const parent = groups.get(entry.group.key)
      if (parent) {
        parent.items = [...(parent.items ?? []), item]
        continue
      }
      const group: NavItem = {
        title: t(entry.group.titleKey),
        href: item.href,
        icon: entry.group.icon,
        items: [item],
      }
      groups.set(entry.group.key, group)
      items.push(group)
    }

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
    ].filter(allowed)

    if (accessControlItems.length > 0) {
      items.push({
        title: t('nav.accessControl'),
        href: accessControlItems[0].href,
        icon: ShieldCheckIcon,
        items: accessControlItems,
      })
    }

    const systemItems: NavItem[] = [
      {
        title: t('nav.auditLog'),
        href: '/admin/audit-log',
        icon: HistoryIcon,
        permission: PermissionKeys.audit.view,
      },
      {
        title: t('nav.settings'),
        href: '/admin/settings',
        icon: Settings2Icon,
        permission: PermissionKeys.settings.view,
      },
    ].filter(allowed)

    if (systemItems.length > 0) {
      items.push({
        title: t('nav.system'),
        href: systemItems[0].href,
        icon: Settings2Icon,
        items: systemItems,
      })
    }

    return items.filter(allowed)
  }, [hasPermission, moduleItems, t])
}
