import { z } from 'zod'

function todayIsoDate(): string {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

export const socialLinkSchema = z.object({
  label: z
    .string()
    .max(50, 'Label must be at most 50 characters')
    .optional()
    .or(z.literal('')),
  url: z
    .string()
    .min(1, 'URL is required')
    .url({ message: 'Enter a valid URL' })
    .refine(
      (value) => /^https?:\/\//i.test(value),
      'URL must start with http:// or https://',
    ),
})

export const personalFieldsSchema = z.object({
  phone: z
    .string()
    .max(50, 'Phone must be at most 50 characters')
    .optional()
    .or(z.literal('')),
  general: z
    .string()
    .max(1000, 'General must be at most 1000 characters')
    .optional()
    .or(z.literal('')),
  birthday: z
    .string()
    .optional()
    .or(z.literal(''))
    .refine(
      (value) => !value || value <= todayIsoDate(),
      'Birthday cannot be in the future',
    ),
  address: z
    .string()
    .max(500, 'Address must be at most 500 characters')
    .optional()
    .or(z.literal('')),
  // Required array (default [] in form defaultValues). Avoid .default([]) —
  // that makes input optional and breaks zodResolver + react-hook-form types.
  socialLinks: z
    .array(socialLinkSchema)
    .max(5, 'You can add at most 5 social links'),
})

export type PersonalFieldsFormValues = z.infer<typeof personalFieldsSchema>
export type SocialLinkFormValues = z.infer<typeof socialLinkSchema>

export const emptyPersonalFields = (): PersonalFieldsFormValues => ({
  phone: '',
  general: '',
  birthday: '',
  address: '',
  socialLinks: [],
})
