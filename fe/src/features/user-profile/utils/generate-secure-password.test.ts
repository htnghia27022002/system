import { describe, expect, it } from 'vitest'

import { generateSecurePassword } from './generate-secure-password'

describe('generateSecurePassword', () => {
  it('returns a mixed 16-character password', () => {
    const password = generateSecurePassword()
    expect(password).toHaveLength(16)
    expect(password).toMatch(/[A-Z]/)
    expect(password).toMatch(/[a-z]/)
    expect(password).toMatch(/[0-9]/)
    expect(password).toMatch(/[!@#$%^&*\-_=+]/)
  })

  it('does not return the same value twice in a row', () => {
    const first = generateSecurePassword()
    const second = generateSecurePassword()
    expect(first).not.toEqual(second)
  })
})
