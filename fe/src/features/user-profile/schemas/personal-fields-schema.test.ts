import { describe, expect, it } from 'vitest'

import {
  personalFieldsSchema,
  socialLinkSchema,
} from './personal-fields-schema'
import {
  changePasswordSchema,
  profileFormSchema,
  validateAvatarFile,
} from './profile-schemas'

describe('personal fields schema', () => {
  it('rejects future birthday', () => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    const iso = tomorrow.toISOString().slice(0, 10)
    const result = personalFieldsSchema.safeParse({
      phone: '',
      general: '',
      birthday: iso,
      address: '',
      socialLinks: [],
    })
    expect(result.success).toBe(false)
  })

  it('accepts today as birthday', () => {
    const today = new Date()
    const iso = [
      today.getFullYear(),
      String(today.getMonth() + 1).padStart(2, '0'),
      String(today.getDate()).padStart(2, '0'),
    ].join('-')
    const result = personalFieldsSchema.safeParse({
      phone: '+84',
      general: 'About me',
      birthday: iso,
      address: 'Somewhere',
      socialLinks: [],
    })
    expect(result.success).toBe(true)
  })

  it('rejects non-http social URL', () => {
    const result = socialLinkSchema.safeParse({
      label: 'Site',
      url: 'ftp://example.com',
    })
    expect(result.success).toBe(false)
  })

  it('rejects more than 5 social links', () => {
    const links = Array.from({ length: 6 }, (_, i) => ({
      label: `L${i}`,
      url: `https://example.com/${i}`,
    }))
    const result = personalFieldsSchema.safeParse({
      phone: '',
      general: '',
      birthday: '',
      address: '',
      socialLinks: links,
    })
    expect(result.success).toBe(false)
  })

  it('rejects over-long general text', () => {
    const result = personalFieldsSchema.safeParse({
      phone: '',
      general: 'x'.repeat(1001),
      birthday: '',
      address: '',
      socialLinks: [],
    })
    expect(result.success).toBe(false)
  })
})

describe('profile schemas', () => {
  it('requires name on profile form', () => {
    const result = profileFormSchema.safeParse({
      name: 'A',
      phone: '',
      general: '',
      birthday: '',
      address: '',
      socialLinks: [],
    })
    expect(result.success).toBe(false)
  })

  it('requires matching passwords', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'oldpassword',
      newPassword: 'newpassword',
      confirmPassword: 'different',
    })
    expect(result.success).toBe(false)
  })

  it('validates avatar file type and size', () => {
    const ok = new File(['x'], 'a.png', { type: 'image/png' })
    expect(validateAvatarFile(ok)).toBeNull()

    const badType = new File(['x'], 'a.gif', { type: 'image/gif' })
    expect(validateAvatarFile(badType)).toMatch(/JPEG|PNG|WebP/i)

    const big = new File([new Uint8Array(2 * 1024 * 1024 + 1)], 'a.png', {
      type: 'image/png',
    })
    expect(validateAvatarFile(big)).toMatch(/2 MB/i)
  })
})
