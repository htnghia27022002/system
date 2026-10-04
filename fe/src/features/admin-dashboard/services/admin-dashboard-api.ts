import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { loadDashboardMock } from '@/services/mock'

import type { DashboardOverview } from '../types'

export const adminDashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    if (env.USE_MOCK_API) {
      return (await loadDashboardMock()).getOverview()
    }

    const { data } = await apiClient.get<DashboardOverview>(
      '/admin/dashboard/overview',
    )
    return data
  },
}
