import { env } from '@/config/env'
import type { DashboardOverview } from '@/features/admin-dashboard/types'
import { mockDelay } from '@/utils/mock-delay'

export const mockDashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    await mockDelay(env.VITE_MOCK_API_DELAY_MS)

    return {
      updatedAt: new Date().toISOString(),
      kpis: [
        { key: 'totalUsers', value: 1284, changePercent: 8.2 },
        { key: 'activeSessions', value: 342, changePercent: -3.1 },
        { key: 'revenue', value: 48250, changePercent: 12.4 },
        { key: 'ordersToday', value: 96, changePercent: 5.6 },
      ],
      weeklyActivity: [
        { label: 'Mon', value: 42 },
        { label: 'Tue', value: 58 },
        { label: 'Wed', value: 51 },
        { label: 'Thu', value: 67 },
        { label: 'Fri', value: 73 },
        { label: 'Sat', value: 49 },
        { label: 'Sun', value: 61 },
      ],
    }
  },
}
