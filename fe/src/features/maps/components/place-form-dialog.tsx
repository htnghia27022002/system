'use client'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'

import type {
  LocationRecord,
  PlaceCategoryKey,
  PlaceDetail,
  PlacePin,
  PlaceStatus,
  PlaceWriteInput,
} from '../types'

const CATEGORIES: PlaceCategoryKey[] = [
  'roomRental',
  'restaurant',
  'hotel',
  'eatery',
  'uncategorized',
]

type PlaceFormDialogProps = {
  open: boolean
  mode: 'create' | 'edit'
  place?: PlacePin | PlaceDetail
  locations: LocationRecord[]
  fixedLocationId?: string
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (input: PlaceWriteInput) => void
}

export function PlaceFormDialog({
  open,
  mode,
  place,
  locations,
  fixedLocationId,
  isPending,
  onOpenChange,
  onSubmit,
}: PlaceFormDialogProps) {
  const { t } = useTranslation('admin')
  const [locationId, setLocationId] = useState('')
  const [name, setName] = useState('')
  const [category, setCategory] = useState('eatery')
  const [unit, setUnit] = useState('')
  const [placeKey, setPlaceKey] = useState('')
  const [status, setStatus] = useState<PlaceStatus>('pending')
  const [lat, setLat] = useState('')
  const [lng, setLng] = useState('')

  useEffect(() => {
    if (!open) return
    const detail = place && 'location' in place ? place : undefined
    const nextLocationId =
      fixedLocationId ||
      detail?.location.id ||
      (place && 'locationId' in place ? place.locationId : '') ||
      ''
    const host = locations.find((location) => location.id === nextLocationId)
    setLocationId(nextLocationId)
    setName(place?.name ?? '')
    setCategory(place?.category ?? 'eatery')
    setUnit(place && 'unit' in place ? (place.unit ?? '') : '')
    setPlaceKey(place && 'placeKey' in place ? (place.placeKey ?? '') : '')
    setStatus(
      place?.status ??
        (host?.lat != null && host.lng != null ? 'active' : 'pending'),
    )
    setLat(
      place?.lat != null
        ? String(place.lat)
        : host?.lat != null
          ? String(host.lat)
          : '',
    )
    setLng(
      place?.lng != null
        ? String(place.lng)
        : host?.lng != null
          ? String(host.lng)
          : '',
    )
  }, [fixedLocationId, locations, open, place])

  const parseCoord = (value: string) => {
    const trimmed = value.trim()
    if (!trimmed) return undefined
    const n = Number(trimmed)
    return Number.isFinite(n) ? n : undefined
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>
            {mode === 'edit' ? t('maps.places.editTitle') : t('maps.places.createTitle')}
          </DialogTitle>
        </DialogHeader>
        <form
          className="grid gap-4 sm:grid-cols-2"
          onSubmit={(event) => {
            event.preventDefault()
            onSubmit({
              locationId,
              name: name.trim(),
              category,
              unit: unit.trim(),
              placeKey: placeKey.trim() || undefined,
              status,
              lat: parseCoord(lat),
              lng: parseCoord(lng),
            })
          }}
        >
          {fixedLocationId ? (
            <p className="text-sm text-muted-foreground sm:col-span-2">
              {t('maps.create.parentLocation', {
                name:
                  locations.find((location) => location.id === fixedLocationId)
                    ?.name ?? '',
              })}
            </p>
          ) : (
            <Field className="sm:col-span-2">
              <FieldLabel>{t('maps.places.fields.location')}</FieldLabel>
              <Select value={locationId} onValueChange={setLocationId}>
                <SelectTrigger aria-label={t('maps.places.fields.location')}>
                  <SelectValue placeholder={t('maps.places.fields.location')} />
                </SelectTrigger>
                <SelectContent>
                  {locations.map((location) => (
                    <SelectItem key={location.id} value={location.id}>
                      {location.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>
          )}
          <Field>
            <FieldLabel htmlFor="place-name">{t('maps.places.fields.name')}</FieldLabel>
            <Input
              id="place-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              required
            />
          </Field>
          <Field>
            <FieldLabel>{t('maps.filters.category')}</FieldLabel>
            <Select value={category} onValueChange={setCategory}>
              <SelectTrigger aria-label={t('maps.filters.category')}>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {CATEGORIES.map((key) => (
                  <SelectItem key={key} value={key}>
                    {t(`maps.categories.${key}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="place-unit">{t('maps.places.fields.unit')}</FieldLabel>
            <Input
              id="place-unit"
              value={unit}
              onChange={(event) => setUnit(event.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="place-key">{t('maps.places.fields.placeKey')}</FieldLabel>
            <Input
              id="place-key"
              value={placeKey}
              onChange={(event) => setPlaceKey(event.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel>{t('maps.detail.status')}</FieldLabel>
            <Select
              value={status}
              onValueChange={(value) => setStatus(value as PlaceStatus)}
            >
              <SelectTrigger aria-label={t('maps.detail.status')}>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="pending">{t('maps.status.pending')}</SelectItem>
                <SelectItem value="active">{t('maps.status.active')}</SelectItem>
                <SelectItem value="hidden">{t('maps.status.hidden')}</SelectItem>
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="place-lat">{t('maps.locations.fields.lat')}</FieldLabel>
            <Input
              id="place-lat"
              value={lat}
              onChange={(event) => setLat(event.target.value)}
              inputMode="decimal"
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="place-lng">{t('maps.locations.fields.lng')}</FieldLabel>
            <Input
              id="place-lng"
              value={lng}
              onChange={(event) => setLng(event.target.value)}
              inputMode="decimal"
            />
          </Field>
          <DialogFooter className="sm:col-span-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {t('access.actions.cancel')}
            </Button>
            <Button type="submit" disabled={isPending || !name.trim() || !locationId}>
              {isPending ? <Spinner className="size-4" /> : null}
              {mode === 'edit' ? t('access.actions.save') : t('maps.places.createAction')}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
