import { render, screen, within } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import '@/config/i18n'

import type { AccountSession } from '../types'

const currentId = 'sess-current'
const otherId = 'sess-other'

const { mockState } = vi.hoisted(() => ({
  mockState: {
    query: {
      isLoading: false,
      isError: false,
      isFetching: false,
      refetch: vi.fn(),
    },
    items: [] as AccountSession[],
    hasOtherSessions: true,
    revokeSession: { mutate: vi.fn(), isPending: false },
    revokeOthers: { mutate: vi.fn(), isPending: false },
  },
}))

vi.mock('../hooks/use-account-sessions', () => ({
  useAccountSessions: () => mockState,
}))

import { AccountSessionsCard } from './account-sessions-card'

function session(partial: Partial<AccountSession> & Pick<AccountSession, 'id' | 'current'>): AccountSession {
  return {
    createdAt: '2026-09-04T10:00:00Z',
    expiresAt: '2026-09-18T10:00:00Z',
    lastUsedAt: '2026-09-04T12:30:00Z',
    ipAddress: '203.0.113.10',
    userAgent: 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36',
    ...partial,
  }
}

describe('AccountSessionsCard', () => {
  beforeEach(() => {
    mockState.query.isLoading = false
    mockState.query.isError = false
    mockState.hasOtherSessions = true
    mockState.items = [
      session({ id: currentId, current: true }),
      session({
        id: otherId,
        current: false,
        ipAddress: null,
        userAgent: null,
      }),
    ]
  })

  it('highlights the current row and does not offer Revoke on it', () => {
    render(<AccountSessionsCard />)

    const currentRows = screen.getAllByTestId(`session-row-${currentId}`)
    expect(currentRows.length).toBeGreaterThan(0)
    for (const currentRow of currentRows) {
      expect(currentRow).toHaveAttribute('data-current', 'true')
      expect(within(currentRow).getByText('This device')).toBeInTheDocument()
      expect(
        within(currentRow).queryByRole('button', { name: 'Revoke' }),
      ).not.toBeInTheDocument()
    }

    const otherRows = screen.getAllByTestId(`session-row-${otherId}`)
    expect(otherRows.length).toBeGreaterThan(0)
    for (const otherRow of otherRows) {
      expect(otherRow).toHaveAttribute('data-current', 'false')
    }
    expect(
      screen.getAllByRole('button', { name: 'Revoke' }).length,
    ).toBeGreaterThan(0)
    expect(screen.getAllByText('Unknown').length).toBeGreaterThan(0)

    expect(
      screen.getByRole('button', { name: 'Revoke all other sessions' }),
    ).toBeInTheDocument()
  })

  it('hides revoke-all-others when there are no other sessions', () => {
    mockState.hasOtherSessions = false
    mockState.items = [session({ id: currentId, current: true })]

    render(<AccountSessionsCard />)

    expect(
      screen.queryByRole('button', { name: 'Revoke all other sessions' }),
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Revoke' }),
    ).not.toBeInTheDocument()
  })
})
