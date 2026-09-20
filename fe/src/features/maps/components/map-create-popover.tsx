'use client'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Spinner } from '@/components/ui/spinner'

import type { LocationWriteInput } from '../types'
import type { LatLng } from '../providers/map-provider'

type MapCreatePopoverProps = {
  point: LatLng
  isPending: boolean
  onSubmit: (input: LocationWriteInput) => void
  onClose: () => void
}

export function MapCreatePopover({
  point,
  isPending,
  onSubmit,
  onClose,
}: MapCreatePopoverProps) {
  const { t } = useTranslation('admin')
  const [name, setName] = useState('')

  useEffect(() => {
    setName('')
  }, [point.lat, point.lng])

  const coordLabel = `${point.lat.toFixed(5)}, ${point.lng.toFixed(5)}`

  return (
    <div className="pointer-events-none absolute inset-x-3 bottom-4 z-[1200] flex justify-center">
      <div
        role="form"
        aria-label={t('maps.create.popoverTitle')}
        className="pointer-events-auto w-full max-w-md rounded-2xl border border-border bg-card p-3 shadow-xl"
      >
        <p className="mb-1 text-sm font-semibold">{t('maps.create.popoverTitle')}</p>
        <p className="mb-3 text-xs text-muted-foreground">
          {t('maps.create.coords', { coords: coordLabel })}
        </p>
        <form
          className="grid gap-2"
          onSubmit={(event) => {
            event.preventDefault()
            const trimmed = name.trim()
            if (!trimmed) return
            onSubmit({
              name: trimmed,
              countryCode: 'VN',
              lat: point.lat,
              lng: point.lng,
            })
          }}
        >
          <Input
            id="map-create-name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder={t('maps.create.namePlaceholder')}
            aria-label={t('maps.locations.fields.name')}
            required
            autoFocus
          />
          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" size="sm" onClick={onClose}>
              {t('access.actions.cancel')}
            </Button>
            <Button type="submit" size="sm" disabled={isPending || !name.trim()}>
              {isPending ? <Spinner className="size-4" /> : null}
              {t('maps.locations.createAction')}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
