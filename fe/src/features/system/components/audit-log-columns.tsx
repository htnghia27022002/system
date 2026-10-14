'use client'

import type { ColumnDef } from '@tanstack/react-table'
import type { TFunction } from 'i18next'

import { DataTableColumnHeader } from '@/components/common/data-table'
import { Badge } from '@/components/ui/badge'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

import type { AuditLogEntry } from '../types'

/** Failed or destructive actions get the destructive badge; the rest stay neutral. */
export function actionVariant(action: string): 'destructive' | 'secondary' {
  return action.endsWith('_failed') || action.endsWith('.delete') ? 'destructive' : 'secondary'
}

export function formatAuditTime(iso: string, locale: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

export function summarizeMetadata(metadata: Record<string, unknown>): string {
  return Object.entries(metadata)
    .map(([key, value]) => `${key}: ${Array.isArray(value) ? value.join(', ') : String(value)}`)
    .join(' · ')
}

export function auditLogColumns(t: TFunction<'admin'>, locale: string): ColumnDef<AuditLogEntry>[] {
  return [
    {
      accessorKey: 'createdAt',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.time')} />,
      cell: ({ row }) => (
        <span className="whitespace-nowrap text-muted-foreground tabular-nums">
          {formatAuditTime(row.original.createdAt, locale)}
        </span>
      ),
    },
    {
      accessorKey: 'actorEmail',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.actor')} />,
      cell: ({ row }) => row.original.actorEmail || <span className="text-muted-foreground">{t('system.audit.system')}</span>,
    },
    {
      accessorKey: 'action',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.action')} />,
      cell: ({ row }) => (
        <Badge variant={actionVariant(row.original.action)} className="font-mono">
          {row.original.action}
        </Badge>
      ),
    },
    {
      id: 'resource',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.resource')} />,
      cell: ({ row }) => {
        const { resourceType, resourceId } = row.original
        if (!resourceType) return <span className="text-muted-foreground">—</span>
        return (
          <span className="font-mono text-xs">
            {resourceType}
            {resourceId ? <span className="text-muted-foreground"> · {resourceId.slice(0, 8)}</span> : null}
          </span>
        )
      },
    },
    {
      id: 'details',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.details')} />,
      cell: ({ row }) => {
        const summary = summarizeMetadata(row.original.metadata)
        return summary ? (
          <span className="line-clamp-1 max-w-64 text-xs text-muted-foreground">{summary}</span>
        ) : (
          <span className="text-muted-foreground">—</span>
        )
      },
    },
    {
      accessorKey: 'ipAddress',
      enableSorting: false,
      header: ({ column }) => <DataTableColumnHeader column={column} title={t('system.audit.fields.ip')} />,
      cell: ({ row }) =>
        row.original.ipAddress ? (
          <Tooltip>
            <TooltipTrigger asChild>
              <span className="cursor-default font-mono text-xs">{row.original.ipAddress}</span>
            </TooltipTrigger>
            <TooltipContent className="max-w-80 break-words">
              {row.original.userAgent || '—'}
              {row.original.requestId ? ` · ${row.original.requestId}` : ''}
            </TooltipContent>
          </Tooltip>
        ) : (
          <span className="text-muted-foreground">—</span>
        ),
    },
  ]
}
