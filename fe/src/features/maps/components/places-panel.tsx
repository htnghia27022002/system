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
import { StatusBadge } from '@/components/ui/status-badge'
import { MoreHorizontalIcon } from 'lucide-react'

import { useMapLocations } from '../hooks/use-map-locations'
import { useManagedPlaces, usePlaceMutations } from '../hooks/use-managed-places'
import type { PlacePin, PlaceStatus } from '../types'
import { PlaceFormDialog } from './place-form-dialog'

type PlacesPanelProps = {
  onRequestCreate?: () => void
}

function statusVariant(status: PlaceStatus) {
  if (status === 'active') return 'success' as const
  if (status === 'hidden') return 'neutral' as const
  return 'warning' as const
}

export function PlacesPanel({ onRequestCreate }: PlacesPanelProps) {
  const { t } = useTranslation('admin')
  const placesQuery = useManagedPlaces()
  const locationsQuery = useMapLocations()
  const { createPlace, updatePlace, deletePlace } = usePlaceMutations()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [mode, setMode] = useState<'create' | 'edit'>('create')
  const [editing, setEditing] = useState<PlacePin | undefined>()
  const [deleteTarget, setDeleteTarget] = useState<PlacePin | null>(null)

  const items = placesQuery.data?.items ?? []

  return (
    <div className="space-y-4 px-4 pb-6">
      {onRequestCreate ? (
        <PermissionGate permission={PermissionKeys.maps.modify}>
          <Button className="w-full" onClick={onRequestCreate}>
            {t('maps.places.createAction')}
          </Button>
        </PermissionGate>
      ) : null}

      {placesQuery.isError ? (
        <p className="text-sm text-muted-foreground">{t('maps.places.error')}</p>
      ) : null}
      {!placesQuery.isLoading && items.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('maps.places.empty')}</p>
      ) : null}

      <ul className="space-y-3">
        {items.map((place) => (
          <li key={place.id} className="rounded-lg border border-border p-3">
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0 space-y-1">
                <p className="truncate text-sm font-medium">{place.name}</p>
                <p className="truncate text-xs text-muted-foreground">
                  {place.locationName}
                </p>
                <StatusBadge
                  variant={statusVariant(place.status)}
                  label={t(`maps.status.${place.status}`)}
                />
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
                        setEditing(place)
                        setDialogOpen(true)
                      }}
                    >
                      {t('access.actions.edit')}
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      variant="destructive"
                      onSelect={() => setDeleteTarget(place)}
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

      <PlaceFormDialog
        open={dialogOpen}
        mode={mode}
        place={editing}
        locations={locationsQuery.data?.items ?? []}
        isPending={createPlace.isPending || updatePlace.isPending}
        onOpenChange={setDialogOpen}
        onSubmit={(input) => {
          if (mode === 'edit' && editing) {
            updatePlace.mutate(
              { id: editing.id, input },
              { onSuccess: () => setDialogOpen(false) },
            )
            return
          }
          createPlace.mutate(input, { onSuccess: () => setDialogOpen(false) })
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
            <AlertDialogTitle>{t('maps.places.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('maps.places.deleteDescription', {
                name: deleteTarget?.name ?? '',
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('access.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                if (!deleteTarget) return
                deletePlace.mutate(deleteTarget.id, {
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
