import type { LocationRecord, PlacePin } from './types'
import type { LatLng } from './providers/map-provider'

export type MapSearchKind = 'place' | 'location' | 'geocode'

export type MapSearchHit = {
  id: string
  kind: MapSearchKind
  title: string
  subtitle: string
  position: LatLng
  placeId?: string
  locationId?: string
}

function hasCoords(lat?: number | null, lng?: number | null): lat is number {
  return lat != null && lng != null && Number.isFinite(lat) && Number.isFinite(lng)
}

function needle(query: string) {
  return query.trim().toLowerCase()
}

function includesQuery(value: string | null | undefined, query: string) {
  if (!value) return false
  return value.toLowerCase().includes(query)
}

export function matchSavedPlaces(
  places: PlacePin[],
  query: string,
  limit = 6,
): MapSearchHit[] {
  const q = needle(query)
  if (q.length < 1) return []

  const hits: MapSearchHit[] = []
  for (const place of places) {
    if (hits.length >= limit) break
    if (!hasCoords(place.lat, place.lng)) continue
    if (
      !includesQuery(place.name, q) &&
      !includesQuery(place.locationName, q) &&
      !includesQuery(place.unit, q)
    ) {
      continue
    }
    hits.push({
      id: `place:${place.id}`,
      kind: 'place',
      title: place.name,
      subtitle: [place.locationName, place.unit].filter(Boolean).join(' · '),
      position: { lat: place.lat, lng: place.lng },
      placeId: place.id,
    })
  }
  return hits
}

export function matchSavedLocations(
  locations: LocationRecord[],
  query: string,
  limit = 6,
): MapSearchHit[] {
  const q = needle(query)
  if (q.length < 1) return []

  const hits: MapSearchHit[] = []
  for (const location of locations) {
    if (hits.length >= limit) break
    if (!hasCoords(location.lat, location.lng)) continue
    if (
      !includesQuery(location.name, q) &&
      !includesQuery(location.formatted, q) &&
      !includesQuery(location.street, q)
    ) {
      continue
    }
    hits.push({
      id: `location:${location.id}`,
      kind: 'location',
      title: location.name,
      subtitle: location.formatted || location.street || '',
      position: { lat: location.lat, lng: location.lng },
      locationId: location.id,
    })
  }
  return hits
}

export function mergeMapSearchHits(groups: MapSearchHit[][]): MapSearchHit[] {
  const seen = new Set<string>()
  const hits: MapSearchHit[] = []

  for (const group of groups) {
    for (const hit of group) {
      const key = hit.placeId
        ? `place:${hit.placeId}`
        : hit.locationId
          ? `location:${hit.locationId}`
          : hit.id
      if (seen.has(key)) continue
      seen.add(key)
      hits.push(hit)
    }
  }

  return hits
}
