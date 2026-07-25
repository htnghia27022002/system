import { describe, expect, it } from 'vitest'

import { createRoleSchema, createUserSchema } from './access-control-schemas'

describe('access-control schemas', () => {
  it('requires valid create user fields', () => {
    const result = createUserSchema.safeParse({
      name: '',
      email: 'not-an-email',
      password: 'short',
      roleId: '',
      status: 'active',
      phone: '',
      general: '',
      birthday: '',
      address: '',
      socialLinks: [],
    })

    expect(result.success).toBe(false)
  })

  it('accepts personal fields on create', () => {
    const result = createUserSchema.safeParse({
      name: 'Ada Lovelace',
      email: 'ada@example.com',
      password: 'password123',
      roleId: 'role-user',
      status: 'active',
      phone: '123',
      general: 'Bio',
      birthday: '1990-01-01',
      address: 'London',
      socialLinks: [{ label: 'Site', url: 'https://example.com' }],
    })
    expect(result.success).toBe(true)
  })

  it('accepts a valid role payload', () => {
    const result = createRoleSchema.safeParse({
      name: 'Editor',
      slug: 'editor',
      permissionKeys: ['dashboard:view'],
    })

    expect(result.success).toBe(true)
  })
})
