'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { translateApiError } from '@/lib/api-error'

import { systemApi } from '../services/system-api'
import type { ListAuditLogsParams, SettingValues } from '../types'

const systemKey = ['admin', 'system'] as const
const settingsKey = [...systemKey, 'settings'] as const

export function useAuditLogs(params: ListAuditLogsParams) {
  return useQuery({
    queryKey: [...systemKey, 'audit-logs', params],
    queryFn: () => systemApi.listAuditLogs(params),
    meta: { skipNavLoading: true },
  })
}

export function useSettings() {
  return useQuery({ queryKey: settingsKey, queryFn: () => systemApi.listSettings() })
}

export function useUpdateSettings() {
  const { t } = useTranslation('admin')
  const { t: tc } = useTranslation('common')
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (values: SettingValues) => systemApi.updateSettings(values),
    onSuccess: (items) => {
      queryClient.setQueryData(settingsKey, items)
      // Public settings (e.g. sign-up) are served through /app/info.
      void queryClient.invalidateQueries({ queryKey: ['app', 'info'] })
      toast.success(t('system.settings.saved'))
    },
    onError: (error) => {
      toast.error(translateApiError(error, tc, t('system.settings.saveFailed')))
    },
  })
}
