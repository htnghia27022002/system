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

  it('maps legacy read/crud keys during transition', () => {
    const granted = new Set(['dashboard:read', 'users:update'])

    expect(canViewResource(granted, 'dashboard')).toBe(true)
    expect(canModifyResource(granted, 'users')).toBe(true)
  })
})
