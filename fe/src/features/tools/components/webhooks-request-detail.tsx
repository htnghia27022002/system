'use client'

import { ChevronDownIcon } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Badge } from '@/components/ui/badge'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import { Skeleton } from '@/components/ui/skeleton'
import { cn } from '@/lib/utils'

import type { WebhookRequestDetail } from '../types/webhooks'

type WebhooksRequestDetailProps = {
  detail: WebhookRequestDetail | undefined
  isLoading: boolean
  isError: boolean
}

function countEntries(
  value: Record<string, unknown> | string | null | undefined,
): number {
  if (value == null) return 0
  if (typeof value === 'string') return value.length > 0 ? 1 : 0
  return Object.keys(value).length
}

function CollapsibleBlock({
  title,
  value,
  defaultOpen = false,
  emptyLabel,
  badge,
}: {
  title: string
  value: Record<string, unknown> | string | null | undefined
  defaultOpen?: boolean
  emptyLabel: string
  badge?: string | null
}) {
  const [open, setOpen] = useState(defaultOpen)
  const empty =
    value == null ||
    (typeof value === 'object' && Object.keys(value).length === 0) ||
    (typeof value === 'string' && value.length === 0)

  let text = ''
  if (!empty) {
    text =
      typeof value === 'string' ? value : JSON.stringify(value, null, 2)
  }

  const count = countEntries(value)

  return (
    <Collapsible
      open={open}
      onOpenChange={setOpen}
      className="rounded-lg border border-border"
    >
      <CollapsibleTrigger asChild>
        <button
          type="button"
          className="flex w-full items-center gap-2 px-3 py-2.5 text-left hover:bg-muted/50"
        >
          <ChevronDownIcon
            className={cn(
              'size-4 shrink-0 text-muted-foreground transition-transform',
              open && 'rotate-180',
            )}
          />
          <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
            {title}
          </span>
          {badge ? (
            <Badge variant="outline" className="h-5 px-1.5 text-[10px]">
              {badge}
            </Badge>
          ) : null}
          <Badge
            variant="secondary"
            className="ml-auto h-5 px-1.5 text-[10px] tabular-nums"
          >
            {empty ? '0' : count}
          </Badge>
        </button>
      </CollapsibleTrigger>
      <CollapsibleContent>
        <div className="border-t border-border px-3 py-3">
          {empty ? (
            <p className="text-sm text-muted-foreground">{emptyLabel}</p>
          ) : (
            <pre className="max-h-[min(28rem,50vh)] overflow-auto rounded-md bg-muted/40 p-3 font-mono text-xs leading-relaxed whitespace-pre-wrap break-all">
              {text}
            </pre>
          )}
        </div>
      </CollapsibleContent>
    </Collapsible>
  )
}

export function WebhooksRequestDetail({
  detail,
  isLoading,
  isError,
}: WebhooksRequestDetailProps) {
  const { t } = useTranslation('common')

  if (!detail && !isLoading && !isError) {
    return (
      <div className="flex h-full min-h-48 items-center justify-center p-6 text-sm text-muted-foreground">
        {t('tools.webhooks.detail.selectPrompt')}
      </div>
    )
  }

  if (isLoading && !detail) {
    return (
      <div className="space-y-3 p-4">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-40 w-full" />
      </div>
    )
  }

  if (isError || !detail) {
    return (
      <div className="flex h-full min-h-48 items-center justify-center p-6 text-sm text-destructive">
        {t('tools.webhooks.detail.loadError')}
      </div>
    )
  }

  const bodyBadge =
    detail.bodyEncoding && detail.bodyEncoding !== 'utf-8'
      ? detail.bodyEncoding
      : detail.isBinary
        ? t('tools.webhooks.detail.binary')
        : detail.bodyTruncated
          ? t('tools.webhooks.detail.truncated')
          : null

  return (
    <div className="flex h-full min-h-0 flex-col overflow-y-auto px-3 py-3 sm:px-4">
      <div className="space-y-2">
        <CollapsibleBlock
          title={t('tools.webhooks.detail.body')}
          value={
            detail.isBinary
              ? t('tools.webhooks.detail.binaryBodyHint')
              : detail.body || null
          }
          defaultOpen
          emptyLabel="—"
          badge={bodyBadge}
        />
        <CollapsibleBlock
          title={t('tools.webhooks.detail.query')}
          value={detail.query}
          emptyLabel="—"
        />
        <CollapsibleBlock
          title={t('tools.webhooks.detail.form')}
          value={detail.form}
          emptyLabel="—"
        />
        <CollapsibleBlock
          title={t('tools.webhooks.detail.headers')}
          value={detail.headers}
          emptyLabel="—"
        />
      </div>
    </div>
  )
}
