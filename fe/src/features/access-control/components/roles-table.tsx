'use client'

import { useState, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PencilIcon } from 'lucide-react'
import type { ColumnDef } from '@tanstack/react-table'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'
import { DataTable, DataTableColumnHeader, DataTableSkeleton } from '@/components/common/data-table'
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

import { AccessControlPageHeader } from './access-control-page-header'
import { ListRowActionsMenu } from './list-row-actions-menu'
import { MobileRecordCard } from './mobile-record-card'
import { PermissionGate } from './permission-gate'
import { RoleFormDialog } from './role-form-dialog'
import { PermissionKeys } from '../permission-keys'
import { usePermissions } from '../hooks/use-permissions'
import { useRoleMutations, useRolesList } from '../hooks/use-roles'
import type {
  CreateRoleFormValues,
  UpdateRoleFormValues,
} from '../schemas/access-control-schemas'
import type { Role } from '../types'

export function RolesTable() {
  const { t } = useTranslation('admin')
  const rolesQuery = useRolesList()
  const { canModify } = usePermissions()
  const { createRole, updateRole, deleteRole } = useRoleMutations()
  const canManageRoles = canModify('roles')

  const [dialogOpen, setDialogOpen] = useState(false)
  const [dialogMode, setDialogMode] = useState<'create' | 'edit'>('create')
  const [editingRole, setEditingRole] = useState<Role | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<Role | undefined>()

  const openCreate = () => {
    setDialogMode('create')
    setEditingRole(undefined)
    setDialogOpen(true)
  }

  const openEdit = (role: Role) => {
    setDialogMode('edit')
    setEditingRole(role)
    setDialogOpen(true)
  }

  const renderActions = (role: Role) =>
    canManageRoles ? (
    <ListRowActionsMenu>
      <PermissionGate permission={PermissionKeys.roles.modify}>
        <DropdownMenuItem onSelect={() => openEdit(role)}>
          {t('access.actions.edit')}
        </DropdownMenuItem>
      </PermissionGate>
      <PermissionGate permission={PermissionKeys.roles.modify}>
        <DropdownMenuItem
          variant="destructive"
          onSelect={() => setDeleteTarget(role)}
        >
          {t('access.actions.delete')}
        </DropdownMenuItem>
      </PermissionGate>
    </ListRowActionsMenu>
    ) : null

  const columns = useMemo<ColumnDef<Role>[]>(
    () => [
      {
        accessorKey: 'name',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.roles.fields.name')}
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
        accessorKey: 'slug',
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.roles.fields.slug')}
          />
        ),
        cell: ({ row }) => (
          <Badge variant="outline">{row.original.slug}</Badge>
        ),
      },
      {
        id: 'permissions',
        accessorFn: (row) => row.permissionKeys.length,
        header: ({ column }) => (
          <DataTableColumnHeader
            column={column}
            title={t('access.roles.fields.permissions')}
          />
        ),
        cell: ({ row }) => {
          const keys = row.original.permissionKeys
          if (keys.length === 0) {
            return (
              <span className="text-sm text-muted-foreground">
                {t('access.roles.noPermissions', { defaultValue: 'No permissions' })}
              </span>
            )
          }
          const groups = [...new Set(keys.map((k) => k.split(':')[0]))]
          return (
            <div className="flex flex-wrap items-center gap-1.5">
              {groups.map((g) => (
                <Badge
                  key={g}
                  variant="secondary"
                  className="h-5 px-1.5 text-[11px] capitalize"
                >
                  {g}
                </Badge>
              ))}
              <span className="text-xs tabular-nums text-muted-foreground">
                ({keys.length})
              </span>
            </div>
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
    [t],
  )

  if (rolesQuery.isLoading) {
    return (
      <div className="flex flex-col gap-4 p-4">
        <AccessControlPageHeader
          title={t('access.roles.title')}
          description={t('access.roles.description')}
        />
        <DataTableSkeleton columns={4} />
      </div>
    )
  }

  if (rolesQuery.isError) {
    return (
      <div className="flex flex-col gap-4 p-4">
        <AccessControlPageHeader
          title={t('access.roles.title')}
          description={t('access.roles.description')}
        />
        <Card>
          <CardContent className="flex flex-col gap-3 py-6">
            <p className="text-sm text-muted-foreground">
              {t('access.roles.error')}
            </p>
            <Button variant="outline" onClick={() => void rolesQuery.refetch()}>
              {t('access.actions.retry')}
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  const roles = rolesQuery.data ?? []

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <AccessControlPageHeader
        title={t('access.roles.title')}
        description={t('access.roles.description')}
        actions={
          <PermissionGate permission={PermissionKeys.roles.modify}>
            <Button className="w-full md:w-auto" onClick={openCreate}>
              {t('access.roles.createAction')}
            </Button>
          </PermissionGate>
        }
      />

      <DataTable
        columns={columns}
        data={roles}
        getRowId={(row) => row.id}
        filterKey="name"
        filterPlaceholder={t('access.roles.filterPlaceholder', {
          defaultValue: 'Search roles…',
        })}
        emptyMessage={t('access.roles.empty')}
        renderMobileCard={(role) => (
          <MobileRecordCard
            title={role.name}
            actions={renderActions(role)}
            fields={[
              {
                label: t('access.roles.fields.slug'),
                value: <Badge variant="outline">{role.slug}</Badge>,
              },
              {
                label: t('access.roles.fields.permissions'),
                value: role.permissionKeys.length,
              },
            ]}
          />
        )}
      />

      <RoleFormDialog
        open={dialogOpen}
        mode={dialogMode}
        role={editingRole}
        isPending={createRole.isPending || updateRole.isPending}
        onOpenChange={setDialogOpen}
        onCreate={(values: CreateRoleFormValues) => {
          createRole.mutate(values, {
            onSuccess: () => setDialogOpen(false),
          })
        }}
        onUpdate={(values: UpdateRoleFormValues) => {
          if (!editingRole) return
          updateRole.mutate(
            { id: editingRole.id, input: values },
            { onSuccess: () => setDialogOpen(false) },
          )
        }}
      />

      <AlertDialog
        open={Boolean(deleteTarget)}
        onOpenChange={(open) => !open && setDeleteTarget(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('access.roles.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('access.roles.deleteDescription', {
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
                deleteRole.mutate(deleteTarget.id, {
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
