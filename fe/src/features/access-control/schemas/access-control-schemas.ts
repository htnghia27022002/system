import { z } from 'zod'

import { personalFieldsSchema } from '@/features/user-profile'

export const userStatusSchema = z.enum(['active', 'inactive'])
/** Edit form also round-trips `invited` (shown read-only; never sent to the BE). */
export const editableUserStatusSchema = z.enum(['active', 'inactive', 'invited'])

const personalUserFields = personalFieldsSchema.pick({
  phone: true,
  general: true,
  birthday: true,
  address: true,
  socialLinks: true,
})

export const createUserSchema = z
  .object({
    name: z.string().min(1, 'Name is required'),
    email: z.email({ message: 'Enter a valid email address' }),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    roleId: z.string().min(1, 'Role is required'),
    status: userStatusSchema,
    superAdmin: z.boolean().optional(),
  })
  .merge(personalUserFields)

export const updateUserSchema = z
  .object({
    name: z.string().min(1, 'Name is required'),
    email: z.email({ message: 'Enter a valid email address' }),
    password: z
      .union([
        z.literal(''),
        z.string().min(8, 'Password must be at least 8 characters'),
      ])
      .optional(),
    roleId: z.string().min(1, 'Role is required'),
    status: editableUserStatusSchema,
    superAdmin: z.boolean().optional(),
  })
  .merge(personalUserFields)

export const inviteUserSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.email({ message: 'Enter a valid email address' }),
  roleId: z.string().min(1, 'Role is required'),
})

export const createRoleSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  slug: z
    .string()
    .min(1, 'Slug is required')
    .regex(/^[a-z0-9-]+$/, 'Slug must be lowercase alphanumeric with hyphens'),
  permissionKeys: z.array(z.string()).min(1, 'Select at least one permission'),
})

export const updateRoleSchema = createRoleSchema

export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
export type CreateRoleFormValues = z.infer<typeof createRoleSchema>
export type UpdateRoleFormValues = z.infer<typeof updateRoleSchema>
export type InviteUserFormValues = z.infer<typeof inviteUserSchema>
