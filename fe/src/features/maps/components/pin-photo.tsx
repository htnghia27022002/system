'use client'

import { mapPreviewImageUrl } from '../providers/create-map-provider'
import type { LatLng } from '../providers/map-provider'

type PinPhotoProps = {
  position: LatLng
  alt: string
  className?: string
}

export function PinPhoto({ position, alt, className }: PinPhotoProps) {
  const src = mapPreviewImageUrl(position)
  if (!src) {
    return <div className={className} aria-hidden />
  }
  return (
    // Preview snapshot of the pin area; OSM tiles already load in the map adapter.
    // eslint-disable-next-line @next/next/no-img-element
    <img src={src} alt={alt} className={className} draggable={false} />
  )
}
