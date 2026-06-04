'use client'

import { useTranslation } from 'react-i18next'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

import type { DashboardKpi } from '../types'

type DashboardStatGridProps = {
  kpis?: DashboardKpi[]
  isLoading: boolean
}

export function DashboardStatGrid({ kpis, isLoading }: DashboardStatGridProps) {
  const { t } = useTranslation('admin')

  const placeholders = Array.from({ length: 4 }, (_, index) => index)

  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      {isLoading
        ? placeholders.map((key) => (
            <Card
              key={key}
              className="gap-3 border border-border bg-card py-4 shadow-none ring-0"
            >
              <CardHeader className="px-4">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-8 w-20" />
              </CardHeader>
              <CardContent className="px-4">
                <Skeleton className="h-3 w-16" />
              </CardContent>
            </Card>
          ))
        : kpis?.map((kpi) => {
            const positive = kpi.changePercent >= 0
            return (
              <Card
                key={kpi.key}
                className="gap-3 border border-border bg-card py-4 shadow-none ring-0"
              >
                <CardHeader className="px-4">
                  <CardDescription className="text-[14px] font-semibold leading-[1.29] tracking-[-0.224px] text-muted-foreground">
                    {t(`overview.kpis.${kpi.key}`)}
                  </CardDescription>
                  <CardTitle className="text-[28px] font-semibold tabular-nums">
                    {formatKpiValue(kpi.key, kpi.value)}
                  </CardTitle>
                </CardHeader>
                <CardContent className="px-4">
                  <p
                    className={cn(
                      'text-[14px] leading-[1.29] tracking-[-0.224px]',
                      positive ? 'text-primary' : 'text-destructive',
                    )}
                  >
                    {positive ? '+' : ''}
                    {kpi.changePercent.toFixed(1)}% {t('overview.vsLastWeek')}
                  </p>
                </CardContent>
              </Card>
            )
          })}
    </div>
  )
}

function formatKpiValue(key: DashboardKpi['key'], value: number): string {
  if (key === 'revenue') {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 0,
    }).format(value)
  }

  return new Intl.NumberFormat('en-US').format(value)
}
