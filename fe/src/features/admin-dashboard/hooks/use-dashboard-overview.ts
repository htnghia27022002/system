'use client'

import { useQuery } from '@tanstack/react-query'

import { adminDashboardApi } from '../services/admin-dashboard-api'

export function useDashboardOverview() {
  return useQuery({
    queryKey: ['admin', 'dashboard', 'overview'],
    queryFn: () => adminDashboardApi.getOverview(),
  })
}
