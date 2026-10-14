import { fireEvent, render, screen, within } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import i18n from '@/config/i18n'

import type { DashboardOverview } from '../types'

const { mockState } = vi.hoisted(() => ({
  mockState: {
    data: undefined as DashboardOverview | undefined,
    isLoading: false,
    isError: false,
    isFetching: false,
    refetch: vi.fn(),
  },
}))

vi.mock('../hooks/use-dashboard-overview', () => ({
  useDashboardOverview: () => mockState,
}))

// recharts' ResponsiveContainer needs ResizeObserver, which jsdom lacks.
vi.stubGlobal(
  'ResizeObserver',
  class {
    observe() {}
    unobserve() {}
    disconnect() {}
  },
)

import { AdminDashboardOverview } from './admin-dashboard-overview'

const overview: DashboardOverview = {
  updatedAt: '2026-10-14T09:30:00Z',
  kpis: [
    { key: 'totalUsers', value: 1284, changePercent: 8.2 },
    { key: 'activeSessions', value: 342, changePercent: -3.1 },
    { key: 'newUsers', value: 57, changePercent: 12.4 },
    { key: 'roles', value: 6, changePercent: 0 },
  ],
  weeklyActivity: [
    { date: '2026-10-08', value: 42 },
    { date: '2026-10-09', value: 58 },
    { date: '2026-10-10', value: 51 },
    { date: '2026-10-11', value: 67 },
    { date: '2026-10-12', value: 73 },
    { date: '2026-10-13', value: 49 },
    { date: '2026-10-14', value: 61 },
  ],
}

describe('AdminDashboardOverview', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
    mockState.data = overview
    mockState.isLoading = false
    mockState.isError = false
    mockState.isFetching = false
    mockState.refetch.mockReset()
  })

  it('renders KPI tiles with locale-formatted values and signed change', () => {
    render(<AdminDashboardOverview />)

    expect(
      screen.getByRole('heading', { level: 1, name: 'Dashboard overview' }),
    ).toBeInTheDocument()
    expect(screen.getByText('1,284')).toBeInTheDocument()
    expect(screen.getByText('57')).toBeInTheDocument()
    expect(screen.getByText('+8.2%')).toBeInTheDocument()
    expect(screen.getByText('-3.1%')).toBeInTheDocument()
    // No change is neutral, not a green "up" trend.
    expect(screen.getByText('0.0%')).not.toHaveClass('text-success')
  })

  it('exposes the weekly chart as an accessible table, one row per day', async () => {
    render(<AdminDashboardOverview />)

    const table = await screen.findByRole('table', { name: 'Weekly activity' })
    const rows = within(table).getAllByRole('row')
    // header row + 7 days
    expect(rows).toHaveLength(8)
    expect(within(rows[1]).getByRole('rowheader')).toHaveTextContent('Thursday')
    expect(within(rows[7]).getByRole('cell')).toHaveTextContent('61')
  })

  it('shows an error with retry when the overview fails to load', () => {
    mockState.data = undefined
    mockState.isError = true
    render(<AdminDashboardOverview />)

    expect(screen.getByRole('alert')).toHaveTextContent(
      'Could not load dashboard data.',
    )
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }))
    expect(mockState.refetch).toHaveBeenCalledTimes(1)
  })
})
