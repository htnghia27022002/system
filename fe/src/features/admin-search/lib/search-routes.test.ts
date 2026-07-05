import { describe, expect, it } from 'vitest'

import {
  canNavigateToSearchHit,
  entityTypeLabel,
  searchHitHref,
  type SearchNavigationAccess,
} from './search-routes'

const fullAccess: SearchNavigationAccess = {
  canViewUsers: true,
  canViewRoles: true,
  canViewPermissions: true,
}

const rolesOnlyAccess: SearchNavigationAccess = {
  canViewUsers: false,
  canViewRoles: true,
  canViewPermissions: true,
}

describe('search-routes', () => {
  it('maps entity types to filtered admin list routes', () => {
    expect(
      searchHitHref(
        {
          entityType: 'user',
          entityId: '1',
          title: 'Admin User',
          metadata: { email: 'admin@example.com' },
          updatedAt: '2026-01-01T00:00:00.000Z',
        },
        fullAccess,
      ),
    ).toBe('/admin/users?search=admin%40example.com')

    expect(
      searchHitHref(
        {
          entityType: 'role',
          entityId: '2',
          title: 'Administrator',
          metadata: { slug: 'admin' },
          updatedAt: '2026-01-01T00:00:00.000Z',
        },
        fullAccess,
      ),
    ).toBe('/admin/roles?search=admin')

    expect(
      searchHitHref(
        {
          entityType: 'permission',
          entityId: '3',
          title: 'View users',
          metadata: { key: 'users:view' },
          updatedAt: '2026-01-01T00:00:00.000Z',
        },
        fullAccess,
      ),
    ).toBe('/admin/roles?permissionKey=users%3Aview')
  })

  it('blocks navigation when the caller lacks list-page access', () => {
    const userHit = {
      entityType: 'user' as const,
      entityId: '1',
      title: 'Admin User',
      metadata: { email: 'admin@example.com' },
      updatedAt: '2026-01-01T00:00:00.000Z',
    }

    expect(searchHitHref(userHit, rolesOnlyAccess)).toBeNull()
    expect(canNavigateToSearchHit(userHit, rolesOnlyAccess)).toBe(false)
  })

  it('requires roles view to open permission hits on the roles list', () => {
    const permissionHit = {
      entityType: 'permission' as const,
      entityId: '3',
      title: 'View users',
      metadata: { key: 'users:view' },
      updatedAt: '2026-01-01T00:00:00.000Z',
    }

    expect(
      searchHitHref(permissionHit, {
        canViewUsers: false,
        canViewRoles: false,
        canViewPermissions: true,
      }),
    ).toBeNull()

    expect(searchHitHref(permissionHit, rolesOnlyAccess)).toBe(
      '/admin/roles?permissionKey=users%3Aview',
    )
  })

  it('labels entity types for display', () => {
    expect(entityTypeLabel('user')).toBe('User')
    expect(entityTypeLabel('permission')).toBe('Permission')
  })
})
