/**
 * Google Maps MapProvider stub (P1).
 *
 * This file satisfies the same port as the live OSM adapter so the product UI
 * never imports Google types. Persistence stays lat/lng only.
 *
 * How to replace this stub with a live adapter:
 * 1. Load the Maps JavaScript API using NEXT_PUBLIC_GOOGLE_MAPS_API_KEY
 *    from fe/.env (never root .env or be/.env).
 * 2. Map LatLng / MapPin / MapViewport onto google.maps objects ONLY here.
 * 3. Keep registration in create-map-provider.ts (`case 'google'`).
 * 4. Do not leak Google types into pages, hooks, or domain models.
 * 5. mount() must still fail closed: missing key or load error shows
 *    "Google Maps is not configured" and must not throw uncaught.
 */
import type {
  LatLng,
  MapPin,
  MapProvider,
  MapViewport,
} from './map-provider'

const NOT_CONFIGURED = 'Google Maps is not configured'

export function createGoogleMapProvider(_apiKey?: string): MapProvider {
  return {
    mount(container: HTMLElement, _viewport?: MapViewport) {
      container.replaceChildren()
      const message = document.createElement('div')
      message.setAttribute('role', 'status')
      message.className =
        'flex h-full min-h-[360px] items-center justify-center p-4 text-center text-sm text-muted-foreground'
      message.textContent = NOT_CONFIGURED
      container.appendChild(message)
    },
    unmount() {},
    setCenter(_center: LatLng) {},
    setZoom(_zoom: number) {},
    setBasemap(_basemap) {},
    addPin(_pin: MapPin) {},
    updatePin(_pin: MapPin) {},
    removePin(_id: string) {},
    clearPins() {},
    onPinClick(_handler: (pinId: string) => void) {},
    onPinHover(_handler) {},
    onMapClick(_handler) {},
    fitBounds(_positions: LatLng[]) {},
  }
}
