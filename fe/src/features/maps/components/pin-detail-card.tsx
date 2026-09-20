'use client'

import { useState } from 'react'
import {
  ClockIcon,
  GlobeIcon,
  MapPinIcon,
  PhoneIcon,
  XIcon,
} from 'lucide-react'
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
import { Skeleton } from '@/components/ui/skeleton'
import { StatusBadge } from '@/components/ui/status-badge'

import { useMapLocationMutations } from '../hooks/use-map-locations'
import { usePlaceMutations } from '../hooks/use-managed-places'
import { usePlaceDetail } from '../hooks/use-place-detail'
import type { LocationRecord, PlacePin, PlaceStatus } from '../types'
import { LocationFormDialog } from './location-form-dialog'
import { PlaceQuickCreate } from './place-quick-create'
import { PinPhoto } from './pin-photo'

type PinDetailCardProps = {
  placeId?: string | null
  location?: LocationRecord | null
  places?: PlacePin[]
  onClose: () => void
  onPlaceCreated?: (placeId: string) => void
  onSelectPlace?: (placeId: string) => void
}

function statusBadgeVariant(status: PlaceStatus) {
  if (status === 'active') return 'success' as const
  if (status === 'hidden') return 'neutral' as const
  return 'warning' as const
}

export function PinDetailCard({
  placeId,
  location,
  places = [],
  onClose,
  onPlaceCreated,
  onSelectPlace,
}: PinDetailCardProps) {
  if (placeId) {
    return (
      <PlaceDetailCard
        placeId={placeId}
        places={places}
        onClose={onClose}
        onPlaceCreated={onPlaceCreated}
        onSelectPlace={onSelectPlace}
      />
    )
  }
  if (location) {
    return (
      <LocationDetailCard
        location={location}
        places={places}
        onClose={onClose}
        onPlaceCreated={onPlaceCreated}
        onSelectPlace={onSelectPlace}
      />
    )
  }
  return null
}

