import { renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it } from 'vitest'

import { useAuthStore } from '@/store/auth-store'

import { usePermissions } from './use-permissions'

beforeEach(() => {
  useAuthStore.setState({
    accessToken: null,
    user: null,
    isAdmin: false,
  })
})

describe('usePermissions', () => {
  it('checks permission keys from the auth store', () => {
    useAuthStore.setState({
      accessToken: 'token',
      user: {
        id: '1',
        email: 'a@example.com',
        name: 'Admin',
        role: 'admin',
        roleId: 'role-admin',
        permissions: ['users:read', 'users:create'],
      },
      isAdmin: true,
    })

    const { result } = renderHook(() => usePermissions())

    expect(result.current.hasPermission('users:read')).toBe(true)
    expect(result.current.hasPermission('roles:delete')).toBe(false)
    expect(result.current.hasAny('users:create', 'roles:read')).toBe(true)
  })
})
