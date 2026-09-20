'use client'

import { useEffect, useRef } from 'react'

import { cn } from '@/lib/utils'

import { createMapProvider } from '../providers/create-map-provider'
import {
  DEFAULT_MAP_VIEWPORT,
  type LatLng,
  type MapBasemap,
  type MapPin,
  type MapProvider,
  type MapScreenPoint,
} from '../providers/map-provider'

type MapCanvasProps = {
  pins: MapPin[]
  onPinClick: (pinId: string) => void
  onPinHover?: (pinId: string | null, screen?: MapScreenPoint) => void
  onMapClick?: (point: LatLng, screen: MapScreenPoint) => void
  focus?: LatLng | null
  pickCursor?: boolean
  preserveViewport?: boolean
  basemap?: MapBasemap
  className?: string
}

export function MapCanvas({
  pins,
  onPinClick,
  onPinHover,
  onMapClick,
  focus,
  pickCursor = false,
  preserveViewport = false,
  basemap = 'street',
  className,
}: MapCanvasProps) {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const providerRef = useRef<MapProvider | null>(null)
  const onPinClickRef = useRef(onPinClick)
  const onPinHoverRef = useRef(onPinHover)
  const onMapClickRef = useRef(onMapClick)
  const skipFitRef = useRef(false)
  const fittedRef = useRef(false)

  useEffect(() => {
    onPinClickRef.current = onPinClick
  }, [onPinClick])

  useEffect(() => {
    onPinHoverRef.current = onPinHover
  }, [onPinHover])

  useEffect(() => {
    onMapClickRef.current = onMapClick
  }, [onMapClick])

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const provider = createMapProvider()
    providerRef.current = provider
    provider.onPinClick((pinId) => {
      onPinClickRef.current(pinId)
    })
    provider.onPinHover((pinId, screen) => {
      onPinHoverRef.current?.(pinId, screen)
    })
    provider.onMapClick((point, screen) => {
      onMapClickRef.current?.(point, screen)
    })
    provider.mount(container, DEFAULT_MAP_VIEWPORT)

    return () => {
      provider.unmount()
      providerRef.current = null
    }
  }, [])

  useEffect(() => {
    const provider = providerRef.current
    if (!provider) return

    provider.clearPins()
    for (const pin of pins) {
      provider.addPin(pin)
    }
    if (skipFitRef.current || preserveViewport) {
      skipFitRef.current = false
      return
    }
    const fitPins = pins.filter((pin) => pin.interactive !== false)
    if (fitPins.length === 0) {
      fittedRef.current = false
      return
    }
    if (fittedRef.current) return
    fittedRef.current = true
    provider.fitBounds(fitPins.map((pin) => pin.position))
  }, [pins, preserveViewport])

  useEffect(() => {
    const provider = providerRef.current
    if (!provider || !focus) return
    skipFitRef.current = true
    provider.setCenter(focus)
    provider.setZoom(16)
  }, [focus])

  useEffect(() => {
    providerRef.current?.setBasemap(basemap)
  }, [basemap])

  return (
    <div
      ref={containerRef}
      className={cn(
        'h-full min-h-[360px] w-full overflow-hidden bg-muted/30',
        pickCursor && 'maps-map-pick !cursor-crosshair',
        className,
      )}
    />
  )
}
