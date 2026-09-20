'use client'

import dynamic from 'next/dynamic'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import {
  PermissionKeys,
  usePermissions,
} from '@/features/access-control'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'

import { hasActivePlaceFilters, toMapPlaceFilters } from '../category-filters'
import { useMapSearch } from '../hooks/use-map-search'
import type { MapSearchHit } from '../map-search'
import { isIngestInProgress, useLatestIngest, useStartIngest } from '../hooks/use-map-ingest'
import {
  useMapLocationMutations,
  useMapLocations,
} from '../hooks/use-map-locations'
import { useMapPlaces } from '../hooks/use-map-places'
import { buildMapPins, fromLocationPinId, previewForPin } from '../map-pins'
import {
  CREATE_DRAFT_PIN_ID,
  USER_LOCATION_PIN_ID,
  type LatLng,
  type MapBasemap,
  type MapScreenPoint,
} from '../providers/map-provider'
import { IngestStatus } from './ingest-status'
import { MapActions } from './map-actions'
import { MapCreatePopover } from './map-create-popover'
import { MapFilters, type MapFilterValues } from './map-filters'
import { MapPointTooltip } from './map-point-tooltip'
import { PinDetailCard } from './pin-detail-card'
import { PinHoverCard } from './pin-hover-card'
import { SourcesPanel } from './sources-panel'

const MapCanvas = dynamic(
  () => import('./map-canvas').then((mod) => mod.MapCanvas),
  { ssr: false },
)

const ALL = 'all'

type ManagePanel = 'sources' | null

type CreateDraft = {
  point: LatLng
  screen: MapScreenPoint
  stage: 'inspect' | 'form'
}

