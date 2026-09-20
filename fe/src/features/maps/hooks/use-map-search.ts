'use client'

import { useQuery } from '@tanstack/react-query'

import { useDebouncedValue } from '@/hooks/use-debounced-value'

import {
  matchSavedLocations,
  matchSavedPlaces,
  mergeMapSearchHits,
  type MapSearchHit,
} from '../map-search'
import { mapsApi } from '../services/maps-api'
import type { LocationRecord, PlacePin } from '../types'
import { mapsKeys } from './maps-keys'

const MIN_QUERY = 2

function toHit(item: {
  id: string
  kind: MapSearchHit['kind']
  title: string
  subtitle: string
  lat: number
  lng: number
  placeId?: string
  locationId?: string
}): MapSearchHit {
  return {
    id: item.id,
    kind: item.kind,
    title: item.title,
    subtitle: item.subtitle,
    position: { lat: item.lat, lng: item.lng },
    placeId: item.placeId,
    locationId: item.locationId,
  }
}

export function useMapSearch(
  query: string,
  loaded: {
    places: PlacePin[]
    locations: LocationRecord[]
    countryCode?: string
  },
) {
  const immediate = query.trim()
  const debounced = useDebouncedValue(immediate, 400)
  const remoteEnabled = debounced.length >= MIN_QUERY

  const remote = useQuery({
    queryKey: [
      ...mapsKeys.all,
      'search',
      debounced,
      loaded.countryCode ?? '',
    ],
    queryFn: () =>
      mapsApi.search({
        q: debounced,
        countryCode: loaded.countryCode,
      }),
    enabled: remoteEnabled,
    staleTime: 30_000,
    meta: { skipNavLoading: true },
  })

  const hits: MapSearchHit[] = mergeMapSearchHits([
    matchSavedPlaces(loaded.places, immediate),
    matchSavedLocations(loaded.locations, immediate),
    (remote.data?.items ?? []).map(toHit),
  ])

  return {
    hits,
    loading: remoteEnabled && remote.isFetching,
    open: immediate.length >= 1,
  }
}
