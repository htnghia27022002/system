import { describe, expect, it } from 'vitest'

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

  it('exposes additive maps view and modify keys', () => {
    expect(PermissionKeys.maps.view).toBe('maps:view')
    expect(PermissionKeys.maps.modify).toBe('maps:modify')

    const granted = new Set([PermissionKeys.maps.view])
    expect(canViewResource(granted, 'maps')).toBe(true)
    expect(canModifyResource(granted, 'maps')).toBe(false)
  })

  it('maps legacy read/crud keys during transition', () => {
    const granted = new Set(['dashboard:read', 'users:update'])

    expect(canViewResource(granted, 'dashboard')).toBe(true)
    expect(canModifyResource(granted, 'users')).toBe(true)
  })
})
