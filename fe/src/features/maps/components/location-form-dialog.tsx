'use client'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { AddressSelect } from '@/features/address'
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
import { Spinner } from '@/components/ui/spinner'

import type { LocationRecord, LocationWriteInput } from '../types'

type LocationFormDialogProps = {
  open: boolean
  mode: 'create' | 'edit'
  location?: LocationRecord
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (input: LocationWriteInput) => void
}

export function LocationFormDialog({
  open,
  mode,
  location,
  isPending,
  onOpenChange,
  onSubmit,
}: LocationFormDialogProps) {
  const { t } = useTranslation('admin')
  const [name, setName] = useState('')
  const [locationKey, setLocationKey] = useState('')
  const [countryCode, setCountryCode] = useState('VN')
  const [adminDivisionId, setAdminDivisionId] = useState('')
  const [street, setStreet] = useState('')
  const [postalCode, setPostalCode] = useState('')
  const [formatted, setFormatted] = useState('')
  const [lat, setLat] = useState('')
  const [lng, setLng] = useState('')

  useEffect(() => {
    if (!open) return
    setName(location?.name ?? '')
    setLocationKey(location?.locationKey ?? '')
    setCountryCode(location?.countryCode ?? 'VN')
    setAdminDivisionId(location?.adminDivisionId ?? '')
    setStreet(location?.street ?? '')
    setPostalCode(location?.postalCode ?? '')
    setFormatted(location?.formatted ?? '')
    setLat(location?.lat != null ? String(location.lat) : '')
    setLng(location?.lng != null ? String(location.lng) : '')
  }, [open, location])

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
            {mode === 'edit'
              ? t('maps.locations.editTitle')
              : t('maps.locations.createTitle')}
          </DialogTitle>
        </DialogHeader>
        <form
          className="grid gap-4 sm:grid-cols-2"
          onSubmit={(event) => {
            event.preventDefault()
            onSubmit({
              name: name.trim(),
              locationKey: locationKey.trim() || undefined,
              countryCode: countryCode || undefined,
              adminDivisionId: adminDivisionId || undefined,
              street: street.trim(),
              postalCode: postalCode.trim(),
              formatted: formatted.trim(),
              lat: parseCoord(lat),
              lng: parseCoord(lng),
            })
          }}
        >
          <Field className="sm:col-span-2">
            <FieldLabel htmlFor="loc-name">{t('maps.locations.fields.name')}</FieldLabel>
            <Input
              id="loc-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              required
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="loc-key">{t('maps.locations.fields.locationKey')}</FieldLabel>
            <Input
              id="loc-key"
              value={locationKey}
              onChange={(event) => setLocationKey(event.target.value)}
            />
          </Field>
          <AddressSelect
            className="sm:col-span-2 sm:grid-cols-2"
            allowEmpty
            value={{ countryCode, adminDivisionId }}
            onChange={(next) => {
              setCountryCode(next.countryCode)
              setAdminDivisionId(next.adminDivisionId)
            }}
            contentClassName="z-[1500]"
          />
          <Field>
            <FieldLabel htmlFor="loc-street">{t('maps.locations.fields.street')}</FieldLabel>
            <Input
              id="loc-street"
              value={street}
              onChange={(event) => setStreet(event.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="loc-postal">{t('maps.locations.fields.postalCode')}</FieldLabel>
            <Input
              id="loc-postal"
              value={postalCode}
              onChange={(event) => setPostalCode(event.target.value)}
            />
          </Field>
          <Field className="sm:col-span-2">
            <FieldLabel htmlFor="loc-formatted">
              {t('maps.locations.fields.formatted')}
            </FieldLabel>
            <Input
              id="loc-formatted"
              value={formatted}
              onChange={(event) => setFormatted(event.target.value)}
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="loc-lat">{t('maps.locations.fields.lat')}</FieldLabel>
            <Input
              id="loc-lat"
              value={lat}
              onChange={(event) => setLat(event.target.value)}
              inputMode="decimal"
            />
          </Field>
          <Field>
            <FieldLabel htmlFor="loc-lng">{t('maps.locations.fields.lng')}</FieldLabel>
            <Input
              id="loc-lng"
              value={lng}
              onChange={(event) => setLng(event.target.value)}
              inputMode="decimal"
            />
          </Field>
          <DialogFooter className="sm:col-span-2">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {t('access.actions.cancel')}
            </Button>
            <Button type="submit" disabled={isPending || !name.trim()}>
              {isPending ? <Spinner className="size-4" /> : null}
              {mode === 'edit' ? t('access.actions.save') : t('maps.locations.createAction')}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
