import { env } from '@/config/env'
import type { DashboardOverview } from '@/features/admin-dashboard/types'
import { mockDelay } from '@/utils/mock-delay'

const WEEKLY_VALUES = [42, 58, 51, 67, 73, 49, 61]

/** Last 7 calendar days ending today, oldest first, as `YYYY-MM-DD`. */
function lastSevenDays(today = new Date()): string[] {
  return WEEKLY_VALUES.map((_, index) => {
    const day = new Date(today)
    day.setDate(today.getDate() - (WEEKLY_VALUES.length - 1 - index))
    const month = String(day.getMonth() + 1).padStart(2, '0')
    const date = String(day.getDate()).padStart(2, '0')
    return `${day.getFullYear()}-${month}-${date}`
  })
}

export const mockDashboardApi = {
  async getOverview(): Promise<DashboardOverview> {
    await mockDelay(env.MOCK_API_DELAY_MS)

    return {
      updatedAt: new Date().toISOString(),
      kpis: [
        { key: 'totalUsers', value: 1284, changePercent: 8.2 },
        { key: 'activeSessions', value: 342, changePercent: -3.1 },
        { key: 'revenue', value: 48250, changePercent: 12.4 },
        { key: 'ordersToday', value: 96, changePercent: 5.6 },
      ],
      weeklyActivity: lastSevenDays().map((date, index) => ({
        date,
        value: WEEKLY_VALUES[index],
      })),
    }
  },
}
