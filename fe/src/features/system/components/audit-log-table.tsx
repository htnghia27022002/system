'use client'

import { Suspense, useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { DataTable, DataTableSkeleton } from '@/components/common/data-table'
import { MobileRecordCard } from '@/components/common/mobile-record-card'
import { PageContainer, PageHeader } from '@/components/common/page-layout'
import { ServerListSearch } from '@/components/common/server-list-search'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useListQueryParams } from '@/hooks/use-list-query-params'

import { useAuditLogs } from '../hooks/use-system'
import {
  actionVariant,
  auditLogColumns,
  formatAuditTime,
  summarizeMetadata,
} from './audit-log-columns'

const RESOURCE_TYPES = ['user', 'role', 'session', 'settings'] as const

export function AuditLogTable() {
  return (
    <Suspense
      fallback={
        <PageContainer>
          <DataTableSkeleton columns={6} />
        </PageContainer>
      }
    >
      <AuditLogTableContent />
    </Suspense>
  )
}

function AuditLogTableContent() {
  const { t, i18n } = useTranslation('admin')
  const { values, setParams } = useListQueryParams(['search', 'resourceType'] as const)
  const page = Number(values.page) || 1
  const pageSize = Number(values.pageSize) || 10
  const query = useAuditLogs({
    page,
    pageSize,
    search: values.search || undefined,
    resourceType: values.resourceType || undefined,
  })
  const columns = useMemo(() => auditLogColumns(t, i18n.language), [t, i18n.language])

  const header = (
    <PageHeader title={t('system.audit.title')} description={t('system.audit.description')} />
  )

  if (query.isLoading) {
    return (
      <PageContainer>
        {header}
        <DataTableSkeleton columns={6} />
      </PageContainer>
    )
  }

  if (query.isError) {
    return (
      <PageContainer>
        {header}
        <Card>
          <CardContent className="flex flex-col gap-3 py-6">
            <p className="text-sm text-muted-foreground">{t('system.audit.error')}</p>
            <Button variant="outline" className="w-fit" onClick={() => void query.refetch()}>
              {t('access.actions.retry')}
            </Button>
          </CardContent>
        </Card>
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      {header}
      <DataTable
        columns={columns}
        data={query.data?.items ?? []}
        getRowId={(row) => row.id}
        localPagination={false}
        serverPaging={{
          page,
          pageSize,
          total: query.data?.total ?? 0,
          onPageChange: (next) => setParams({ page: String(next) }),
          onPageSizeChange: (next) => setParams({ pageSize: String(next), page: '1' }),
        }}
        isRefreshing={query.isFetching && !query.isLoading}
        emptyTitle={t('system.audit.emptyTitle')}
        emptyDescription={t('system.audit.emptyDescription')}
        toolbar={
          <div className="flex w-full flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
            <ServerListSearch
              value={values.search}
              placeholder={t('system.audit.searchPlaceholder')}
              onSearch={(next) => setParams({ search: next || undefined })}
              onClear={() => setParams({ search: undefined })}
            />
            <Select
              value={values.resourceType || 'all'}
              onValueChange={(value) => setParams({ resourceType: value === 'all' ? undefined : value })}
            >
              <SelectTrigger className="w-full sm:w-[11rem]" aria-label={t('system.audit.fields.resource')}>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{t('system.audit.allResources')}</SelectItem>
                {RESOURCE_TYPES.map((type) => (
                  <SelectItem key={type} value={type}>
                    {t(`system.audit.resources.${type}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        }
        renderMobileCard={(entry) => (
          <MobileRecordCard
            title={<Badge variant={actionVariant(entry.action)} className="font-mono">{entry.action}</Badge>}
            fields={[
              { label: t('system.audit.fields.time'), value: formatAuditTime(entry.createdAt, i18n.language) },
              { label: t('system.audit.fields.actor'), value: entry.actorEmail || t('system.audit.system') },
              {
                label: t('system.audit.fields.resource'),
                value: entry.resourceType ? `${entry.resourceType} ${entry.resourceId.slice(0, 8)}` : '—',
              },
              { label: t('system.audit.fields.details'), value: summarizeMetadata(entry.metadata) || '—' },
              { label: t('system.audit.fields.ip'), value: entry.ipAddress || '—' },
            ]}
          />
        )}
      />
    </PageContainer>
  )
}
