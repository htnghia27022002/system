import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  ensureRbacSeed,
  mockAccessControlApi,
  MockAccessControlError,
  resolveAuthUserByCredentials,
} from './access-control.mock'

const storage = new Map<string, string>()

beforeEach(() => {
  storage.clear()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => storage.get(key) ?? null,
    setItem: (key: string, value: string) => {
      storage.set(key, value)
    },
    removeItem: (key: string) => {
      storage.delete(key)
    },
    clear: () => storage.clear(),
    key: () => null,
    length: 0,
  })
  ensureRbacSeed()
})

describe('mockAccessControlApi', () => {
  it('seeds admin and demo users with role permissions', async () => {
    const admin = resolveAuthUserByCredentials(
      'admin@example.com',
      'admin1234',
    )
    expect(admin?.role).toBe('admin')
    expect(admin?.permissions).toContain('users:read')

    const demo = resolveAuthUserByCredentials(
      'demo@example.com',
      'password123',
    )
    expect(demo?.role).toBe('user')
    expect(demo?.permissions).toEqual(['dashboard:read'])
  })

  it('creates and lists users', async () => {
    const created = await mockAccessControlApi.createUser({
      name: 'Test User',
      email: 'test@example.com',
      password: 'password123',
      roleId: 'role-user',
    })

    const page = await mockAccessControlApi.listUsers({ search: 'test@' })
    expect(page.items.some((u) => u.id === created.id)).toBe(true)
  })

  it('blocks deleting the last administrator', async () => {
    await expect(
      mockAccessControlApi.deleteUser('admin-user', 'admin-user'),
    ).rejects.toMatchObject({ status: 400 })
  })

  it('blocks deleting a role assigned to users', async () => {
    await expect(mockAccessControlApi.deleteRole('role-admin')).rejects.toThrow(
      MockAccessControlError,
    )
  })
})
