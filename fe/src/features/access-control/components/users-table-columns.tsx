'use client'

import type { ColumnDef } from '@tanstack/react-table'
import type { TFunction } from 'i18next'
import { ShieldCheckIcon } from 'lucide-react'
import type { ReactNode } from 'react'

import { DataTableColumnHeader } from '@/components/common/data-table'
import { Badge } from '@/components/ui/badge'
import { StatusBadge } from '@/components/ui/status-badge'

import { UserIdentity } from './user-identity'
import { USER_STATUS_MAP } from '../enum-maps'
import type { Role, ManagedUser } from '../types'

export function formatUserDate(iso: string | undefined, locale: string, withTime = false) {
  if (!iso) return '—'
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString(locale, withTime
    ? { dateStyle: 'medium', timeStyle: 'short' }
    : { dateStyle: 'medium' })
}

export function UserRoleBadges({
  user,
  role,
  t,
}: {
  user: ManagedUser
  role?: Role
  t: TFunction<'admin'>
}) {
  return (
    <span className="inline-flex flex-wrap items-center gap-1">
      {role ? (
        <Badge variant="secondary">{role.name}</Badge>
      ) : (
        <span className="text-muted-foreground">—</span>
      )}
      {user.superAdmin ? (
        <Badge className="gap-1">
          <ShieldCheckIcon aria-hidden />
          {t('access.users.badge.superAdmin')}
        </Badge>
      ) : null}
    </span>
  )
}

export function UserStatusBadge({ user, t }: { user: ManagedUser; t: TFunction<'admin'> }) {
  const cfg = USER_STATUS_MAP[user.status]
  return (
    <StatusBadge
      variant={cfg.variant}
      pulse={cfg.pulse}
      label={t(`access.users.status.${user.status}`)}
    />
  )
}

type BuildUserColumnsOptions = {
  t: TFunction<'admin'>
  locale: string
  roleById: Map<string, Role>
  currentUserId?: string
  onOpen?: (user: ManagedUser) => void
  renderActions: (user: ManagedUser) => ReactNode
}

export function buildUserColumns({
  t,
  locale,
  roleById,
  currentUserId,
  onOpen,
  renderActions,
}: BuildUserColumnsOptions): ColumnDef<ManagedUser>[] {
  return [
    {
      accessorKey: 'name',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('access.users.fields.user')} />
      ),
      cell: ({ row }) => (
        <UserIdentity
          user={row.original}
          isCurrentUser={row.original.id === currentUserId}
          onOpen={onOpen ? () => onOpen(row.original) : undefined}
        />
      ),
    },
    {
      id: 'role',
      accessorFn: (row) => roleById.get(row.roleId)?.name ?? '',
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('access.users.fields.role')} />
      ),
      cell: ({ row }) => (
        <UserRoleBadges user={row.original} role={roleById.get(row.original.roleId)} t={t} />
      ),
    },
    {
      accessorKey: 'status',
      size: 128,
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('access.users.fields.status')} />
      ),
      cell: ({ row }) => <UserStatusBadge user={row.original} t={t} />,
    },
    {
      accessorKey: 'createdAt',
      size: 160,
      header: ({ column }) => (
        <DataTableColumnHeader column={column} title={t('access.users.fields.createdAt')} />
      ),
      cell: ({ row }) => (
        <time
          dateTime={row.original.createdAt}
          title={formatUserDate(row.original.createdAt, locale, true)}
          className="whitespace-nowrap text-muted-foreground tabular-nums"
        >
          {formatUserDate(row.original.createdAt, locale)}
        </time>
      ),
    },
    {
      id: 'actions',
      size: 48,
      enableSorting: false,
      header: () => <span className="sr-only">{t('access.users.fields.actions')}</span>,
      cell: ({ row }) => <div className="flex justify-end">{renderActions(row.original)}</div>,
    },
  ]
}
