import { DEFAULT_MAP_VIEWPORT, type LatLng, type MapBasemap } from './map-provider'

function tileXY(lat: number, lng: number, zoom: number) {
  const n = 2 ** zoom
  const x = Math.floor(((lng + 180) / 360) * n)
  const latRad = (lat * Math.PI) / 180
  const y = Math.floor(
    ((1 - Math.log(Math.tan(latRad) + 1 / Math.cos(latRad)) / Math.PI) / 2) * n,
  )
  return { x, y }
}

/** Street-map snapshot for pin cards. Tile URLs stay in the provider layer. */
export function osmPreviewImageUrl(position: LatLng, zoom = 16) {
  const { x, y } = tileXY(position.lat, position.lng, zoom)
  return `https://tile.openstreetmap.org/${zoom}/${x}/${y}.png`
}

export function satellitePreviewImageUrl(position: LatLng, zoom = 16) {
  const { x, y } = tileXY(position.lat, position.lng, zoom)
  return `https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/${zoom}/${y}/${x}`
}

export function basemapPreviewImageUrl(
  basemap: MapBasemap,
  position: LatLng = DEFAULT_MAP_VIEWPORT.center,
) {
  return basemap === 'satellite'
    ? satellitePreviewImageUrl(position, 14)
    : osmPreviewImageUrl(position, 14)
}
