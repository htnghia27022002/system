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
    import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:3000/api',
  VITE_APP_NAME: import.meta.env.VITE_APP_NAME ?? 'System App',
  VITE_USE_MOCK_API: import.meta.env.VITE_USE_MOCK_API ?? 'true',
  VITE_MOCK_API_DELAY_MS: import.meta.env.VITE_MOCK_API_DELAY_MS ?? '1200',
})
