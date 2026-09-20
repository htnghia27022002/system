'use client'

import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { PermissionGate, PermissionKeys } from '@/features/access-control'
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
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreHorizontalIcon } from 'lucide-react'

import {
  useMapLocationMutations,
  useMapLocations,
} from '../hooks/use-map-locations'
import type { LocationRecord } from '../types'
import { LocationFormDialog } from './location-form-dialog'

type LocationsPanelProps = {
  onRequestCreate?: () => void
}

export function LocationsPanel({ onRequestCreate }: LocationsPanelProps) {
  const { t } = useTranslation('admin')
  const locationsQuery = useMapLocations()
  const { createLocation, updateLocation, deleteLocation } =
    useMapLocationMutations()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [mode, setMode] = useState<'create' | 'edit'>('create')
  const [editing, setEditing] = useState<LocationRecord | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<LocationRecord | null>(null)

  const items = locationsQuery.data?.items ?? []

  return (
    <div className="space-y-4 px-4 pb-6">
      {onRequestCreate ? (
        <PermissionGate permission={PermissionKeys.maps.modify}>
          <Button className="w-full" onClick={onRequestCreate}>
            {t('maps.locations.createAction')}
          </Button>
        </PermissionGate>
      ) : null}

      {locationsQuery.isError ? (
        <p className="text-sm text-muted-foreground">{t('maps.locations.error')}</p>
      ) : null}
      {!locationsQuery.isLoading && items.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('maps.locations.empty')}</p>
      ) : null}

      <ul className="space-y-3">
        {items.map((location) => (
          <li key={location.id} className="rounded-lg border border-border p-3">
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0">
                <p className="truncate text-sm font-medium">{location.name}</p>
                <p className="truncate text-xs text-muted-foreground">
                  {location.formatted || location.street || '—'}
                </p>
              </div>
              <PermissionGate permission={PermissionKeys.maps.modify}>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon-sm">
                      <MoreHorizontalIcon className="size-4" />
                      <span className="sr-only">{t('access.actions.openMenu')}</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onSelect={() => {
                        setMode('edit')
                        setEditing(location)
                        setDialogOpen(true)
                      }}
                    >
                      {t('access.actions.edit')}
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      variant="destructive"
                      onSelect={() => setDeleteTarget(location)}
                    >
                      {t('access.actions.delete')}
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </PermissionGate>
            </div>
          </li>
        ))}
      </ul>

      <LocationFormDialog
        open={dialogOpen}
        mode={mode}
        location={editing}
        isPending={createLocation.isPending || updateLocation.isPending}
        onOpenChange={setDialogOpen}
        onSubmit={(input) => {
          if (mode === 'edit' && editing) {
            updateLocation.mutate(
              { id: editing.id, input },
              { onSuccess: () => setDialogOpen(false) },
            )
            return
          }
          createLocation.mutate(input, { onSuccess: () => setDialogOpen(false) })
        }}
      />

      <AlertDialog
        open={Boolean(deleteTarget)}
        onOpenChange={(open) => {
          if (!open) setDeleteTarget(null)
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('maps.locations.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('maps.locations.deleteDescription', {
                name: deleteTarget?.name ?? '',
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('access.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (!deleteTarget) return
                deleteLocation.mutate(deleteTarget.id, {
                  onSuccess: () => setDeleteTarget(null),
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
