import { env } from '@/config/env'
import { mockDelay } from '@/utils/mock-delay'

export type HealthResponse = {
  status: 'ok' | 'offline'
}

export const mockHealthApi = {
  async getHealth(): Promise<HealthResponse> {
    await mockDelay(env.MOCK_API_DELAY_MS)
    return { status: 'ok' }
  },
}
