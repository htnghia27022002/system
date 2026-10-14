'use client'

import { useMemo } from 'react'
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
import { useReducedMotion } from '@/hooks/use-reduced-motion'

import type { DashboardChartPoint } from '../types'

type DashboardChartsProps = {
  weeklyActivity?: DashboardChartPoint[]
  isLoading: boolean
}

const AXIS_TICK = { fill: 'var(--muted-foreground)', fontSize: 12 }

/** Parse `YYYY-MM-DD` as a local calendar day (not UTC midnight). */
function toLocalDate(isoDay: string): Date {
  const [year, month, day] = isoDay.split('-').map(Number)
  return new Date(year, month - 1, day)
}

export function DashboardCharts({
  weeklyActivity,
  isLoading,
}: DashboardChartsProps) {
  const { t, i18n } = useTranslation('admin')
  const reducedMotion = useReducedMotion()

  const rows = useMemo(() => {
    const weekday = new Intl.DateTimeFormat(i18n.language, { weekday: 'short' })
    const fullDay = new Intl.DateTimeFormat(i18n.language, {
      weekday: 'long',
      day: 'numeric',
      month: 'short',
    })
    return (weeklyActivity ?? []).map((point) => {
      const date = toLocalDate(point.date)
      return {
        ...point,
        tick: weekday.format(date),
        fullLabel: fullDay.format(date),
      }
    })
  }, [weeklyActivity, i18n.language])

  const number = new Intl.NumberFormat(i18n.language)

  return (
    <Card className="gap-0">
      <CardHeader className="px-4">
        <CardTitle className="text-base font-semibold">
          {t('overview.chart.title')}
        </CardTitle>
        <CardDescription>{t('overview.chart.description')}</CardDescription>
      </CardHeader>
      <CardContent className="h-72 px-2 pt-4 sm:px-4">
        {isLoading ? (
          <Skeleton className="size-full rounded-lg" />
        ) : rows.length === 0 ? (
          <p className="flex size-full items-center justify-center text-sm text-muted-foreground">
            {t('overview.chart.empty')}
          </p>
        ) : (
          <>
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={rows}
                margin={{ top: 4, right: 4, bottom: 0, left: 0 }}
                accessibilityLayer
              >
                <CartesianGrid vertical={false} stroke="var(--border)" />
                <XAxis
                  dataKey="tick"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  tick={AXIS_TICK}
                />
                <YAxis
                  width={36}
                  allowDecimals={false}
                  tickLine={false}
                  axisLine={false}
                  tick={AXIS_TICK}
                />
                <Tooltip
                  cursor={{ fill: 'var(--muted)', opacity: 0.6 }}
                  separator=": "
                  labelFormatter={(_, payload) =>
                    payload?.[0]?.payload?.fullLabel ?? ''
                  }
                  formatter={(value) => [
                    number.format(Number(value)),
                    t('overview.chart.series'),
                  ]}
                  contentStyle={{
                    background: 'var(--popover)',
                    border: '1px solid var(--border)',
                    borderRadius: 'var(--radius)',
                    color: 'var(--popover-foreground)',
                    fontSize: 12,
                  }}
                  labelStyle={{ color: 'var(--muted-foreground)' }}
                  itemStyle={{ color: 'var(--popover-foreground)' }}
                />
                <Bar
                  dataKey="value"
                  fill="var(--chart-1)"
                  radius={[4, 4, 0, 0]}
                  maxBarSize={40}
                  isAnimationActive={!reducedMotion}
                />
              </BarChart>
            </ResponsiveContainer>

            <table className="sr-only">
              <caption>{t('overview.chart.title')}</caption>
              <thead>
                <tr>
                  <th scope="col">{t('overview.chart.day')}</th>
                  <th scope="col">{t('overview.chart.series')}</th>
                </tr>
              </thead>
              <tbody>
                {rows.map((row) => (
                  <tr key={row.date}>
                    <th scope="row">{row.fullLabel}</th>
                    <td>{number.format(row.value)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        )}
      </CardContent>
    </Card>
  )
}
