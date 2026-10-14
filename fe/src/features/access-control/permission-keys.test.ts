import { describe, expect, it } from 'vitest'

import { definePermissions } from '@/lib/permissions'

import {
  PermissionKeys,
  canModifyResource,
  canViewResource,
  isPermissionAllowed,
} from './permission-keys'

describe('permission-keys', () => {
  it('matches view and modify keys', () => {
    const granted = new Set([
      PermissionKeys.users.view,
      PermissionKeys.roles.modify,
    ])

    expect(isPermissionAllowed(granted, PermissionKeys.users.view)).toBe(true)
    expect(isPermissionAllowed(granted, PermissionKeys.users.modify)).toBe(false)
    expect(canViewResource(granted, 'users')).toBe(true)
    expect(canModifyResource(granted, 'roles')).toBe(true)
  })

  it('checks module resources defined with definePermissions', () => {
    const widgets = definePermissions('widgets')
    expect(widgets).toEqual({ view: 'widgets:view', modify: 'widgets:modify' })

    const granted = new Set([widgets.view])
    expect(canViewResource(granted, 'widgets')).toBe(true)
    expect(canModifyResource(granted, 'widgets')).toBe(false)
  })

  it('maps legacy read/crud keys during transition', () => {
    const granted = new Set(['dashboard:read', 'users:update'])

    expect(canViewResource(granted, 'dashboard')).toBe(true)
    expect(canModifyResource(granted, 'users')).toBe(true)
  })
})
