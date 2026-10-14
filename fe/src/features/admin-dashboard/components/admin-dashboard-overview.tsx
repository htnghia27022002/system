'use client'

import { RefreshCwIcon } from 'lucide-react'
import dynamic from 'next/dynamic'
import { useTranslation } from 'react-i18next'

import { PageContainer, PageHeader } from '@/components/common/page-layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

import { useDashboardOverview } from '../hooks/use-dashboard-overview'
import { DashboardStatGrid } from './dashboard-stat-grid'

// recharts is ~100 kB; keep it out of the initial /admin bundle.
const DashboardCharts = dynamic(
  () => import('./dashboard-charts').then((mod) => mod.DashboardCharts),
  { ssr: false, loading: () => <Skeleton className="h-96 w-full rounded-xl" /> },
)

export function AdminDashboardOverview() {
  const { t, i18n } = useTranslation('admin')
  const query = useDashboardOverview()
  const overview = query.data
  const isRefreshing = query.isFetching && !query.isLoading

  const updatedAt = overview
    ? new Intl.DateTimeFormat(i18n.language, {
        timeStyle: 'short',
      }).format(new Date(overview.updatedAt))
    : null

  const actions = (
    <>
      {updatedAt ? (
        <p className="text-sm text-muted-foreground" aria-live="polite">
          {isRefreshing
            ? t('overview.refreshing')
            : t('overview.lastUpdated', { time: updatedAt })}
        </p>
      ) : null}
      <Button
        variant="outline"
        className="w-full md:w-auto"
        onClick={() => void query.refetch()}
        disabled={query.isFetching}
      >
        <RefreshCwIcon
          className={cn('size-4', query.isFetching && 'motion-safe:animate-spin')}
          aria-hidden
        />
        {t('overview.refresh')}
      </Button>
    </>
  )

  return (
    <PageContainer>
      <PageHeader
        title={t('overview.title')}
        description={t('overview.description')}
        actions={actions}
      />

      {query.isError && !overview ? (
        <Card>
          <CardContent className="flex flex-col items-start gap-3 py-2">
            <p className="text-sm text-destructive" role="alert">
              {t('overview.error')}
            </p>
            <Button variant="outline" onClick={() => void query.refetch()}>
              {t('overview.retry')}
            </Button>
          </CardContent>
        </Card>
      ) : (
        <>
          <DashboardStatGrid kpis={overview?.kpis} isLoading={query.isLoading} />
          <DashboardCharts
            weeklyActivity={overview?.weeklyActivity}
            isLoading={query.isLoading}
          />
        </>
      )}
    </PageContainer>
  )
}
