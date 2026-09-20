'use client'

import { useQuery } from '@tanstack/react-query'

import { mapsApi } from '../services/maps-api'
import type { MapPlaceFilters } from '../types'
import { mapsKeys } from './maps-keys'

export function useMapPlaces(filters: MapPlaceFilters) {
  return useQuery({
    queryKey: mapsKeys.places(filters),
    queryFn: () => mapsApi.listPlaces(filters),
    meta: { skipNavLoading: true },
  })
}
