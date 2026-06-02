import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'

import { useDashboardOverview } from '../hooks/use-dashboard-overview'
import { DashboardCharts } from './dashboard-charts'
import { DashboardPageHeader } from './dashboard-page-header'
import { DashboardStatGrid } from './dashboard-stat-grid'

export function AdminDashboardOverview() {
  const { t } = useTranslation('admin')
  const { data, isLoading, isError, refetch, isFetching } = useDashboardOverview()

  return (
    <div className="mx-auto max-w-7xl">
      <DashboardPageHeader
        title={t('overview.title')}
        description={t('overview.description')}
      />

      {isError ? (
        <div className="mb-6 rounded-[var(--ds-radius-lg)] border border-border bg-card p-6">
          <p className="text-[17px] text-destructive">{t('overview.error')}</p>
          <Button
            variant="outline"
            size="sm"
            className="mt-4"
            onClick={() => void refetch()}
          >
            {t('overview.retry')}
          </Button>
        </div>
      ) : null}

      <DashboardStatGrid kpis={data?.kpis} isLoading={isLoading} />
      <DashboardCharts
        weeklyActivity={data?.weeklyActivity}
        isLoading={isLoading}
      />

      {data?.updatedAt && !isLoading ? (
        <p className="mt-4 text-[12px] leading-none tracking-[-0.12px] text-muted-foreground">
          {t('overview.lastUpdated', {
            time: new Date(data.updatedAt).toLocaleString(),
          })}
          {isFetching ? ` · ${t('overview.refreshing')}` : ''}
        </p>
      ) : null}
    </div>
  )
}
