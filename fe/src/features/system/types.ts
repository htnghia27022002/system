import type { AuditLogResponse } from '@/types/api.gen'

/** Generated from the BE OpenAPI schema (pnpm gen:api). */
export type AuditLogEntry = AuditLogResponse

export type ListAuditLogsParams = {
  page?: number
  pageSize?: number
  search?: string
  action?: string
  resourceType?: string
}

export type SettingType = 'string' | 'bool' | 'int'

/** A setting declared on the BE (core or module) with its effective value. */
export type AppSetting = {
  key: string
  type: SettingType
  default: string | boolean | number
  group: string
  label: string
  description?: string
  public: boolean
  min?: number
  max?: number
  maxLength?: number
  value: string | boolean | number
  overridden: boolean
}

export type SettingValues = Record<string, string | boolean | number>
