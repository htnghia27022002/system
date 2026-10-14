import { renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import '@/config/i18n'

const { granted } = vi.hoisted(() => ({ granted: { keys: new Set<string>() } }))

vi.mock('@/features/access-control/hooks/use-permissions', () => ({
  usePermissions: () => ({ hasPermission: (key: string) => granted.keys.has(key) }),
}))

import { useAdminNavItems } from '@/features/access-control/hooks/use-admin-nav-items'

import {
  appModules,
  navItemsForBackend,
  moduleBreadcrumbs,
  moduleNavItems,
  modulePermissions,
  moduleSitemapPaths,
} from './modules'

describe('module registry', () => {
  it('registers the tools module and collects its contributions', () => {
    // Projects append their own modules; the reference module must stay wired.
    expect(appModules.map((m) => m.id)).toContain('tools')
    expect(modulePermissions.map((p) => p.key)).toEqual(
      expect.arrayContaining(['webhooks:view', 'webhooks:modify']),
    )
    expect(moduleBreadcrumbs.map((b) => b.path)).toContain('/admin/tools/webhooks')
    expect(moduleSitemapPaths).toContain('/tools')
  })

  it('groups module nav entries and hides them without permission', () => {
    granted.keys = new Set(['dashboard:view', 'users:view', 'webhooks:view'])
    const { result: withAccess } = renderHook(() => useAdminNavItems(moduleNavItems))
    const tools = withAccess.current.find((item) => item.title === 'Tools')
    expect(tools?.items?.map((item) => item.href)).toEqual(['/admin/tools/webhooks'])
    expect(withAccess.current.map((item) => item.title)).toEqual([
      'Dashboard',
      'Tools',
      'Access control',
    ])

    granted.keys = new Set(['dashboard:view'])
    const { result: noAccess } = renderHook(() => useAdminNavItems(moduleNavItems))
    expect(noAccess.current.map((item) => item.title)).toEqual(['Dashboard'])
  })

  it('hides module nav when its BE module is off or still unknown', () => {
    expect(navItemsForBackend(new Set(['webhooks'])).map((i) => i.href)).toEqual([
      '/admin/tools/webhooks',
    ])
    expect(navItemsForBackend(new Set())).toEqual([])
    expect(navItemsForBackend(undefined)).toEqual([])
  })
})
