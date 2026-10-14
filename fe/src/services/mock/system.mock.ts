import { env } from '@/config/env'
import type {
  AppSetting,
  AuditLogEntry,
  ListAuditLogsParams,
  SettingValues,
} from '@/features/system/types'
import type { PaginatedResponse } from '@/types/pagination'
import { mockDelay } from '@/utils/mock-delay'

import { MockAccessControlError } from './mock-errors'

// Mirrors be/internal/settings CoreDefinitions().
const SETTINGS: AppSetting[] = [
  {
    key: 'app.supportEmail', type: 'string', default: '', group: 'app', public: true, maxLength: 255,
    label: 'Support email', description: 'Shown to users who need help. Leave empty to hide.',
    value: '', overridden: false,
  },
  {
    key: 'audit.retentionDays', type: 'int', default: 180, group: 'audit', public: false, min: 7, max: 3650,
    label: 'Audit log retention (days)', description: 'Entries older than this are deleted by the daily cleanup.',
    value: 180, overridden: false,
  },
  {
    key: 'auth.registrationEnabled', type: 'bool', default: true, group: 'auth', public: true,
    label: 'Allow self sign-up',
    description: 'When off, only invited users can join (password and OAuth sign-up are refused).',
    value: true, overridden: false,
  },
]

const ACTIONS = ['auth.login', 'auth.login_failed', 'user.update', 'role.update', 'settings.update', 'user.invite']

function seedAuditLogs(): AuditLogEntry[] {
  const now = Date.now()
  return Array.from({ length: 36 }, (_, i) => {
    const action = ACTIONS[i % ACTIONS.length]
    const admin = action !== 'auth.login_failed'
    return {
      id: `audit-${i}`,
      actorId: admin ? 'admin-user' : 'demo-user',
      actorEmail: admin ? 'admin@example.com' : 'demo@example.com',
      action,
      resourceType: action.startsWith('settings') ? 'settings' : action.startsWith('role') ? 'role' : 'user',
      resourceId: action.startsWith('settings') ? '' : action.startsWith('role') ? 'role-user' : 'demo-user',
      metadata: action === 'user.update' ? { fields: ['name', 'status'] } : action === 'auth.login' ? { method: 'password' } : {},
      ipAddress: '203.0.113.' + (10 + (i % 5)),
      userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 Safari/605.1.15',
      requestId: `req-${1000 + i}`,
      createdAt: new Date(now - i * 47 * 60 * 1000).toISOString(),
    }
  })
}

let auditLogs: AuditLogEntry[] | null = null
const settings = SETTINGS.map((s) => ({ ...s }))

export const mockSystemApi = {
  async listAuditLogs(params: ListAuditLogsParams): Promise<PaginatedResponse<AuditLogEntry>> {
    await mockDelay(env.MOCK_API_DELAY_MS)
    auditLogs ??= seedAuditLogs()
    const term = params.search?.toLowerCase().trim()
    const items = auditLogs.filter(
      (e) =>
        (!params.action || e.action === params.action) &&
        (!params.resourceType || e.resourceType === params.resourceType) &&
        (!term || [e.action, e.actorEmail, e.resourceId].some((v) => v.toLowerCase().includes(term))),
    )
    const page = params.page ?? 1
    const pageSize = params.pageSize ?? 10
    return {
      items: items.slice((page - 1) * pageSize, page * pageSize),
      total: items.length,
      page,
      pageSize,
    }
  },

  async listSettings(): Promise<AppSetting[]> {
    await mockDelay(env.MOCK_API_DELAY_MS)
    return settings.map((s) => ({ ...s }))
  },

  async updateSettings(values: SettingValues): Promise<AppSetting[]> {
    await mockDelay(env.MOCK_API_DELAY_MS)
    for (const [key, value] of Object.entries(values)) {
      const setting = settings.find((s) => s.key === key)
      if (!setting) throw new MockAccessControlError(`unknown setting "${key}"`, 400)
      if (setting.type === 'int' && (typeof value !== 'number' || (setting.max !== undefined && (value < (setting.min ?? 0) || value > setting.max)))) {
        throw new MockAccessControlError(`${key}: must be between ${setting.min} and ${setting.max}`, 400)
      }
      setting.value = value
      setting.overridden = true
    }
    return settings.map((s) => ({ ...s }))
  },

  /** Public subset, as returned by GET /api/app/info. */
  publicSettings(): Record<string, string | boolean | number> {
    return Object.fromEntries(settings.filter((s) => s.public).map((s) => [s.key, s.value]))
  },
}
