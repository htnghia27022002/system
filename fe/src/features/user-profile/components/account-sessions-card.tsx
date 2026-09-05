'use client'

import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { ColumnDef } from '@tanstack/react-table'

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
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import {
  DataTable,
  DataTableColumnHeader,
  DataTableSkeleton,
} from '@/components/common/data-table'

import { useAccountSessions } from '../hooks/use-account-sessions'
import type { AccountSession } from '../types'

function formatSessionDate(iso: string, locale: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

function RevokeSessionButton({
  disabled,
  onConfirm,
}: {
  disabled: boolean
  onConfirm: () => void
}) {
  const { t } = useTranslation('admin')
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="text-destructive hover:text-destructive"
          disabled={disabled}
        >
          {t('profile.sessions.revoke')}
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t('profile.sessions.revokeTitle')}</AlertDialogTitle>
          <AlertDialogDescription>
            {t('profile.sessions.revokeDescription')}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t('profile.sessions.cancel')}</AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            onClick={onConfirm}
            disabled={disabled}
          >
            {t('profile.sessions.revokeConfirm')}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export function AccountSessionsCard() {
  const { t, i18n } = useTranslation('admin')
  const locale = i18n.language
  const { query, items, hasOtherSessions, revokeSession, revokeOthers } =
    useAccountSessions()
  const [revokeOthersOpen, setRevokeOthersOpen] = useState(false)
  const unknown = t('profile.sessions.unknown')

  const columns = useMemo<ColumnDef<AccountSession>[]>(
    () => [
      {
        accessorKey: 'current',
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title={t('profile.sessions.device')} />
        ),
        cell: ({ row }) => (
          <div
            data-testid={`session-row-${row.original.id}`}
            data-current={row.original.current ? 'true' : 'false'}
            className="flex items-center gap-2"
          >
            {row.original.current ? (
              <Badge variant="secondary">{t('profile.sessions.current')}</Badge>
            ) : (
              <span className="text-sm text-muted-foreground">
                {t('profile.sessions.otherDevice')}
              </span>
            )}
          </div>
        ),
      },
      {
        accessorKey: 'userAgent',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('profile.sessions.userAgent')}
          />
        ),
        cell: ({ row }) => {
          const value = row.original.userAgent?.trim() || unknown
          return (
            <span className="line-clamp-2 max-w-xl break-all text-sm" title={value}>
              {value}
            </span>
          )
        },
      },
      {
        accessorKey: 'ipAddress',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('profile.sessions.networkAddress')}
          />
        ),
        cell: ({ row }) => (
          <span className="text-sm">{row.original.ipAddress?.trim() || unknown}</span>
        ),
      },
      {
        accessorKey: 'lastUsedAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('profile.sessions.lastActivity')}
          />
        ),
        cell: ({ row }) => (
          <span className="text-sm whitespace-nowrap">
            {formatSessionDate(row.original.lastUsedAt, locale)}
          </span>
        ),
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('profile.sessions.started')}
          />
        ),
        cell: ({ row }) => (
          <span className="text-sm whitespace-nowrap">
            {formatSessionDate(row.original.createdAt, locale)}
          </span>
        ),
      },
      {
        accessorKey: 'expiresAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('profile.sessions.expires')}
          />
        ),
        cell: ({ row }) => (
          <span className="text-sm whitespace-nowrap">
            {formatSessionDate(row.original.expiresAt, locale)}
          </span>
        ),
      },
      {
        id: 'actions',
        enableSorting: false,
        header: () => (
          <span className="sr-only">{t('profile.sessions.actions')}</span>
        ),
        cell: ({ row }) =>
          row.original.current ? null : (
            <RevokeSessionButton
              disabled={revokeSession.isPending}
              onConfirm={() => revokeSession.mutate(row.original.id)}
            />
          ),
      },
    ],
    [locale, revokeSession, t, unknown],
  )

  if (query.isLoading && items.length === 0) {
    return <DataTableSkeleton columns={6} />
  }

  if (query.isError && items.length === 0) {
    return (
      <div className="flex flex-col gap-3">
        <p className="text-sm text-destructive" role="alert">
          {t('profile.sessions.loadFailed')}
        </p>
        <Button type="button" variant="outline" className="w-fit" onClick={() => void query.refetch()}>
          {t('profile.sessions.retry')}
        </Button>
      </div>
    )
  }

  return (
    <div className="flex min-w-0 flex-col gap-3">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <p className="max-w-3xl text-sm text-muted-foreground">
          {items.length > 0
            ? `${t('profile.sessions.count', { count: items.length })}. ${t('profile.sessions.description')}`
            : t('profile.sessions.description')}
        </p>
        {hasOtherSessions ? (
          <AlertDialog open={revokeOthersOpen} onOpenChange={setRevokeOthersOpen}>
            <AlertDialogTrigger asChild>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="shrink-0 text-destructive hover:text-destructive"
                disabled={revokeOthers.isPending}
              >
                {t('profile.sessions.revokeOthers')}
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>
                  {t('profile.sessions.revokeOthersTitle')}
                </AlertDialogTitle>
                <AlertDialogDescription>
                  {t('profile.sessions.revokeOthersDescription')}
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>{t('profile.sessions.cancel')}</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  disabled={revokeOthers.isPending}
                  onClick={() => {
                    revokeOthers.mutate(undefined, {
                      onSuccess: () => setRevokeOthersOpen(false),
                    })
                  }}
                >
                  {t('profile.sessions.revokeOthersConfirm')}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        ) : null}
      </div>

      <DataTable
        columns={columns}
        data={items}
        getRowId={(row) => row.id}
        filterKey="userAgent"
        filterPlaceholder={t('profile.sessions.searchPlaceholder')}
        emptyTitle={t('profile.sessions.empty')}
        isRefreshing={query.isFetching && items.length > 0}
        renderMobileCard={(row) => (
          <Card
            data-testid={`session-row-${row.id}`}
            data-current={row.current ? 'true' : 'false'}
            className="overflow-hidden py-0"
          >
            <CardHeader className="flex flex-row items-start justify-between gap-3 space-y-0 border-b px-4 py-3">
              <div className="min-w-0 flex-1">
                {row.current ? (
                  <Badge variant="secondary">{t('profile.sessions.current')}</Badge>
                ) : (
                  <span className="font-medium">{t('profile.sessions.otherDevice')}</span>
                )}
              </div>
              {row.current ? null : (
                <RevokeSessionButton
                  disabled={revokeSession.isPending}
                  onConfirm={() => revokeSession.mutate(row.id)}
                />
              )}
            </CardHeader>
            <CardContent className="px-4 py-3">
              <dl className="flex flex-col gap-2">
                {(
                  [
                    [t('profile.sessions.userAgent'), row.userAgent?.trim() || unknown],
                    [t('profile.sessions.networkAddress'), row.ipAddress?.trim() || unknown],
                    [t('profile.sessions.lastActivity'), formatSessionDate(row.lastUsedAt, locale)],
                    [t('profile.sessions.started'), formatSessionDate(row.createdAt, locale)],
                    [t('profile.sessions.expires'), formatSessionDate(row.expiresAt, locale)],
                  ] as const
                ).map(([label, value], index) => (
                  <div key={label}>
                    {index > 0 ? <Separator className="mb-2" /> : null}
                    <div className="flex items-start justify-between gap-4 text-sm">
                      <dt className="shrink-0 text-muted-foreground">{label}</dt>
                      <dd className="min-w-0 break-all text-right font-medium">{value}</dd>
                    </div>
                  </div>
                ))}
              </dl>
            </CardContent>
          </Card>
        )}
      />
    </div>
  )
}
