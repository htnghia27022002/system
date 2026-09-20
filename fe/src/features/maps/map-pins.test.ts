import { describe, expect, it } from 'vitest'

import {
  buildMapPins,
  isLocationPinId,
  previewForPin,
  toLocationPinId,
} from './map-pins'
import type { LocationRecord, PlacePin } from './types'

const labels = { user: 'You', draft: 'Draft' }

function location(partial: Partial<LocationRecord> & Pick<LocationRecord, 'id' | 'name'>): LocationRecord {
  return {
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    countryCode: 'VN',
    lat: 10.38,
    lng: 106.42,
    ...partial,
  }
}

function place(partial: Partial<PlacePin> & Pick<PlacePin, 'id' | 'name' | 'locationId'>): PlacePin {
  return {
    categoryId: 'cat-1',
    category: 'eatery',
    status: 'active',
    lat: 16.04,
    lng: 108.2,
    locationName: 'Site',
    ...partial,
  }
}

describe('buildMapPins', () => {
  it('shows a Location pin when the site has coordinates and no Place yet', () => {
    const pins = buildMapPins({
      places: [],
      locations: [location({ id: 'loc-1', name: 'My Tho site' })],
      countryCode: 'VN',
      userLocation: null,
      labels,
    })

    expect(pins).toEqual([
      {
        id: toLocationPinId('loc-1'),
        position: { lat: 10.38, lng: 106.42 },
        title: 'My Tho site',
        variant: 'location',
      },
    ])
    expect(isLocationPinId(pins[0].id)).toBe(true)
    expect(
      previewForPin({
        pinId: toLocationPinId('loc-1'),
        places: [],
        locations: [location({ id: 'loc-1', name: 'My Tho site' })],
      })?.title,
    ).toBe('My Tho site')
  })

  it('hides a Location pin once a Place at that Location is pinned', () => {
    const pins = buildMapPins({
      places: [
        place({
          id: 'place-1',
          name: 'Cafe',
          locationId: 'loc-1',
          lat: 10.38,
          lng: 106.42,
        }),
      ],
      locations: [location({ id: 'loc-1', name: 'My Tho site' })],
      userLocation: null,
      labels,
    })

    expect(pins.map((pin) => pin.id)).toEqual(['place-1'])
  })

  it('offsets two Places that share the same coordinates so both pins stay clickable', () => {
    const pins = buildMapPins({
      places: [
        place({
          id: 'place-1',
          name: 'Cafe',
          locationId: 'loc-1',
          lat: 10.38,
          lng: 106.42,
        }),
        place({
          id: 'place-2',
          name: 'Shop',
          locationId: 'loc-1',
          lat: 10.38,
          lng: 106.42,
        }),
      ],
      locations: [location({ id: 'loc-1', name: 'My Tho site' })],
      userLocation: null,
      labels,
    })

    expect(pins.map((pin) => pin.id)).toEqual(['place-1', 'place-2'])
    expect(pins[0].position).not.toEqual(pins[1].position)
  })
})