function PlaceDetailCard({
  placeId,
  places,
  onClose,
  onPlaceCreated,
  onSelectPlace,
}: {
  placeId: string
  places: PlacePin[]
  onClose: () => void
  onPlaceCreated?: (placeId: string) => void
  onSelectPlace?: (placeId: string) => void
}) {
  const { t } = useTranslation('admin')
  const detailQuery = usePlaceDetail(placeId)
  const { createPlace } = usePlaceMutations()
  const [addingPlace, setAddingPlace] = useState(false)
  const place = detailQuery.data
  const siblings = places.filter(
    (item) => item.locationId === place?.location.id,
  )
  const position =
    place?.lat != null && place.lng != null
      ? { lat: place.lat, lng: place.lng }
      : place?.location.lat != null && place.location.lng != null
        ? { lat: place.location.lat, lng: place.location.lng }
        : null

  return (
    <article className="pointer-events-auto absolute inset-x-3 bottom-4 z-[1200] mx-auto flex max-h-[min(36rem,calc(100%-5.5rem))] w-full max-w-md flex-col overflow-hidden rounded-2xl border border-border bg-card shadow-xl">
      <CardClose onClose={onClose} label={t('access.actions.cancel')} />
      {detailQuery.isLoading ? (
        <div className="space-y-3 p-4">
          <Skeleton className="h-36 w-full rounded-xl" />
          <Skeleton className="h-6 w-2/3" />
          <Skeleton className="h-4 w-1/2" />
        </div>
      ) : null}
      {detailQuery.isError ? (
        <div className="space-y-3 p-4">
          <p className="text-sm text-muted-foreground">{t('maps.detail.error')}</p>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void detailQuery.refetch()}
          >
            {t('maps.retry')}
          </Button>
        </div>
      ) : null}
      {place ? (
        <div className="min-h-0 flex-1 overflow-y-auto">
          {position ? (
            <div className="relative h-40 bg-muted">
              <PinPhoto
                position={position}
                alt={place.name}
                className="size-full object-cover"
              />
              <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-card to-transparent px-4 pb-3 pt-10">
                <h2 className="text-lg font-semibold leading-tight text-foreground">
                  {place.name}
                </h2>
              </div>
            </div>
          ) : (
            <div className="px-4 pt-10">
              <h2 className="text-lg font-semibold">{place.name}</h2>
            </div>
          )}
          <div className="space-y-4 px-4 pb-4 pt-2">
            <div className="flex flex-wrap items-center gap-2">
              <p className="text-sm text-muted-foreground">
                {t(`maps.categories.${place.category}`, {
                  defaultValue: place.category,
                })}
              </p>
              <StatusBadge
                variant={statusBadgeVariant(place.status)}
                label={t(`maps.status.${place.status}`)}
              />
            </div>

            <p className="flex items-start gap-2 text-sm text-foreground">
              <MapPinIcon className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
              <span>
                {place.location.name}
                {place.location.formatted
                  ? ` · ${place.location.formatted}`
                  : ''}
                {place.unit ? ` · ${place.unit}` : ''}
              </span>
            </p>

            {place.details?.description ? (
              <p className="text-sm text-muted-foreground">
                {place.details.description}
              </p>
            ) : null}

            {place.details?.phone ? (
              <p className="flex items-center gap-2 text-sm">
                <PhoneIcon className="size-4 text-muted-foreground" />
                {place.details.phone}
              </p>
            ) : null}
            {place.details?.hours ? (
              <p className="flex items-center gap-2 text-sm">
                <ClockIcon className="size-4 text-muted-foreground" />
                {place.details.hours}
              </p>
            ) : null}
            {place.details?.website ? (
              <p className="flex items-center gap-2 text-sm">
                <GlobeIcon className="size-4 text-muted-foreground" />
                <a
                  href={place.details.website}
                  target="_blank"
                  rel="noreferrer"
                  className="text-primary underline-offset-4 hover:underline"
                >
                  {t('maps.detail.website')}
                </a>
              </p>
            ) : null}
            {place.details?.priceRange ? (
              <p className="text-sm text-muted-foreground">
                {t('maps.detail.priceRange')}: {place.details.priceRange}
              </p>
            ) : null}

            <PlacesAtLocation
              places={siblings}
              selectedId={place.id}
              onSelect={onSelectPlace}
            />

            <PermissionGate permission={PermissionKeys.maps.modify}>
              {addingPlace ? (
                <PlaceQuickCreate
                  locationId={place.location.id}
                  locationName={place.location.name}
                  lat={place.lat ?? place.location.lat}
                  lng={place.lng ?? place.location.lng}
                  isPending={createPlace.isPending}
                  onCancel={() => setAddingPlace(false)}
                  onSubmit={(input) =>
                    createPlace.mutate(input, {
                      onSuccess: (created) => {
                        setAddingPlace(false)
                        onPlaceCreated?.(created.id)
                      },
                    })
                  }
                />
              ) : (
                <Button
                  type="button"
                  size="sm"
                  onClick={() => setAddingPlace(true)}
                >
                  {t('maps.places.addHere')}
                </Button>
              )}
            </PermissionGate>

            {place.news.length > 0 ? (
              <section className="space-y-2">
                <h3 className="text-sm font-medium">{t('maps.detail.news')}</h3>
                <ul className="space-y-2">
                  {place.news.map((item) => (
                    <li
                      key={item.id}
                      className="rounded-lg border border-border px-3 py-2"
                    >
                      <p className="text-sm font-medium">{item.title}</p>
                      <p className="text-xs text-muted-foreground">
                        {item.sourceName}
                      </p>
                      {item.originalUrl ? (
                        <a
                          href={item.originalUrl}
                          target="_blank"
                          rel="noreferrer"
                          className="text-sm text-primary underline-offset-4 hover:underline"
                        >
                          {t('maps.detail.originalUrl')}
                        </a>
                      ) : null}
                    </li>
                  ))}
                </ul>
              </section>
            ) : null}

          </div>
        </div>
      ) : null}
    </article>
  )
}

