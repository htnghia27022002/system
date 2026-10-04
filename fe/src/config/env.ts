import { z } from 'zod'

const envSchema = z.object({
  API_BASE_URL: z.url(),
  SITE_URL: z.url(),
  APP_NAME: z.string().min(1),
  USE_MOCK_API: z
    .enum(['true', 'false'])
    .optional()
    .transform((value) => value === 'true'),
  MOCK_API_DELAY_MS: z.coerce.number().int().positive().default(1200),
  MAP_PROVIDER: z.enum(['osm', 'google']).default('osm'),
  GOOGLE_MAPS_API_KEY: z.string().optional(),
})

export const env = envSchema.parse({
  API_BASE_URL:
    process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:3000/api',
  SITE_URL:
    process.env.NEXT_PUBLIC_SITE_URL ?? 'http://localhost:3000',
  APP_NAME: process.env.NEXT_PUBLIC_APP_NAME ?? 'System App',
  // Mock API is opt-in: a missing variable must never serve fake auth in a real deploy.
  USE_MOCK_API: process.env.NEXT_PUBLIC_USE_MOCK_API ?? 'false',
  MOCK_API_DELAY_MS: process.env.NEXT_PUBLIC_MOCK_API_DELAY_MS ?? '1200',
  MAP_PROVIDER: process.env.NEXT_PUBLIC_MAP_PROVIDER || 'osm',
  GOOGLE_MAPS_API_KEY: process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || undefined,
})
