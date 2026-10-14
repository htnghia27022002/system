'use client'

import { MailIcon, MailOpenIcon, Trash2Icon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { PermissionGate } from '@/features/access-control'
import { cn } from '@/lib/utils'

import { WebhooksPermissions } from '../permissions'
import type { WebhookRequestListItem } from '../types/webhooks'

type WebhooksRequestListProps = {
  items: WebhookRequestListItem[]
  selectedId: string | null
  onSelect: (id: string) => void
  onDelete: (id: string) => void
  onToggleRead: (id: string, isRead: boolean) => void
  deletingId?: string | null
  togglingReadId?: string | null
  isLoading?: boolean
  emptyTitle: string
  emptyBody: string
  publicUrl?: string
}

function methodBadgeClass(method: string): string {
  const m = method.toUpperCase()
  if (m === 'GET') return 'bg-success/15 text-success'
  if (m === 'POST') return 'bg-primary/15 text-primary'
  if (m === 'PUT' || m === 'PATCH') return 'bg-warning/15 text-warning'
  if (m === 'DELETE') return 'bg-destructive/15 text-destructive'
  return 'bg-muted text-muted-foreground'
}

function pathSnippet(url: string): string {
  try {
    if (url.startsWith('http')) {
      const parsed = new URL(url)
      return `${parsed.pathname}${parsed.search}` || '/'
    }
  } catch {
    /* fall through */
  }
  return url.length > 72 ? `${url.slice(0, 69)}…` : url
}

function formatTimestamp(iso: string, locale: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString(locale, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

export function WebhooksRequestList({
  items,
  selectedId,
  onSelect,
  onDelete,
  onToggleRead,
  deletingId,
  togglingReadId,
  isLoading,
  emptyTitle,
  emptyBody,
  publicUrl,
}: WebhooksRequestListProps) {
  const { t, i18n } = useTranslation('common')

  if (isLoading && items.length === 0) {
    return (
      <div className="flex h-full min-h-48 items-center justify-center p-6 text-sm text-muted-foreground">
        {t('tools.webhooks.list.loading')}
      </div>
    )
  }

  if (items.length === 0) {
    return (
      <div className="flex h-full min-h-48 flex-col items-start justify-center gap-3 p-6">
        <h3 className="font-heading text-base font-semibold">{emptyTitle}</h3>
        <p className="max-w-prose text-sm text-muted-foreground">{emptyBody}</p>
        {publicUrl ? (
          <code className="max-w-full break-all rounded-md bg-muted px-2 py-1 text-xs">
            {publicUrl}
          </code>
        ) : null}
      </div>
    )
  }

  return (
    <ul className="h-full divide-y divide-border overflow-y-auto" role="listbox">
      {items.map((item) => {
        const selected = item.id === selectedId
        const unread = !item.isRead
        return (
          <li key={item.id} className="group relative">
            <button
              type="button"
              role="option"
              aria-selected={selected}
              onClick={() => onSelect(item.id)}
              className={cn(
                'flex w-full flex-col gap-1 px-3 py-2.5 pr-16 text-left transition-colors',
                'hover:bg-muted/70 focus-visible:bg-muted/70 focus-visible:outline-none',
                selected && 'bg-muted',
                unread && 'border-l-2 border-l-primary pl-[10px]',
                !unread && 'pl-3',
              )}
            >
              <div className="flex items-center gap-2">
                {unread ? (
                  <span
                    className="size-2 shrink-0 rounded-full bg-primary"
                    aria-hidden
                  />
                ) : (
                  <span className="size-2 shrink-0" aria-hidden />
                )}
                <Badge
                  variant="outline"
                  className={cn(
                    'border-0 font-mono text-xs',
                    methodBadgeClass(item.method),
                  )}
                >
                  {item.method.toUpperCase()}
                </Badge>
                <Badge
                  variant={unread ? 'default' : 'secondary'}
                  className="h-5 px-1.5 text-xs font-medium"
                >
                  {unread
                    ? t('tools.webhooks.badges.unread')
                    : t('tools.webhooks.badges.read')}
                </Badge>
                <span className="ml-auto shrink-0 text-xs text-muted-foreground tabular-nums">
                  {formatTimestamp(item.createdAt, i18n.language)}
                </span>
              </div>
              <p
                className={cn(
                  'truncate font-mono text-xs',
                  unread ? 'font-semibold text-foreground' : 'text-foreground/80',
                )}
              >
                {pathSnippet(item.url)}
              </p>
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                {item.clientIp ? (
                  <span className="shrink-0 font-mono">{item.clientIp}</span>
                ) : null}
                {item.snippet ? (
                  <span className="truncate opacity-80">{item.snippet}</span>
                ) : null}
              </div>
            </button>
            <div
              className={cn(
                'absolute top-1.5 right-1 flex gap-0.5 opacity-0 transition-opacity',
                'group-hover:opacity-100 focus-within:opacity-100',
                selected && 'opacity-100',
              )}
            >
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="size-7 text-muted-foreground"
                aria-label={
                  unread
                    ? t('tools.webhooks.actions.markRead')
                    : t('tools.webhooks.actions.markUnread')
                }
                disabled={togglingReadId === item.id}
                onClick={(e) => {
                  e.stopPropagation()
                  onToggleRead(item.id, unread)
                }}
              >
                {unread ? (
                  <MailOpenIcon className="size-3.5" />
                ) : (
                  <MailIcon className="size-3.5" />
                )}
              </Button>
              <PermissionGate permission={WebhooksPermissions.modify}>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="size-7 text-muted-foreground hover:text-destructive"
                  aria-label={t('tools.webhooks.actions.remove')}
                  disabled={deletingId === item.id}
                  onClick={(e) => {
                    e.stopPropagation()
                    onDelete(item.id)
                  }}
                >
                  <Trash2Icon className="size-3.5" />
                </Button>
              </PermissionGate>
            </div>
          </li>
        )
      })}
    </ul>
  )
}
