export type LatLng = {
  lat: number
  lng: number
}

export type MapPinVariant = 'default' | 'user' | 'draft' | 'location'

export type MapPin = {
  id: string
  position: LatLng
  title?: string
  interactive?: boolean
  variant?: MapPinVariant
  category?: string
}

export const USER_LOCATION_PIN_ID = 'user-location'
export const CREATE_DRAFT_PIN_ID = 'create-draft'

export type MapScreenPoint = {
  x: number
  y: number
}

export type MapBasemap = 'street' | 'satellite'

export type MapViewport = {
  center: LatLng
  zoom: number
}

export const DEFAULT_MAP_VIEWPORT: MapViewport = {
  center: { lat: 16.047079, lng: 108.20623 },
  zoom: 6,
}

/**
 * Vendor-agnostic map port. Product UI talks only to this interface.
 * Persistence and domain models stay lat/lng — no OSM or Google types here.
 */
export type MapProvider = {
  mount: (container: HTMLElement, viewport?: MapViewport) => void
  unmount: () => void
  setCenter: (center: LatLng) => void
  setZoom: (zoom: number) => void
  setBasemap: (basemap: MapBasemap) => void
  addPin: (pin: MapPin) => void
  updatePin: (pin: MapPin) => void
  removePin: (id: string) => void
  clearPins: () => void
  onPinClick: (handler: (pinId: string) => void) => void
  onPinHover: (
    handler: (pinId: string | null, screen?: MapScreenPoint) => void,
  ) => void
  onMapClick: (handler: (point: LatLng, screen: MapScreenPoint) => void) => void
  fitBounds: (positions: LatLng[]) => void
}
