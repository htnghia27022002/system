export type AuditLogEntry = {
  id: string
  actorId: string | null
  actorEmail: string
  /** `<resource>.<verb>`, e.g. `user.update`, `auth.login_failed`. */
  action: string
  resourceType: string
  resourceId: string
  metadata: Record<string, unknown>
  ipAddress: string
  userAgent: string
  requestId: string
  createdAt: string
}

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
