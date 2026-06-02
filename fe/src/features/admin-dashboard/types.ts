export type DashboardKpiKey =
  | 'totalUsers'
  | 'activeSessions'
  | 'revenue'
  | 'ordersToday'

export type DashboardKpi = {
  key: DashboardKpiKey
  value: number
  changePercent: number
}

export type DashboardChartPoint = {
  label: string
  value: number
}

export type DashboardOverview = {
  kpis: DashboardKpi[]
  weeklyActivity: DashboardChartPoint[]
  updatedAt: string
}