function LocationDetailCard({
  location,
  places,
  onClose,
  onPlaceCreated,
  onSelectPlace,
}: {
  location: LocationRecord
  places: PlacePin[]
  onClose: () => void
  onPlaceCreated?: (placeId: string) => void
  onSelectPlace?: (placeId: string) => void
}) {
  const { t } = useTranslation('admin')
  const { updateLocation, deleteLocation } = useMapLocationMutations()
  const { createPlace } = usePlaceMutations()
  const [editOpen, setEditOpen] = useState(false)
  const [addingPlace, setAddingPlace] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState(false)
  const placesHere = places.filter((item) => item.locationId === location.id)
  const position =
    location.lat != null && location.lng != null
      ? { lat: location.lat, lng: location.lng }
      : null

  return (
    <article className="pointer-events-auto absolute inset-x-3 bottom-4 z-[1200] mx-auto flex max-h-[min(36rem,calc(100%-5.5rem))] w-full max-w-md flex-col overflow-hidden rounded-2xl border border-border bg-card shadow-xl">
      <CardClose onClose={onClose} label={t('access.actions.cancel')} />
      <div className="min-h-0 flex-1 overflow-y-auto">
        {position ? (
          <div className="relative h-40 bg-muted">
            <PinPhoto
              position={position}
              alt={location.name}
              className="size-full object-cover"
            />
            <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-card to-transparent px-4 pb-3 pt-10">
              <h2 className="text-lg font-semibold leading-tight">
                {location.name}
              </h2>
            </div>
          </div>
        ) : (
          <div className="px-4 pt-10">
            <h2 className="text-lg font-semibold">{location.name}</h2>
          </div>
        )}
        <div className="space-y-3 px-4 pb-4 pt-2">
          <p className="text-sm text-muted-foreground">
            {t('maps.locations.title')}
            {location.countryCode ? ` · ${location.countryCode}` : ''}
          </p>
          <p className="text-sm text-muted-foreground">
            {t('maps.locations.addPlaceHint')}
          </p>
          {location.formatted || location.street ? (
            <p className="flex items-start gap-2 text-sm">
              <MapPinIcon className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
              <span>{location.formatted || location.street}</span>
            </p>
          ) : null}
          {position ? (
            <p className="text-xs text-muted-foreground">
              {position.lat.toFixed(5)}, {position.lng.toFixed(5)}
            </p>
          ) : null}
          <PlacesAtLocation
            places={placesHere}
            onSelect={onSelectPlace}
          />

          <PermissionGate permission={PermissionKeys.maps.modify}>
            {addingPlace ? (
              <PlaceQuickCreate
                locationId={location.id}
                locationName={location.name}
                lat={location.lat}
                lng={location.lng}
                isPending={createPlace.isPending}
                onCancel={() => setAddingPlace(false)}
                onSubmit={(input) =>
                  createPlace.mutate(input, {
                    onSuccess: (created) => {
                      setAddingPlace(false)
                      onPlaceCreated?.(created.id)
                    },
                  })
                }
              />
            ) : null}
            <div className="flex flex-wrap gap-2">
              {!addingPlace ? (
                <Button
                  type="button"
                  size="sm"
                  onClick={() => setAddingPlace(true)}
                >
                  {t('maps.places.addHere')}
                </Button>
              ) : null}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setEditOpen(true)}
              >
                {t('access.actions.edit')}
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setConfirmDelete(true)}
              >
                {t('access.actions.delete')}
              </Button>
            </div>
          </PermissionGate>
        </div>
      </div>

      <LocationFormDialog
        open={editOpen}
        mode="edit"
        location={location}
        isPending={updateLocation.isPending}
        onOpenChange={setEditOpen}
        onSubmit={(input) =>
          updateLocation.mutate(
            { id: location.id, input },
            { onSuccess: () => setEditOpen(false) },
          )
        }
      />
      <AlertDialog open={confirmDelete} onOpenChange={setConfirmDelete}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('maps.locations.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('maps.locations.deleteDescription', { name: location.name })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('access.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction
              onClick={() =>
                deleteLocation.mutate(location.id, { onSuccess: onClose })
              }
            >
              {t('access.actions.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </article>
  )
}

function PlacesAtLocation({
  places,
  selectedId,
  onSelect,
}: {
  places: PlacePin[]
  selectedId?: string
  onSelect?: (placeId: string) => void
}) {
  const { t } = useTranslation('admin')
  if (places.length === 0) return null
  if (selectedId && places.length < 2) return null
  return (
    <section className="space-y-2">
      <h3 className="text-sm font-medium">{t('maps.places.alsoHere')}</h3>
      <ul className="flex flex-wrap gap-2">
        {places.map((item) => {
          const selected = item.id === selectedId
          return (
            <li key={item.id}>
              <Button
                type="button"
                size="sm"
                variant={selected ? 'default' : 'outline'}
                disabled={selected || !onSelect}
                onClick={() => onSelect?.(item.id)}
              >
                {item.name}
              </Button>
            </li>
          )
        })}
      </ul>
    </section>
  )
}

function CardClose({ onClose, label }: { onClose: () => void; label: string }) {
  return (
    <Button
      type="button"
      size="icon-sm"
      variant="secondary"
      className="absolute top-3 right-3 z-10 rounded-full"
      onClick={onClose}
      aria-label={label}
    >
      <XIcon className="size-4" />
    </Button>
  )
}
