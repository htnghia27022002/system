'use client'

import { useTranslation } from 'react-i18next'
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

import type { DashboardChartPoint } from '../types'

type DashboardChartsProps = {
  weeklyActivity?: DashboardChartPoint[]
  isLoading: boolean
}

export function DashboardCharts({
  weeklyActivity,
  isLoading,
}: DashboardChartsProps) {
  const { t } = useTranslation('admin')

  return (
    <Card className="mt-6 gap-0 border border-border bg-card py-0 shadow-none ring-0">
      <CardHeader className="px-6 pt-6">
        <CardTitle className="text-[21px] font-semibold leading-[1.19] tracking-[0.231px]">
          {t('overview.chart.title')}
        </CardTitle>
        <CardDescription className="text-[17px] leading-[1.47] tracking-[-0.374px]">
          {t('overview.chart.description')}
        </CardDescription>
      </CardHeader>
      <CardContent className="h-72 px-6 pb-6 pt-2">
        {isLoading ? (
          <Skeleton className="size-full rounded-lg" />
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={weeklyActivity ?? []}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="label" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="value" fill="var(--chart-1)" radius={4} />
            </BarChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
