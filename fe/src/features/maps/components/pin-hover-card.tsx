'use client'

import { useTranslation } from 'react-i18next'

import type { PinPreview } from '../map-pins'
import type { MapScreenPoint } from '../providers/map-provider'
import { PinPhoto } from './pin-photo'

type PinHoverCardProps = {
  preview: PinPreview
  screen: MapScreenPoint
}

export function PinHoverCard({ preview, screen }: PinHoverCardProps) {
  const { t } = useTranslation('admin')
  const kindLabel =
    preview.kind === 'location'
      ? t('maps.locations.title')
      : t(`maps.categories.${preview.category ?? 'uncategorized'}`, {
          defaultValue: preview.category,
        })

  return (
    <div
      role="tooltip"
      className="pointer-events-none absolute z-[1150] w-[15.5rem] overflow-hidden rounded-xl border border-border bg-card shadow-lg"
      style={{ left: screen.x + 16, top: screen.y - 8 }}
    >
      <div className="relative h-24 bg-muted">
        <PinPhoto
          position={preview.position}
          alt=""
          className="size-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-t from-card/90 to-transparent" />
      </div>
      <div className="space-y-1 px-3 py-2.5">
        <p className="truncate text-sm font-semibold text-foreground">
          {preview.title}
        </p>
        <p className="truncate text-xs text-muted-foreground">{kindLabel}</p>
        {preview.subtitle ? (
          <p className="truncate text-xs text-muted-foreground">
            {preview.subtitle}
          </p>
        ) : null}
      </div>
    </div>
  )
}
