import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import i18n from '@/config/i18n'

import type { AppSetting } from '../types'
import { actionVariant, summarizeMetadata } from './audit-log-columns'

const { api, perms } = vi.hoisted(() => ({
  api: { listSettings: vi.fn(), updateSettings: vi.fn(), listAuditLogs: vi.fn() },
  perms: { canModify: true },
}))

vi.mock('../services/system-api', () => ({ systemApi: api }))
vi.mock('@/features/access-control', () => ({
  usePermissions: () => ({ canModify: () => perms.canModify }),
}))

import { SettingsPage } from './settings-page'

const SETTINGS: AppSetting[] = [
  { key: 'auth.registrationEnabled', type: 'bool', default: true, group: 'auth', label: 'Allow self sign-up', public: true, value: true, overridden: false },
  { key: 'audit.retentionDays', type: 'int', default: 180, group: 'audit', label: 'Retention', public: false, min: 7, max: 3650, value: 180, overridden: false },
  { key: 'billing.currency', type: 'string', default: 'USD', group: 'billing', label: 'Currency (module)', public: true, value: 'USD', overridden: false },
]

function wrap(ui: ReactNode) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>)
}

describe('audit log helpers', () => {
  it('flags failed and destructive actions', () => {
    expect(actionVariant('auth.login_failed')).toBe('destructive')
    expect(actionVariant('user.delete')).toBe('destructive')
    expect(actionVariant('user.update')).toBe('secondary')
  })

  it('summarizes metadata without nesting', () => {
    expect(summarizeMetadata({ fields: ['name', 'status'], method: 'password' })).toBe('fields: name, status · method: password')
    expect(summarizeMetadata({})).toBe('')
  })
})

describe('SettingsPage', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
    api.listSettings.mockReset().mockResolvedValue(SETTINGS)
    api.updateSettings.mockReset()
    perms.canModify = true
  })

  it('renders every declared setting by type and saves only what changed', async () => {
    api.updateSettings.mockResolvedValue(SETTINGS)
    wrap(<SettingsPage />)

    const toggle = await screen.findByRole('checkbox', { name: 'Allow self sign-up' })
    expect(screen.getByLabelText('Audit log retention (days)')).toHaveValue(180)
    expect(screen.getByLabelText('Currency (module)')).toHaveValue('USD') // module setting: BE label
    expect(screen.getByRole('button', { name: 'Save' })).toBeDisabled()

    fireEvent.click(toggle)
    fireEvent.change(screen.getByLabelText('Audit log retention (days)'), { target: { value: '30' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    await waitFor(() =>
      expect(api.updateSettings).toHaveBeenCalledWith({
        'auth.registrationEnabled': false,
        'audit.retentionDays': 30,
      }),
    )
  })

  it('is read-only without settings:modify', async () => {
    perms.canModify = false
    wrap(<SettingsPage />)
    expect(await screen.findByRole('checkbox', { name: 'Allow self sign-up' })).toBeDisabled()
    expect(screen.queryByRole('button', { name: 'Save' })).not.toBeInTheDocument()
  })
})
