import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { mockDashboardApi } from '@/services/mock/dashboard.mock'

import type { DashboardOverview } from '../types'

export const adminDashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    if (env.VITE_USE_MOCK_API) {
      return mockDashboardApi.getOverview()
    }

    const { data } = await apiClient.get<DashboardOverview>(
      '/admin/dashboard/overview',
    )
    return data
  },
}
