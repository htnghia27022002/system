import { z } from 'zod'

export const userStatusSchema = z.enum(['active', 'inactive'])

export const createUserSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.email({ message: 'Enter a valid email address' }),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  roleId: z.string().min(1, 'Role is required'),
  status: userStatusSchema,
})

export const updateUserSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.email({ message: 'Enter a valid email address' }),
  password: z
    .union([
      z.literal(''),
      z.string().min(8, 'Password must be at least 8 characters'),
    ])
    .optional(),
  roleId: z.string().min(1, 'Role is required'),
  status: userStatusSchema,
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
