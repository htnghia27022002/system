import { env } from '@/config/env'
import { normalizePaginatedResponse } from '@/lib/normalize-paginated-response'
import { apiClient } from '@/services/api-client'
import { loadSystemMock } from '@/services/mock'
import type { PaginatedResponse } from '@/types/pagination'

import type { AppSetting, AuditLogEntry, ListAuditLogsParams, SettingValues } from '../types'

export const systemApi = {
  async listAuditLogs(params: ListAuditLogsParams): Promise<PaginatedResponse<AuditLogEntry>> {
    if (env.USE_MOCK_API) return (await loadSystemMock()).listAuditLogs(params)
    const { data } = await apiClient.get<PaginatedResponse<AuditLogEntry>>('/admin/audit-logs', {
      params,
      skipNavLoading: true,
    })
    return normalizePaginatedResponse<AuditLogEntry>(data, params.pageSize)
  },

  async listSettings(): Promise<AppSetting[]> {
    if (env.USE_MOCK_API) return (await loadSystemMock()).listSettings()
    const { data } = await apiClient.get<{ items: AppSetting[] }>('/admin/settings')
    return data.items
  },

  /** Partial update; returns every setting with its new effective value. */
  async updateSettings(values: SettingValues): Promise<AppSetting[]> {
    if (env.USE_MOCK_API) return (await loadSystemMock()).updateSettings(values)
    const { data } = await apiClient.patch<{ items: AppSetting[] }>('/admin/settings', { values })
    return data.items
  },
}
