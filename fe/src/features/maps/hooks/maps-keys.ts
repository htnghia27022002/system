import type { ListLocationsParams, MapPlaceFilters } from '../types'

export const mapsKeys = {
  all: ['admin', 'maps'] as const,
  places: (filters: MapPlaceFilters) =>
    ['admin', 'maps', 'places', filters] as const,
  managedPlaces: (locationId?: string) =>
    ['admin', 'maps', 'places', 'manage', locationId ?? 'all'] as const,
  place: (id: string) => ['admin', 'maps', 'places', 'detail', id] as const,
  locations: (params?: ListLocationsParams) =>
    ['admin', 'maps', 'locations', params ?? {}] as const,
  sources: (page?: number) => ['admin', 'maps', 'sources', page] as const,
  ingest: ['admin', 'maps', 'ingest'] as const,
  categories: ['admin', 'maps', 'categories'] as const,
}
