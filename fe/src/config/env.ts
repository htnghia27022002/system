import { z } from 'zod'

const envSchema = z.object({
  VITE_API_BASE_URL: z.url(),
  VITE_APP_NAME: z.string().min(1),
  VITE_USE_MOCK_API: z
    .enum(['true', 'false'])
    .optional()
    .transform((value) => value === 'true'),
  VITE_MOCK_API_DELAY_MS: z.coerce.number().int().positive().default(1200),
})

export const env = envSchema.parse({
  VITE_API_BASE_URL:
    process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:3000/api',
  VITE_APP_NAME: process.env.NEXT_PUBLIC_APP_NAME ?? 'System App',
  VITE_USE_MOCK_API: process.env.NEXT_PUBLIC_USE_MOCK_API ?? 'true',
  VITE_MOCK_API_DELAY_MS: process.env.NEXT_PUBLIC_MOCK_API_DELAY_MS ?? '1200',
})
