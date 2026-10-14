import { z } from 'zod'

export const loginSchema = z.object({
  email: z.email({ message: 'Enter a valid email address' }),
  password: z.string().min(8, 'Password must be at least 8 characters'),
})

export const registerSchema = z
  .object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.email({ message: 'Enter a valid email address' }),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string().min(8, 'Confirm your password'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  })

export const forgotPasswordSchema = z.object({
  email: z.email({ message: 'Enter a valid email address' }),
})

const newPasswordFields = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string().min(8, 'Confirm your password'),
})

const passwordsMatch = {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
}

export const resetPasswordSchema = newPasswordFields.refine(
  (data) => data.password === data.confirmPassword,
  passwordsMatch,
)

export const acceptInviteSchema = newPasswordFields
  .extend({ name: z.string().min(2, 'Name must be at least 2 characters') })
  .refine((data) => data.password === data.confirmPassword, passwordsMatch)

export type LoginFormValues = z.infer<typeof loginSchema>
export type RegisterFormValues = z.infer<typeof registerSchema>
export type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>
export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>
export type AcceptInviteFormValues = z.infer<typeof acceptInviteSchema>
