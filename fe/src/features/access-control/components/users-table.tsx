import { useState, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import type { ColumnDef } from '@tanstack/react-table'

import { PencilIcon } from 'lucide-react'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'
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
import { UserFormDialog } from './user-form-dialog'
import { usePermissions } from '../hooks/use-permissions'
import { useRolesList } from '../hooks/use-roles'
import { useUserMutations, useUsersList } from '../hooks/use-users'
import type {
  CreateUserFormValues,
  UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import type { ManagedUser } from '../types'

export function UsersTable() {
  const { t } = useTranslation('admin')
  const usersQuery = useUsersList({ page: 1, pageSize: 50 })
  const rolesQuery = useRolesList()
  const { hasAny } = usePermissions()
  const { createUser, updateUser, deleteUser } = useUserMutations()
  const canManageUsers = hasAny('users:update', 'users:delete')

  const [dialogOpen, setDialogOpen] = useState(false)
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create')
  const [editingUser, setEditingUser] = useState<ManagedUser | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<ManagedUser | undefined>()

  const roles = rolesQuery.data ?? []
  const roleById = useMemo(
    () => new Map(roles.map((r) => [r.id, r])),
    [roles],
  )

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
    createUser.mutate(values, {
      onSuccess: () => setDialogOpen(false),
    })
  }

  const handleUpdate = (values: UpdateUserFormValues) => {
    if (!editingUser) return
    const input = {
      name: values.name,
      email: values.email,
      roleId: values.roleId,
      status: values.status,
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
        <PermissionGate permission="users:update">
          <DropdownMenuItem onSelect={() => openEdit(user)}>
            {t('access.actions.edit')}
          </DropdownMenuItem>
        </PermissionGate>
        <PermissionGate permission="users:delete">
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
        id: 'actions',
        size: 48,
        enableSorting: false,
        cell: ({ row }) => renderActions(row.original),
      },
    ],
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [t, roleById, canManageUsers],
  )

  if (usersQuery.isLoading || rolesQuery.isLoading) {
    return (
      <div className="flex flex-col gap-4 p-4">
        <AccessControlPageHeader
          title={t('access.users.title')}
          description={t('access.users.description')}
        />
        <DataTableSkeleton columns={5} />
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

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <AccessControlPageHeader
        title={t('access.users.title')}
        description={t('access.users.description')}
        actions={
          <PermissionGate permission="users:create">
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
        filterKey="name"
        filterPlaceholder={t('access.users.filterPlaceholder', {
          defaultValue: 'Search users…',
        })}
        emptyMessage={t('access.users.empty')}
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
