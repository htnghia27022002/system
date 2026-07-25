import { z } from 'zod'

import { personalFieldsSchema } from './personal-fields-schema'

export const profileFormSchema = personalFieldsSchema.extend({
  name: z.string().min(2, 'Name must be at least 2 characters'),
})

export const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(1, 'Current password is required'),
    newPassword: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string().min(8, 'Confirm your password'),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

export type ProfileFormValues = z.infer<typeof profileFormSchema>
export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>

export const AVATAR_MAX_BYTES = 2 * 1024 * 1024
export const AVATAR_ALLOWED_TYPES = [
  'image/jpeg',
  'image/png',
  'image/webp',
] as const

export function validateAvatarFile(file: File): string | null {
  if (
    !AVATAR_ALLOWED_TYPES.includes(
      file.type as (typeof AVATAR_ALLOWED_TYPES)[number],
    )
  ) {
    return 'Avatar must be a JPEG, PNG, or WebP image'
  }
  if (file.size > AVATAR_MAX_BYTES) {
    return 'Avatar must be 2 MB or smaller'
  }
  return null
}
