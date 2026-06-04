'use client'

import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { fetchHealth } from '@/services/health-api'

const chartData = [
  { name: 'Mon', value: 12 },
  { name: 'Tue', value: 19 },
  { name: 'Wed', value: 8 },
  { name: 'Thu', value: 15 },
  { name: 'Fri', value: 22 },
]

export function HomePage() {
  const { t } = useTranslation('common')
  const healthQuery = useQuery({
    queryKey: ['health'],
    queryFn: fetchHealth,
  })

  return (
    <div className="grid gap-6">
      <section className="space-y-2">
        <h1 className="text-3xl font-semibold tracking-tight">
          {t('home.title')}
        </h1>
        <p className="max-w-2xl text-muted-foreground">{t('home.description')}</p>
        <div className="flex flex-wrap gap-2">
          <Button onClick={() => toast.success(t('home.toast'))}>
            Show toast
          </Button>
          <Button variant="outline" onClick={() => void healthQuery.refetch()}>
            Refetch API health
          </Button>
        </div>
        <p className="text-sm text-muted-foreground">
          API status: {healthQuery.isLoading ? 'checking...' : healthQuery.data?.status}
        </p>
      </section>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Weekly activity</CardTitle>
            <CardDescription>Recharts sample chart</CardDescription>
          </CardHeader>
          <CardContent className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="value" fill="var(--chart-1)" radius={4} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Calendar</CardTitle>
            <CardDescription>react-day-picker via shadcn</CardDescription>
          </CardHeader>
          <CardContent className="flex justify-center">
            <Calendar mode="single" />
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
