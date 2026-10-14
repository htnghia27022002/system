export type DashboardKpiKey =
  | 'totalUsers'
  | 'activeSessions'
  | 'newUsers'
  | 'roles'

export type DashboardKpi = {
  key: DashboardKpiKey
  value: number
  changePercent: number
}

export type DashboardChartPoint = {
  /** Calendar day, ISO `YYYY-MM-DD`. The UI formats it in the active locale. */
  date: string
  value: number
}

export type DashboardOverview = {
  kpis: DashboardKpi[]
  weeklyActivity: DashboardChartPoint[]
  updatedAt: string
}
