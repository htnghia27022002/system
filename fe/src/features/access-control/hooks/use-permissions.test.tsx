import { renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { PermissionKeys } from '../permission-keys'
import { usePermissions } from './use-permissions'

vi.mock('@/store/auth-store', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) =>
    selector({
      user: {
        role: 'admin',
        roleId: 'role-admin',
        permissions: [
          PermissionKeys.users.view,
          PermissionKeys.users.modify,
        ],
      },
      sessionSynced: true,
    }),
  selectPermissions: (state: { user?: { permissions?: string[] } }) =>
    state.user?.permissions ?? [],
}))

describe('usePermissions', () => {
  it('checks view and modify permissions', () => {
    const { result } = renderHook(() => usePermissions())

    expect(result.current.hasPermission(PermissionKeys.users.view)).toBe(true)
    expect(result.current.hasPermission(PermissionKeys.roles.modify)).toBe(false)
    expect(result.current.canModify('users')).toBe(true)
    expect(result.current.canView('roles')).toBe(false)
  })
})
