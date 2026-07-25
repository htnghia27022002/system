'use client'

import {
  CheckIcon,
  CopyIcon,
  RefreshCwIcon,
  RotateCcwIcon,
} from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { PermissionGate, PermissionKeys } from '@/features/access-control'
import { useDebouncedValue } from '@/hooks/use-debounced-value'

import {
  useWebhookInbox,
  useWebhookInboxMutations,
  useWebhookRequestDetail,
  useWebhookRequests,
} from '../hooks/use-webhooks-inbox'
import type {
  WebhookMethodFilter,
  WebhookReadFilter,
} from '../types/webhooks'
import { buildWebhookPublicUrl } from '../utils/webhooks-url'
import { WebhooksListToolbar } from './webhooks-list-toolbar'
import { WebhooksRequestDetail } from './webhooks-request-detail'
import { WebhooksRequestList } from './webhooks-request-list'

const PAGE_SIZE = 20

export function WebhooksInbox() {
  const { t } = useTranslation('common')
  const [method, setMethod] = useState<WebhookMethodFilter>('ALL')
  const [readFilter, setReadFilter] = useState<WebhookReadFilter>('all')
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [copied, setCopied] = useState(false)

  const debouncedSearch = useDebouncedValue(search, 300)

  const listParams = useMemo(
    () => ({
      method: method === 'ALL' ? undefined : method,
      q: debouncedSearch.trim() || undefined,
      read: readFilter,
      page,
      limit: PAGE_SIZE,
    }),
    [method, debouncedSearch, readFilter, page],
  )

  const inboxQuery = useWebhookInbox()
  const listQuery = useWebhookRequests(listParams)
  const detailQuery = useWebhookRequestDetail(selectedId)
  const { regenerate, softDelete, clearAll, setRead, invalidateAll } =
    useWebhookInboxMutations()

  const inbox = inboxQuery.data
  const list = listQuery.data
  const items = useMemo(() => list?.items ?? [], [list?.items])

  const publicUrl = useMemo(() => {
    if (!inbox) return ''
    const origin =
      typeof window !== 'undefined' ? window.location.origin : ''
    return buildWebhookPublicUrl(origin, inbox.publicPath || inbox.publicUuid)
  }, [inbox])

  const activeCount = inbox?.activeCount ?? list?.activeCount ?? 0
  const lifetimeReceived =
    inbox?.lifetimeReceived ?? list?.lifetimeReceived ?? 0

  useEffect(() => {
    setPage(1)
  }, [method, debouncedSearch, readFilter])

  useEffect(() => {
    if (items.length === 0) {
      setSelectedId(null)
      return
    }
    if (selectedId && items.some((i) => i.id === selectedId)) return
    setSelectedId(items[0].id)
  }, [items, selectedId])

  async function handleCopy() {
    if (!publicUrl) return
    try {
      await navigator.clipboard.writeText(publicUrl)
      setCopied(true)
      toast.success(t('tools.webhooks.toasts.copied'))
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      toast.error(t('tools.webhooks.toasts.copyFailed'))
    }
  }

  function handleRefresh() {
    void invalidateAll()
    void inboxQuery.refetch()
    void listQuery.refetch()
    if (selectedId) void detailQuery.refetch()
  }

  function handleSoftDelete(id?: string | null) {
    const target = id ?? selectedId
    if (!target) return
    softDelete.mutate(target, {
      onSuccess: () => {
        if (selectedId === target) setSelectedId(null)
      },
    })
  }

  function handleClearAll() {
    clearAll.mutate(undefined, {
      onSuccess: () => setSelectedId(null),
    })
  }

  if (inboxQuery.isLoading && !inbox) {
    return (
      <div className="space-y-3">
        <div className="rounded-2xl border border-border bg-card p-3">
          <Skeleton className="h-8 w-full" />
        </div>
        <Skeleton className="h-[40rem] w-full rounded-2xl" />
      </div>
    )
  }

  if (inboxQuery.isError || !inbox) {
    return (
      <div className="rounded-2xl border border-destructive/40 bg-destructive/5 px-6 py-10 text-sm text-destructive">
        {t('tools.webhooks.errors.inboxLoad')}
      </div>
    )
  }

  const isRefreshing = listQuery.isFetching || inboxQuery.isFetching
  const copyLabel = copied
    ? t('tools.webhooks.actions.copied')
    : t('tools.webhooks.actions.copy')

  return (
    <div className="space-y-3">
      <section
        className="rounded-2xl border border-border bg-card p-3"
        aria-label={t('tools.webhooks.urlTitle')}
      >
        <p id="webhooks-url-hint" className="sr-only">
          {t('tools.webhooks.urlHint')}
        </p>
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
          <Input
            readOnly
            value={publicUrl}
            aria-label={t('tools.webhooks.publicUrlLabel')}
            aria-describedby="webhooks-url-hint"
            className="h-8 min-w-0 flex-1 font-mono text-xs"
          />
          <div className="flex shrink-0 flex-wrap items-center gap-1.5">
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  size="icon-sm"
                  variant="outline"
                  onClick={() => void handleCopy()}
                  aria-label={copyLabel}
                >
                  {copied ? (
                    <CheckIcon className="size-3.5" />
                  ) : (
                    <CopyIcon className="size-3.5" />
                  )}
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">{copyLabel}</TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  type="button"
                  size="icon-sm"
                  variant="outline"
                  onClick={handleRefresh}
                  disabled={isRefreshing}
                  aria-label={t('tools.webhooks.actions.refresh')}
                >
                  <RefreshCwIcon
                    className={
                      isRefreshing ? 'size-3.5 animate-spin' : 'size-3.5'
                    }
                  />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">
                {t('tools.webhooks.actions.refresh')}
              </TooltipContent>
            </Tooltip>
            <PermissionGate permission={PermissionKeys.webhooks.modify}>
              <AlertDialog>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <AlertDialogTrigger asChild>
                      <Button
                        type="button"
                        size="icon-sm"
                        variant="outline"
                        aria-label={t('tools.webhooks.actions.regenerate')}
                      >
                        <RotateCcwIcon className="size-3.5" />
                      </Button>
                    </AlertDialogTrigger>
                  </TooltipTrigger>
                  <TooltipContent side="bottom">
                    {t('tools.webhooks.actions.regenerate')}
                  </TooltipContent>
                </Tooltip>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>
                      {t('tools.webhooks.regenerate.title')}
                    </AlertDialogTitle>
                    <AlertDialogDescription>
                      {t('tools.webhooks.regenerate.description')}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>
                      {t('tools.webhooks.actions.cancel')}
                    </AlertDialogCancel>
                    <AlertDialogAction
                      onClick={() => regenerate.mutate()}
                      disabled={regenerate.isPending}
                    >
                      {t('tools.webhooks.actions.confirmRegenerate')}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </PermissionGate>
            <p className="ml-1 text-xs whitespace-nowrap text-muted-foreground">
              {t('tools.webhooks.totals', {
                active: activeCount,
                lifetime: lifetimeReceived,
              })}
            </p>
          </div>
        </div>
      </section>

      <section
        className="overflow-hidden rounded-2xl border border-border bg-card"
        aria-label={t('tools.webhooks.inboxLabel')}
      >
        <div className="grid min-h-[40rem] lg:grid-cols-[minmax(18rem,26rem)_minmax(0,1fr)] lg:h-[min(78vh,52rem)]">
          <div className="flex min-h-0 flex-col border-b border-border lg:border-r lg:border-b-0">
            <WebhooksListToolbar
              search={search}
              onSearchChange={setSearch}
              method={method}
              onMethodChange={setMethod}
              readFilter={readFilter}
              onReadFilterChange={setReadFilter}
              shown={items.length}
              total={list?.total ?? 0}
              page={page}
              activeCount={activeCount}
              onClearAll={handleClearAll}
              clearAllPending={clearAll.isPending}
            />
            <div className="min-h-0 flex-1 overflow-hidden">
              <WebhooksRequestList
                items={items}
                selectedId={selectedId}
                onSelect={setSelectedId}
                onDelete={(id) => handleSoftDelete(id)}
                onToggleRead={(id, isRead) =>
                  setRead.mutate({ id, isRead })
                }
                deletingId={softDelete.isPending ? softDelete.variables : null}
                togglingReadId={
                  setRead.isPending ? setRead.variables?.id : null
                }
                isLoading={listQuery.isLoading}
                emptyTitle={t('tools.webhooks.empty.title')}
                emptyBody={t('tools.webhooks.empty.body')}
                publicUrl={publicUrl}
              />
            </div>
            {(list?.hasMore || page > 1) && (
              <div className="flex items-center justify-between gap-2 border-t border-border px-3 py-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={page <= 1 || listQuery.isFetching}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  {t('tools.webhooks.pagination.prev')}
                </Button>
                <span className="text-xs text-muted-foreground">
                  {t('tools.webhooks.pagination.page', { page })}
                </span>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  disabled={!list?.hasMore || listQuery.isFetching}
                  onClick={() => setPage((p) => p + 1)}
                >
                  {t('tools.webhooks.pagination.next')}
                </Button>
              </div>
            )}
          </div>
          <div className="min-h-0 overflow-hidden">
            <WebhooksRequestDetail
              detail={detailQuery.data}
              isLoading={detailQuery.isLoading}
              isError={detailQuery.isError}
            />
          </div>
        </div>
      </section>
    </div>
  )
}