export function MapsPage() {
  const { t } = useTranslation('admin')
  const { hasPermission } = usePermissions()
  const canModify = hasPermission(PermissionKeys.maps.modify)
  const mapFrameRef = useRef<HTMLDivElement | null>(null)

  const [filterValues, setFilterValues] = useState<MapFilterValues>({
    q: '',
    category: ALL,
    countryCode: 'VN',
    adminDivisionId: '',
  })
  const [selectedPlaceId, setSelectedPlaceId] = useState<string | null>(null)
  const [selectedLocationId, setSelectedLocationId] = useState<string | null>(
    null,
  )
  const [hoverPin, setHoverPin] = useState<{
    pinId: string
    screen: MapScreenPoint
  } | null>(null)
  const [managePanel, setManagePanel] = useState<ManagePanel>(null)
  const [userLocation, setUserLocation] = useState<LatLng | null>(null)
  const [focus, setFocus] = useState<LatLng | null>(null)
  const [locating, setLocating] = useState(false)
  const [createDraft, setCreateDraft] = useState<CreateDraft | null>(null)
  const [basemap, setBasemap] = useState<MapBasemap>('street')
  const [skipNextFit, setSkipNextFit] = useState(false)
  const [awaitingAutoLocate, setAwaitingAutoLocate] = useState(true)
  const didAutoLocate = useRef(false)

  const filters = toMapPlaceFilters({
    category: filterValues.category,
    adminDivisionId: filterValues.adminDivisionId,
    countryCode: filterValues.countryCode || undefined,
  })

  const placesQuery = useMapPlaces(filters)
  const locationsQuery = useMapLocations()
  const { createLocation } = useMapLocationMutations()
  const ingestQuery = useLatestIngest()
  const startIngest = useStartIngest()

  const ingestRun = ingestQuery.data?.run
  const ingesting = isIngestInProgress(ingestRun?.status)
  const pins = useMemo(
    () => placesQuery.data?.items ?? [],
    [placesQuery.data?.items],
  )
  const filteredEmpty =
    !placesQuery.isLoading &&
    pins.length === 0 &&
    hasActivePlaceFilters({
      ...filters,
      countryCode: filterValues.adminDivisionId
        ? filters.countryCode
        : undefined,
    })
  const closeCreate = useCallback(() => {
    setCreateDraft(null)
  }, [])

  const openManage = useCallback((panel: ManagePanel) => {
    closeCreate()
    setSelectedPlaceId(null)
    setSelectedLocationId(null)
    setHoverPin(null)
    setManagePanel(panel)
  }, [closeCreate])

  const closeSelection = useCallback(() => {
    setSelectedPlaceId(null)
    setSelectedLocationId(null)
    setHoverPin(null)
  }, [])

  const handleMapClick = useCallback(
    (point: LatLng, screen: MapScreenPoint) => {
      setManagePanel(null)
      closeSelection()
      setSkipNextFit(true)
      setCreateDraft({
        point,
        screen,
        stage: 'inspect',
      })
    },
    [closeSelection],
  )

  useEffect(() => {
    if (!createDraft) return
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        closeCreate()
        closeSelection()
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [closeCreate, closeSelection, createDraft])

  useEffect(() => {
    if (createDraft || !skipNextFit) return
    const frame = window.requestAnimationFrame(() => setSkipNextFit(false))
    return () => window.cancelAnimationFrame(frame)
  }, [createDraft, skipNextFit])

  const mapPins = useMemo(
    () =>
      buildMapPins({
        places: pins,
        locations: locationsQuery.data?.items ?? [],
        countryCode: filterValues.countryCode || undefined,
        userLocation,
        draftPoint: createDraft?.point,
        labels: {
          user: t('maps.locate.here'),
          draft: t('maps.create.draftPin'),
        },
      }),
    [
      createDraft?.point,
      filterValues.countryCode,
      locationsQuery.data?.items,
      pins,
      t,
      userLocation,
    ],
  )
  const selectedLocation = useMemo(
    () =>
      (locationsQuery.data?.items ?? []).find(
        (item) => item.id === selectedLocationId,
      ) ?? null,
    [locationsQuery.data?.items, selectedLocationId],
  )
  const hoverPreview = useMemo(() => {
    if (!hoverPin) return null
    return previewForPin({
      pinId: hoverPin.pinId,
      places: pins,
      locations: locationsQuery.data?.items ?? [],
    })
  }, [hoverPin, locationsQuery.data?.items, pins])
  const hoverOffset = useMemo(() => {
    if (!hoverPin) return null
    const width = mapFrameRef.current?.clientWidth ?? 800
    const height = mapFrameRef.current?.clientHeight ?? 600
    return {
      x: Math.min(Math.max(12, hoverPin.screen.x + 14), width - 260),
      y: Math.min(Math.max(12, hoverPin.screen.y - 12), height - 180),
    }
  }, [hoverPin])
  const noPinsYet =
    !placesQuery.isLoading &&
    !locationsQuery.isLoading &&
    !filteredEmpty &&
    !ingesting &&
    mapPins.every((pin) => pin.interactive === false)

  const search = useMapSearch(filterValues.q, {
    places: pins,
    locations: locationsQuery.data?.items ?? [],
    countryCode: filterValues.countryCode || undefined,
  })

  const locate = useCallback(
    (options?: { silent?: boolean }) => {
      const silent = options?.silent === true
      const finish = () => {
        setLocating(false)
        setAwaitingAutoLocate(false)
      }
      if (!navigator.geolocation) {
        finish()
        if (!silent) toast.error(t('maps.locate.unsupported'))
        return
      }
      setLocating(true)
      try {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            const next = {
              lat: position.coords.latitude,
              lng: position.coords.longitude,
            }
            setUserLocation(next)
            setSkipNextFit(true)
            setFocus({ ...next })
            finish()
          },
          (error) => {
            finish()
            if (silent) return
            if (typeof window !== 'undefined' && !window.isSecureContext) {
              toast.error(t('maps.locate.insecure'))
              return
            }
            if (error.code === 1) {
              toast.error(t('maps.locate.denied'))
              return
            }
            if (error.code === 2) {
              toast.error(t('maps.locate.unavailable'))
              return
            }
            if (error.code === 3) {
              toast.error(t('maps.locate.timeout'))
              return
            }
            toast.error(t('maps.locate.failed'))
          },
          { enableHighAccuracy: true, timeout: 15000, maximumAge: 30000 },
        )
      } catch {
        finish()
        if (silent) return
        toast.error(
          typeof window !== 'undefined' && !window.isSecureContext
            ? t('maps.locate.insecure')
            : t('maps.locate.failed'),
        )
      }
    },
    [t],
  )

  useEffect(() => {
    if (didAutoLocate.current) return
    didAutoLocate.current = true
    locate({ silent: true })
  }, [locate])

  const handleSearchSelect = useCallback(
    (hit: MapSearchHit) => {
      setSkipNextFit(true)
      setCreateDraft(null)
      setHoverPin(null)
      setManagePanel(null)
      setFocus({ ...hit.position })
      if (hit.kind === 'place' && hit.placeId) {
        setSelectedLocationId(null)
        setSelectedPlaceId(hit.placeId)
        return
      }
      if (hit.kind === 'location' && hit.locationId) {
        setSelectedPlaceId(null)
        setSelectedLocationId(hit.locationId)
        return
      }
      closeSelection()
    },
    [closeSelection],
  )

  const emptyMessage = noPinsYet
    ? canModify
      ? t('maps.create.emptyHint')
      : t('maps.emptyPins')
    : null

  return (
    <div
      ref={mapFrameRef}
      className="relative h-[calc(100svh-4rem)] min-h-0 flex-1 overflow-hidden"
    >
      <MapCanvas
        pins={mapPins}
        focus={focus}
        pickCursor={canModify}
        basemap={basemap}
        preserveViewport={
          Boolean(createDraft) || skipNextFit || awaitingAutoLocate
        }
        onPinClick={(pinId) => {
          if (pinId === USER_LOCATION_PIN_ID || pinId === CREATE_DRAFT_PIN_ID) {
            return
          }
          setCreateDraft(null)
          setHoverPin(null)
          setManagePanel(null)
          const locationId = fromLocationPinId(pinId)
          if (locationId) {
            setSelectedPlaceId(null)
            setSelectedLocationId(locationId)
            return
          }
          setSelectedLocationId(null)
          setSelectedPlaceId(pinId)
        }}
        onPinHover={(pinId, screen) => {
          if (!pinId || !screen) {
            setHoverPin(null)
            return
          }
          setHoverPin({ pinId, screen })
        }}
        onMapClick={handleMapClick}
        className="absolute inset-0 h-full min-h-0"
      />

      <MapFilters
        values={filterValues}
        locating={locating}
        onLocate={() => locate()}
        searchHits={search.hits}
        searchLoading={search.loading}
        onSearchSelect={handleSearchSelect}
        basemap={basemap}
        onBasemapChange={setBasemap}
        onChange={setFilterValues}
      />

      <MapActions
        ingesting={startIngest.isPending || ingesting}
        onLoadData={() => startIngest.mutate()}
        onOpenSources={() => openManage('sources')}
      />

      <div className="pointer-events-none absolute top-[4.75rem] right-3 z-[1100] sm:top-3 sm:right-[4.25rem]">
        <div className="pointer-events-auto">
          <IngestStatus run={ingestRun} />
        </div>
      </div>

      {createDraft?.stage === 'inspect' ? (
        <MapPointTooltip
          point={createDraft.point}
          canCreate={canModify}
          onCreate={() =>
            setCreateDraft((prev) =>
              prev ? { ...prev, stage: 'form' } : prev,
            )
          }
          onClose={closeCreate}
        />
      ) : null}

      {createDraft?.stage === 'form' && canModify ? (
        <MapCreatePopover
          point={createDraft.point}
          isPending={createLocation.isPending}
          onSubmit={(input) => {
            createLocation.mutate(input, {
              onSuccess: (created) => {
                setSkipNextFit(true)
                closeCreate()
                setSelectedLocationId(created.id)
              },
            })
          }}
          onClose={closeCreate}
        />
      ) : null}

      {placesQuery.isError ? (
        <div className="pointer-events-none absolute inset-x-3 bottom-3 z-[1100] flex justify-start">
          <div className="pointer-events-auto flex max-w-sm flex-col gap-2 rounded-lg border border-border bg-card/95 p-3 shadow-md backdrop-blur-sm">
            <p className="text-sm text-muted-foreground">{t('maps.error')}</p>
            <Button
              variant="outline"
              size="sm"
              onClick={() => void placesQuery.refetch()}
            >
              {t('maps.retry')}
            </Button>
          </div>
        </div>
      ) : null}

      {!createDraft && emptyMessage ? (
        <div className="pointer-events-none absolute inset-x-3 bottom-3 z-[1100] flex justify-start">
          <p className="max-w-sm rounded-lg border border-border bg-card/95 px-3 py-2 text-sm text-muted-foreground shadow-md backdrop-blur-sm">
            {emptyMessage}
          </p>
        </div>
      ) : null}
      {!createDraft && filteredEmpty ? (
        <div className="pointer-events-none absolute inset-x-3 bottom-3 z-[1100] flex justify-start">
          <p className="max-w-sm rounded-lg border border-border bg-card/95 px-3 py-2 text-sm text-muted-foreground shadow-md backdrop-blur-sm">
            {t('maps.emptyFilter')}
          </p>
        </div>
      ) : null}

      <Sheet
        modal={false}
        open={managePanel !== null}
        onOpenChange={(open) => {
          if (!open) setManagePanel(null)
        }}
      >
        <SheetContent
          side="right"
          showOverlay={false}
          className="z-[1300] w-full sm:max-w-sm"
        >
          <SheetHeader>
            <SheetTitle>{t('maps.sources.title')}</SheetTitle>
            <SheetDescription>{t('maps.sources.description')}</SheetDescription>
          </SheetHeader>
          <SourcesPanel />
        </SheetContent>
      </Sheet>

      {hoverPreview && hoverOffset && !selectedPlaceId && !selectedLocationId && !createDraft ? (
        <PinHoverCard preview={hoverPreview} screen={hoverOffset} />
      ) : null}

      {selectedPlaceId || selectedLocation ? (
        <PinDetailCard
          placeId={selectedPlaceId}
          location={selectedLocation}
          places={pins}
          onClose={closeSelection}
          onSelectPlace={(placeId) => {
            setSelectedLocationId(null)
            setSelectedPlaceId(placeId)
          }}
          onPlaceCreated={(placeId) => {
            setSkipNextFit(true)
            setSelectedLocationId(null)
            setSelectedPlaceId(placeId)
          }}
        />
      ) : null}
    </div>
  )
}
