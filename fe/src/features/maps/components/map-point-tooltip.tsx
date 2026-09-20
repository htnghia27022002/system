'use client'

import { PlusIcon, XIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'

import { mapPreviewImageUrl } from '../providers/create-map-provider'
import type { LatLng } from '../providers/map-provider'

type MapPointTooltipProps = {
  point: LatLng
  canCreate: boolean
  onCreate: () => void
  onClose: () => void
}

export function MapPointTooltip({
  point,
  canCreate,
  onCreate,
  onClose,
}: MapPointTooltipProps) {
  const { t } = useTranslation('admin')
  const coords = `${point.lat.toFixed(5)}, ${point.lng.toFixed(5)}`
  const preview = mapPreviewImageUrl(point)
  const title = t('maps.create.droppedPin')

  return (
    <div className="pointer-events-none absolute inset-x-3 bottom-4 z-[1200] flex justify-center">
      <div
        role="dialog"
        aria-label={t('maps.create.pointTitle')}
        className="pointer-events-auto flex w-full max-w-md items-center gap-3 rounded-2xl border border-border bg-card p-2 pr-2 shadow-xl"
      >
        {preview ? (
          // Preview tiles already used by the map adapter.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={preview}
            alt=""
            className="size-16 shrink-0 rounded-xl object-cover"
            draggable={false}
          />
        ) : (
          <div className="size-16 shrink-0 rounded-xl bg-muted" />
        )}
        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-semibold text-foreground">
            {title}
          </p>
          <p className="truncate text-xs text-muted-foreground">
            {t('maps.create.locationHint')}
          </p>
          <p className="truncate text-xs text-muted-foreground">{coords}</p>
        </div>
        <div className="flex shrink-0 items-center gap-1">
          {canCreate ? (
            <Button
              type="button"
              size="icon"
              className="rounded-full"
              onClick={onCreate}
              aria-label={t('maps.create.here')}
            >
              <PlusIcon className="size-4" />
            </Button>
          ) : null}
          <Button
            type="button"
            size="icon"
            variant="ghost"
            className="rounded-full"
            onClick={onClose}
            aria-label={t('access.actions.cancel')}
          >
            <XIcon className="size-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
