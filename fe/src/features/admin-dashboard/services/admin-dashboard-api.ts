import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { loadDashboardMock } from '@/services/mock'

import type { DashboardOverview } from '../types'

/**
 * The BE does not expose `GET /admin/dashboard/overview` yet, so the overview is
 * served from the mock regardless of `NEXT_PUBLIC_USE_MOCK_API`. Flip this to
 * `true` once the endpoint ships (response shape: `DashboardOverview`).
 */
const OVERVIEW_ENDPOINT_AVAILABLE = false

export const adminDashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    if (env.USE_MOCK_API || !OVERVIEW_ENDPOINT_AVAILABLE) {
      return (await loadDashboardMock()).getOverview()
    }

    const { data } = await apiClient.get<DashboardOverview>(
      '/admin/dashboard/overview',
    )
    return data
  },
}
