'use client'

import { useQuery } from '@tanstack/react-query'

import { mapsApi } from '../services/maps-api'
import { mapsKeys } from './maps-keys'

export function useMapCategories() {
  return useQuery({
    queryKey: mapsKeys.categories,
    queryFn: () => mapsApi.listCategories(),
    meta: { skipNavLoading: true },
  })
}
