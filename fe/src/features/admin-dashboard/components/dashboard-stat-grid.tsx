'use client'

import { MinusIcon, TrendingDownIcon, TrendingUpIcon } from 'lucide-react'
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
  const { t, i18n } = useTranslation('admin')

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: 4 }, (_, index) => (
          <Card key={index} className="gap-3">
            <CardHeader className="px-4">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-8 w-20" />
            </CardHeader>
            <CardContent className="px-4">
              <Skeleton className="h-3 w-28" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      {kpis?.map((kpi) => {
        const trend = Math.sign(kpi.changePercent)
        const TrendIcon =
          trend > 0 ? TrendingUpIcon : trend < 0 ? TrendingDownIcon : MinusIcon
        return (
          <Card key={kpi.key} className="gap-3">
            <CardHeader className="px-4">
              <CardDescription className="font-medium">
                {t(`overview.kpis.${kpi.key}`)}
              </CardDescription>
              <CardTitle className="text-2xl font-semibold tabular-nums">
                {formatKpiValue(kpi.value, i18n.language)}
              </CardTitle>
            </CardHeader>
            <CardContent className="px-4">
              <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <span
                  className={cn(
                    'inline-flex items-center gap-1 font-medium tabular-nums',
                    trend > 0 && 'text-success',
                    trend < 0 && 'text-destructive',
                  )}
                >
                  <TrendIcon className="size-4" aria-hidden />
                  {formatChange(kpi.changePercent, i18n.language)}
                </span>
                {t('overview.vsLastWeek')}
              </p>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}

function formatKpiValue(value: number, locale: string): string {
  return new Intl.NumberFormat(locale).format(value)
}

function formatChange(changePercent: number, locale: string): string {
  return new Intl.NumberFormat(locale, {
    style: 'percent',
    signDisplay: 'exceptZero',
    minimumFractionDigits: 1,
    maximumFractionDigits: 1,
  }).format(changePercent / 100)
}
