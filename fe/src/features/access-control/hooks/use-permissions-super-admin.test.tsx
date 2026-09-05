import { renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { PermissionKeys } from '../permission-keys'
import { usePermissions } from './use-permissions'

vi.mock('@/store/auth-store', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) =>
    selector({
      user: {
        role: 'user',
        roleId: 'role-user',
        permissions: [],
        superAdmin: true,
      },
      sessionSynced: true,
    }),
  selectPermissions: (state: { user?: { permissions?: string[] } }) =>
    state.user?.permissions ?? [],
  selectIsSuperAdmin: (state: { user?: { superAdmin?: boolean } }) =>
    Boolean(state.user?.superAdmin),
}))

describe('usePermissions super admin', () => {
  it('allows every permission without role keys', () => {
    const { result } = renderHook(() => usePermissions())

    expect(result.current.isSuperAdmin).toBe(true)
    expect(result.current.hasPermission(PermissionKeys.users.view)).toBe(true)
    expect(result.current.hasPermission(PermissionKeys.roles.modify)).toBe(true)
    expect(result.current.canView('webhooks')).toBe(true)
    expect(result.current.canModify('permissions')).toBe(true)
    expect(result.current.hasAny(PermissionKeys.dashboard.view)).toBe(true)
  })
})
