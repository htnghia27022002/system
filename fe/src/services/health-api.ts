import { env } from '@/config/env'
import { loadHealthMock } from '@/services/mock'

import { apiClient } from './api-client'

export type HealthResponse = {
  status: string
}

export async function fetchHealth(): Promise<HealthResponse> {
  if (env.USE_MOCK_API) {
    return (await loadHealthMock()).getHealth()
  }

  try {
    const { data } = await apiClient.get<HealthResponse>('/health')
    return data
  } catch {
    return { status: 'offline' }
  }
}
