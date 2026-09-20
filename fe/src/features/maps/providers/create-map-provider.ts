import { env } from '@/config/env'

import { createGoogleMapProvider } from './google-map-provider'
import type { LatLng, MapBasemap, MapProvider } from './map-provider'
import {
  basemapPreviewImageUrl,
  osmPreviewImageUrl,
} from './map-preview-image'
import { createOsmMapProvider } from './osm-map-provider'

export type MapProviderKind = 'osm' | 'google'

export type CreateMapProviderOptions = {
  provider?: MapProviderKind
  googleMapsApiKey?: string
}

export function createMapProvider(
  options: CreateMapProviderOptions = {},
): MapProvider {
  const provider = options.provider ?? env.MAP_PROVIDER
  const googleMapsApiKey =
    options.googleMapsApiKey ?? env.GOOGLE_MAPS_API_KEY ?? ''

  if (provider === 'google') {
    return createGoogleMapProvider(googleMapsApiKey)
  }

  return createOsmMapProvider()
}

export function mapPreviewImageUrl(position: LatLng) {
  const provider = env.MAP_PROVIDER
  if (provider === 'google') return undefined
  return osmPreviewImageUrl(position)
}

export function mapBasemapPreviewUrl(basemap: MapBasemap) {
  const provider = env.MAP_PROVIDER
  if (provider === 'google') return undefined
  return basemapPreviewImageUrl(basemap)
}
