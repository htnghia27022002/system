import { env } from '@/config/env'
import { mockHealthApi } from '@/services/mock/health.mock'

import { apiClient } from './api-client'

export type HealthResponse = {
  status: string
}

export async function fetchHealth(): Promise<HealthResponse> {
  if (env.VITE_USE_MOCK_API) {
    return mockHealthApi.getHealth()
  }

  try {
    const { data } = await apiClient.get<HealthResponse>('/health')
    return data
  } catch {
    return { status: 'offline' }
  }
}
