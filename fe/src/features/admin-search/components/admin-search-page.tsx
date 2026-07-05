'use client'

import { AlertTriangleIcon, SearchIcon } from 'lucide-react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { useEffect, useMemo, useState, type FormEvent } from 'react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/lib/utils'

import { useAdminSearch } from '../hooks/use-admin-search'
import { useSearchNavigation } from '../hooks/use-search-navigation'
import type { SearchEntityType } from '../types'
import { SearchHitList } from './search-hit-list'

const ENTITY_TYPES: Array<{ value: 'all' | SearchEntityType; labelKey: string }> = [
  { value: 'all', labelKey: 'search.filters.all' },
  { value: 'user', labelKey: 'search.filters.users' },
  { value: 'role', labelKey: 'search.filters.roles' },
  { value: 'permission', labelKey: 'search.filters.permissions' },
]

export function AdminSearchPage() {
  const { t } = useTranslation('admin')
  const router = useRouter()
  const searchParams = useSearchParams()
  const { access } = useSearchNavigation()

  const qParam = searchParams.get('q') ?? ''
  const typesParam = searchParams.get('types') ?? 'all'
  const pageParam = Number(searchParams.get('page') ?? '1')

  const availableEntityTypes = useMemo(() => {
    return ENTITY_TYPES.filter((item) => {
      if (item.value === 'all') return true
      if (item.value === 'user') return access.canViewUsers
      if (item.value === 'role') return access.canViewRoles
      if (item.value === 'permission') return access.canViewPermissions
      return false
    })
  }, [access.canViewPermissions, access.canViewRoles, access.canViewUsers])

  const defaultEntityFilter =
    availableEntityTypes.find((item) => item.value === typesParam)?.value ??
    availableEntityTypes[0]?.value ??
    'all'

  const [query, setQuery] = useState(qParam)
  const [entityFilter, setEntityFilter] = useState(defaultEntityFilter)

  useEffect(() => {
    setQuery(qParam)
  }, [qParam])

  useEffect(() => {
    setEntityFilter(defaultEntityFilter)
  }, [defaultEntityFilter])

  const searchParamsMemo = useMemo(() => {
    const types =
      entityFilter === 'all' ? undefined : ([entityFilter] as SearchEntityType[])

    return {
      q: qParam,
      types,
      page: Number.isFinite(pageParam) && pageParam > 0 ? pageParam : 1,
      pageSize: 20,
    }
  }, [entityFilter, pageParam, qParam])

  const { data, isLoading, isError, refetch, isFetching } = useAdminSearch(searchParamsMemo)

  function pushSearch(next: { q?: string; types?: string; page?: number }) {
    const params = new URLSearchParams()
    const nextQuery = next.q ?? qParam
    const nextTypes = next.types ?? entityFilter
    const nextPage = next.page ?? 1

    if (nextQuery.trim()) {
      params.set('q', nextQuery.trim())
    }
    if (nextTypes !== 'all') {
      params.set('types', nextTypes)
    }
    if (nextPage > 1) {
      params.set('page', String(nextPage))
    }

    const suffix = params.toString()
    router.push(suffix ? `/admin/search?${suffix}` : '/admin/search')
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    pushSearch({ q: query, types: entityFilter, page: 1 })
  }

  const pagination = data?.pagination
  const canGoPrev = (pagination?.page ?? 1) > 1
  const canGoNext =
    pagination != null && pagination.page < pagination.totalPages

  const hasQuery = Boolean(qParam.trim())
  const showInitialLoading = hasQuery && isLoading
  const showResults = hasQuery && !isLoading && !isError && data
  const isRefetching = hasQuery && isFetching && !isLoading && Boolean(data)

  return (
    <div className="flex flex-1 flex-col gap-6 p-4 md:p-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">{t('search.title')}</h1>
        <p className="text-sm text-muted-foreground">{t('search.description')}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4 rounded-xl border bg-card p-4">
        <div className="grid gap-4 md:grid-cols-[minmax(0,1fr)_220px_auto] md:items-end">
          <div className="space-y-2">
            <Label htmlFor="admin-search-query">{t('search.queryLabel')}</Label>
            <div className="relative">
              <SearchIcon
                className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
                aria-hidden
              />
              <Input
                id="admin-search-query"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder={t('search.placeholder')}
                className="pl-9"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="admin-search-type">{t('search.typeLabel')}</Label>
            <Select value={entityFilter} onValueChange={setEntityFilter}>
              <SelectTrigger id="admin-search-type" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {availableEntityTypes.map((item) => (
                  <SelectItem key={item.value} value={item.value}>
                    {t(item.labelKey)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <Button type="submit" className="w-full md:w-auto">
            {t('search.submit')}
          </Button>
        </div>
      </form>

      {data?.degraded ? (
        <div
          className="flex gap-3 rounded-lg border border-amber-500/30 bg-amber-500/10 p-4 text-sm"
          role="status"
        >
          <AlertTriangleIcon className="mt-0.5 size-4 shrink-0 text-amber-600 dark:text-amber-400" />
          <div className="space-y-1">
            <p className="font-medium">{t('search.degradedTitle')}</p>
            <p className="text-muted-foreground">{t('search.degradedDescription')}</p>
          </div>
        </div>
      ) : null}

      {!hasQuery ? (
        <p className="rounded-lg border border-dashed px-4 py-10 text-center text-sm text-muted-foreground">
          {t('search.prompt')}
        </p>
      ) : null}

      {showInitialLoading ? (
        <div className="flex items-center justify-center gap-2 py-16 text-sm text-muted-foreground">
          <Spinner className="size-4" />
          {t('search.loading')}
        </div>
      ) : null}

      {hasQuery && isError ? (
        <div className="space-y-3 rounded-lg border border-destructive/30 bg-destructive/5 p-4">
          <p className="text-sm text-destructive">{t('search.error')}</p>
          <Button type="button" variant="outline" size="sm" onClick={() => void refetch()}>
            {t('overview.retry')}
          </Button>
        </div>
      ) : null}

      {showResults ? (
        <div className="relative space-y-4">
          {isRefetching ? (
            <div
              className="pointer-events-none absolute inset-0 z-10 flex items-start justify-center rounded-lg bg-background/50 pt-8"
              aria-hidden
            >
              <Spinner className="size-4 text-muted-foreground" />
            </div>
          ) : null}

          <div
            className={cn(
              'flex flex-wrap items-center justify-between gap-2 text-sm text-muted-foreground',
              isRefetching && 'opacity-60',
            )}
          >
            <p>{t('search.resultsCount', { count: data.pagination.total })}</p>
            {pagination && pagination.totalPages > 1 ? (
              <p>
                {t('search.pageSummary', {
                  page: pagination.page,
                  totalPages: pagination.totalPages,
                })}
              </p>
            ) : null}
          </div>

          <SearchHitList hits={data.hits} emptyMessage={t('search.empty')} />

          {pagination && pagination.totalPages > 1 ? (
            <div className="flex items-center justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={!canGoPrev || isRefetching}
                onClick={() => pushSearch({ page: pagination.page - 1 })}
              >
                {t('search.previous')}
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={!canGoNext || isRefetching}
                onClick={() => pushSearch({ page: pagination.page + 1 })}
              >
                {t('search.next')}
              </Button>
            </div>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}
