import type { LocationRecord, PlacePin } from './types'
import type { LatLng, MapPin } from './providers/map-provider'
import {
  CREATE_DRAFT_PIN_ID,
  USER_LOCATION_PIN_ID,
} from './providers/map-provider'

export const LOCATION_PIN_ID_PREFIX = 'location:'

export function toLocationPinId(locationId: string) {
  return `${LOCATION_PIN_ID_PREFIX}${locationId}`
}

export function isLocationPinId(pinId: string) {
  return pinId.startsWith(LOCATION_PIN_ID_PREFIX)
}

export function fromLocationPinId(pinId: string) {
  return isLocationPinId(pinId)
    ? pinId.slice(LOCATION_PIN_ID_PREFIX.length)
    : null
}

export type PinPreview = {
  kind: 'place' | 'location'
  title: string
  subtitle: string
  category?: string
  position: LatLng
}

export function previewForPin(input: {
  pinId: string
  places: PlacePin[]
  locations: LocationRecord[]
}): PinPreview | null {
  const locationId = fromLocationPinId(input.pinId)
  if (locationId) {
    const location = input.locations.find((item) => item.id === locationId)
    if (!location || !hasCoords(location.lat, location.lng)) return null
    return {
      kind: 'location',
      title: location.name,
      subtitle: location.formatted || location.street || '',
      position: { lat: location.lat as number, lng: location.lng as number },
    }
  }
  const place = input.places.find((item) => item.id === input.pinId)
  if (!place || !hasCoords(place.lat, place.lng)) return null
  return {
    kind: 'place',
    title: place.name,
    subtitle: place.locationName,
    category: String(place.category),
    position: { lat: place.lat as number, lng: place.lng as number },
  }
}

function hasCoords(lat?: number | null, lng?: number | null) {
  return lat != null && lng != null
}

export function buildMapPins(input: {
  places: PlacePin[]
  locations: LocationRecord[]
  countryCode?: string
  userLocation: LatLng | null
  draftPoint?: LatLng | null
  labels: { user: string; draft: string }
}): MapPin[] {
  const items: MapPin[] = []
  const coveredLocationIds = new Set<string>()

  for (const place of input.places) {
    if (!hasCoords(place.lat, place.lng)) continue
    coveredLocationIds.add(place.locationId)
    items.push({
      id: place.id,
      position: { lat: place.lat as number, lng: place.lng as number },
      title: place.name,
      category: String(place.category),
    })
  }

  const country = input.countryCode?.trim().toUpperCase()
  for (const location of input.locations) {
    if (!hasCoords(location.lat, location.lng)) continue
    if (coveredLocationIds.has(location.id)) continue
    if (country && location.countryCode && location.countryCode !== country) {
      continue
    }
    items.push({
      id: toLocationPinId(location.id),
      position: { lat: location.lat as number, lng: location.lng as number },
      title: location.name,
      variant: 'location',
    })
  }

  if (input.userLocation) {
    items.push({
      id: USER_LOCATION_PIN_ID,
      position: input.userLocation,
      title: input.labels.user,
      interactive: false,
      variant: 'user',
    })
  }

  if (input.draftPoint) {
    items.push({
      id: CREATE_DRAFT_PIN_ID,
      position: input.draftPoint,
      title: input.labels.draft,
      interactive: false,
      variant: 'draft',
    })
  }

  return spreadOverlappingPins(items)
}

const OVERLAP_RADIUS_DEG = 0.00014

function spreadOverlappingPins(pins: MapPin[]): MapPin[] {
  const groups = new Map<string, number[]>()
  pins.forEach((pin, index) => {
    if (pin.interactive === false) return
    const key = `${pin.position.lat.toFixed(5)},${pin.position.lng.toFixed(5)}`
    const ids = groups.get(key) ?? []
    ids.push(index)
    groups.set(key, ids)
  })

  const next = pins.map((pin) => ({ ...pin, position: { ...pin.position } }))
  for (const indexes of groups.values()) {
    if (indexes.length < 2) continue
    indexes.forEach((index, i) => {
      const pin = next[index]
      const angle = (2 * Math.PI * i) / indexes.length - Math.PI / 2
      const latRad = (pin.position.lat * Math.PI) / 180
      pin.position = {
        lat: pin.position.lat + OVERLAP_RADIUS_DEG * Math.cos(angle),
        lng:
          pin.position.lng +
          (OVERLAP_RADIUS_DEG * Math.sin(angle)) / Math.cos(latRad),
      }
    })
  }
  return next
}
