import { afterEach, describe, expect, it, vi } from 'vitest'

import { mockDashboardApi } from './dashboard.mock'

describe('mockDashboardApi.getOverview', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it('returns seven consecutive days ending today, oldest first', async () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date(2026, 0, 3, 10, 0)) // 3 Jan — window crosses a year boundary

    const pending = mockDashboardApi.getOverview()
    await vi.runAllTimersAsync()
    const overview = await pending

    expect(overview.weeklyActivity.map((point) => point.date)).toEqual([
      '2025-12-28',
      '2025-12-29',
      '2025-12-30',
      '2025-12-31',
      '2026-01-01',
      '2026-01-02',
      '2026-01-03',
    ])
    expect(overview.kpis.map((kpi) => kpi.key)).toEqual([
      'totalUsers',
      'activeSessions',
      'newUsers',
      'roles',
    ])
  })
})
