'use client'

import { Suspense, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import type { ColumnDef } from '@tanstack/react-table'

import { PencilIcon } from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { StatusBadge } from '@/components/ui/status-badge'
import { DataTable, DataTableColumnHeader, DataTableSkeleton } from '@/components/common/data-table'
import { USER_STATUS_MAP } from '../enum-maps'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Card, CardContent } from '@/components/ui/card'

import { AccessControlPageHeader } from './access-control-page-header'
import { ListRowActionsMenu } from './list-row-actions-menu'
import { MobileRecordCard } from './mobile-record-card'
import { PermissionGate } from './permission-gate'
import { ServerListSearch } from './server-list-search'
import { UserFormDialog } from './user-form-dialog'
import { UserOAuthProviders } from './user-oauth-providers'
import { PermissionKeys } from '../permission-keys'
import { useListQueryParams } from '../hooks/use-list-query-params'
import { usePermissions } from '../hooks/use-permissions'
import { useAllRolesList } from '../hooks/use-roles'
import { useUserMutations, useUsersList } from '../hooks/use-users'
import type {
  CreateUserFormValues,
  UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import type { ManagedUser } from '../types'

export function UsersTable() {
  return (
    <Suspense
      fallback={
        <div className="flex flex-col gap-4 p-4">
          <DataTableSkeleton columns={8} />
        </div>
      }
    >
      <UsersTableContent />
    </Suspense>
  )
}

function UsersTableContent() {
  const { t, i18n } = useTranslation('admin')
  const { values, setParams } = useListQueryParams([
    'search',
    'roleId',
    'status',
    'id',
  ] as const)
  const page = Number(values.page) || 1
  const pageSize = 50

  const usersQuery = useUsersList({
    page,
    pageSize,
    search: values.search || undefined,
    roleId: values.roleId || undefined,
    status: (values.status as ManagedUser['status']) || undefined,
    id: values.id || undefined,
  })
  const rolesQuery = useAllRolesList()
  const { canModify } = usePermissions()
  const { createUser, updateUser, deleteUser } = useUserMutations()
  const canManageUsers = canModify('users')

  const [dialogOpen, setDialogOpen] = useState(false)
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create')
  const [editingUser, setEditingUser] = useState<ManagedUser | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<ManagedUser | undefined>()

  const roles = rolesQuery.data ?? []
  const roleById = useMemo(
    () => new Map(roles.map((r) => [r.id, r])),
    [roles],
  )

  const formatCreatedAt = (iso?: string) => {
    if (!iso) return '—'
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return iso
    return d.toLocaleString(i18n.language, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const openCreate = () => {
    setDialogMode('create')
    setEditingUser(undefined)
    setDialogOpen(true)
  }

  const openEdit = (user: ManagedUser) => {
    setDialogMode('edit')
    setEditingUser(user)
    setDialogOpen(true)
  }

  const handleCreate = (values: CreateUserFormValues) => {
    createUser.mutate(
      {
        name: values.name,
        email: values.email,
        password: values.password,
        roleId: values.roleId,
        status: values.status,
        phone: values.phone,
        general: values.general,
        birthday: values.birthday?.trim() ? values.birthday : null,
        address: values.address,
        socialLinks: values.socialLinks,
      },
      {
        onSuccess: () => setDialogOpen(false),
      },
    )
  }

  const handleUpdate = (values: UpdateUserFormValues) => {
    if (!editingUser) return
    const input = {
      name: values.name,
      email: values.email,
      roleId: values.roleId,
      status: values.status,
      phone: values.phone,
      general: values.general,
      birthday: values.birthday?.trim() ? values.birthday : null,
      address: values.address,
      socialLinks: values.socialLinks,
      ...(values.password ? { password: values.password } : {}),
    }
    updateUser.mutate(
      { id: editingUser.id, input },
      { onSuccess: () => setDialogOpen(false) },
    )
  }

  const renderActions = (user: ManagedUser) =>
    canManageUsers ? (
      <ListRowActionsMenu>
        <PermissionGate permission={PermissionKeys.users.modify}>
          <DropdownMenuItem onSelect={() => openEdit(user)}>
            {t('access.actions.edit')}
          </DropdownMenuItem>
        </PermissionGate>
        <PermissionGate permission={PermissionKeys.users.modify}>
          <DropdownMenuItem
            variant="destructive"
            onSelect={() => setDeleteTarget(user)}
          >
            {t('access.actions.delete')}
          </DropdownMenuItem>
        </PermissionGate>
      </ListRowActionsMenu>
    ) : null

  const columns = useMemo<ColumnDef<ManagedUser>[]>(
    () => [
      {
        id: 'index',
        size: 56,
        enableSorting: false,
        header: () => (
          <span className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
            {t('access.users.fields.index')}
          </span>
        ),
        cell: ({ row }) => (
          <span className="tabular-nums text-muted-foreground">
            {(page - 1) * pageSize + row.index + 1}
          </span>
        ),
      },
      {
        accessorKey: 'name',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.name')}
          />
        ),
        cell: ({ row }) => (
          <button
            type="button"
            className="group flex items-center gap-1.5 text-left font-medium hover:text-primary focus-visible:outline-none focus-visible:underline"
            onClick={() => openEdit(row.original)}
          >
            <span className="hover:underline">{row.original.name}</span>
            <PencilIcon
              className="size-3 shrink-0 text-muted-foreground/0 transition-colors group-hover:text-muted-foreground"
              aria-hidden
            />
          </button>
        ),
      },
      {
        accessorKey: 'email',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.email')}
          />
        ),
      },
      {
        id: 'oauth',
        accessorFn: (row) => (row.oauthProviders ?? []).join(','),
        enableSorting: false,
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.oauth')}
          />
        ),
        cell: ({ row }) => (
          <UserOAuthProviders providers={row.original.oauthProviders} />
        ),
      },
      {
        id: 'role',
        accessorFn: (row) => roleById.get(row.roleId)?.name ?? '',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.role')}
          />
        ),
        cell: ({ getValue }) => {
          const name = getValue<string>()
          return name ? (
            <Badge variant="secondary">{name}</Badge>
          ) : (
            <span className="text-muted-foreground">—</span>
          )
        },
      },
      {
        accessorKey: 'status',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.status')}
          />
        ),
        cell: ({ row }) => {
          const status = row.original.status
          const cfg = USER_STATUS_MAP[status]
          return (
            <StatusBadge
              variant={cfg.variant}
              pulse={cfg.pulse}
              label={t(`access.users.status.${status}`)}
            />
          )
        },
      },
      {
        accessorKey: 'createdAt',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.users.fields.createdAt')}
          />
        ),
        cell: ({ row }) => (
          <span className="text-muted-foreground tabular-nums">
            {formatCreatedAt(row.original.createdAt)}
          </span>
        ),
      },
      {
        id: 'actions',
        size: 48,
        enableSorting: false,
        cell: ({ row }) => renderActions(row.original),
      },
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [t, i18n.language, roleById, canManageUsers, page, pageSize],
  )

  if (usersQuery.isLoading || rolesQuery.isLoading) {
    return (
      <div className="flex flex-col gap-4 p-4">
        <AccessControlPageHeader
          title={t('access.users.title')}
          description={t('access.users.description')}
        />
        <DataTableSkeleton columns={8} />
      </div>
    )
  }

  if (usersQuery.isError) {
    return (
      <div className="flex flex-col gap-4 p-4">
        <AccessControlPageHeader
          title={t('access.users.title')}
          description={t('access.users.description')}
        />
        <Card>
          <CardContent className="flex flex-col gap-3 py-6">
            <p className="text-sm text-muted-foreground">
              {t('access.users.error')}
            </p>
            <Button variant="outline" onClick={() => void usersQuery.refetch()}>
              {t('access.actions.retry')}
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  const users = usersQuery.data?.items ?? []
  const isRefreshing = usersQuery.isFetching && !usersQuery.isLoading

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <AccessControlPageHeader
        title={t('access.users.title')}
        description={t('access.users.description')}
        actions={
          <PermissionGate permission={PermissionKeys.users.modify}>
            <Button className="w-full md:w-auto" onClick={openCreate}>
              {t('access.users.createAction')}
            </Button>
          </PermissionGate>
        }
      />

      <DataTable
        columns={columns}
        data={users}
        getRowId={(row) => row.id}
        localPagination={false}
        isRefreshing={isRefreshing}
        emptyTitle={t('access.users.emptyTitle', { defaultValue: 'No users yet' })}
        emptyDescription={t('access.users.emptyDescription', {
          defaultValue: 'Create a user account to get started.',
        })}
        toolbar={
          <div className="flex w-full flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
            <ServerListSearch
              value={values.search}
              placeholder={t('access.users.filterPlaceholder', {
                defaultValue: 'Search users…',
              })}
              onSearch={(next) => setParams({ search: next || undefined })}
              onClear={() => setParams({ search: undefined })}
            />
            <Select
              value={values.roleId || 'all'}
              onValueChange={(value) =>
                setParams({ roleId: value === 'all' ? undefined : value })
              }
            >
              <SelectTrigger
                className="w-full sm:w-[11rem]"
                aria-label={t('access.users.filters.role')}
              >
                <SelectValue placeholder={t('access.users.filters.role')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">
                  {t('access.users.filters.allRoles')}
                </SelectItem>
                {roles.map((role) => (
                  <SelectItem key={role.id} value={role.id}>
                    {role.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select
              value={values.status || 'all'}
              onValueChange={(value) =>
                setParams({ status: value === 'all' ? undefined : value })
              }
            >
              <SelectTrigger
                className="w-full sm:w-[10rem]"
                aria-label={t('access.users.filters.status')}
              >
                <SelectValue placeholder={t('access.users.filters.status')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">
                  {t('access.users.filters.allStatuses')}
                </SelectItem>
                <SelectItem value="active">
                  {t('access.users.status.active')}
                </SelectItem>
                <SelectItem value="inactive">
                  {t('access.users.status.inactive')}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        }
        renderMobileCard={(user) => {
          const role = roleById.get(user.roleId)
          return (
            <MobileRecordCard
              title={user.name}
              actions={renderActions(user)}
              fields={[
                {
                  label: t('access.users.fields.email'),
                  value: (
                    <span className="break-all text-right">{user.email}</span>
                  ),
                },
                {
                  label: t('access.users.fields.oauth'),
                  value: <UserOAuthProviders providers={user.oauthProviders} />,
                },
                {
                  label: t('access.users.fields.role'),
                  value: role ? (
                    <Badge variant="secondary">{role.name}</Badge>
                  ) : (
                    '—'
                  ),
                },
                {
                  label: t('access.users.fields.status'),
                  value: (() => {
                    const cfg = USER_STATUS_MAP[user.status]
                    return (
                      <StatusBadge
                        variant={cfg.variant}
                        pulse={cfg.pulse}
                        label={t(`access.users.status.${user.status}`)}
                      />
                    )
                  })(),
                },
                {
                  label: t('access.users.fields.createdAt'),
                  value: formatCreatedAt(user.createdAt),
                },
              ]}
            />
          )
        }}
      />

      <UserFormDialog
        open={dialogOpen}
        mode={dialogMode}
        user={editingUser}
        roles={roles}
        isPending={createUser.isPending || updateUser.isPending}
        onOpenChange={setDialogOpen}
        onCreate={handleCreate}
        onUpdate={handleUpdate}
      />

      <AlertDialog
        open={Boolean(deleteTarget)}
        onOpenChange={(open) => !open && setDeleteTarget(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('access.users.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('access.users.deleteDescription', {
                name: deleteTarget?.name,
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('access.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={() => {
                if (!deleteTarget) return
                deleteUser.mutate(deleteTarget.id, {
                  onSuccess: () => setDeleteTarget(undefined),
                })
              }}
            >
              {t('access.actions.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
