import { describe, expect, it } from 'vitest'

import { createRoleSchema, createUserSchema } from './access-control-schemas'

describe('access-control schemas', () => {
  it('requires valid create user fields', () => {
    const result = createUserSchema.safeParse({
      name: '',
      email: 'not-an-email',
      password: 'short',
      roleId: '',
    })

    expect(result.success).toBe(false)
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
