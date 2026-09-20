import { describe, expect, it } from 'vitest'

import {
  matchSavedLocations,
  matchSavedPlaces,
  mergeMapSearchHits,
} from './map-search'
import type { LocationRecord, PlacePin } from './types'

const place = (partial: Partial<PlacePin> & Pick<PlacePin, 'id' | 'name'>): PlacePin => ({
  categoryId: 'cat',
  category: 'eatery',
  status: 'active',
  locationId: 'loc-1',
  locationName: 'Han Market',
  lat: 16.07,
  lng: 108.22,
  ...partial,
})

const location = (
  partial: Partial<LocationRecord> & Pick<LocationRecord, 'id' | 'name'>,
): LocationRecord => ({
  createdAt: '2026-01-01T00:00:00Z',
  updatedAt: '2026-01-01T00:00:00Z',
  lat: 16.047,
  lng: 108.206,
  formatted: 'Hai Chau, Da Nang',
  ...partial,
})

describe('map search helpers', () => {
  it('matches saved Places by name or location and skips items without coordinates', () => {
    const hits = matchSavedPlaces(
      [
        place({ id: 'p1', name: 'Com Tam Demo' }),
        place({ id: 'p2', name: 'Pho', locationName: 'Han Market' }),
        place({ id: 'p3', name: 'Hidden', lat: null, lng: null }),
      ],
      'han',
    )

    expect(hits.map((hit) => hit.placeId)).toEqual(['p1', 'p2'])
    expect(hits[0]?.kind).toBe('place')
  })

  it('matches saved Locations by name or address', () => {
    const hits = matchSavedLocations(
      [
        location({
          id: 'l1',
          name: 'Bach Dang hotel',
          formatted: 'Bach Dang, Da Nang',
        }),
        location({ id: 'l2', name: 'Office', formatted: 'Hai Chau, Da Nang' }),
        location({ id: 'l3', name: 'No pin', lat: null, lng: null }),
      ],
      'chau',
    )

    expect(hits.map((hit) => hit.locationId)).toEqual(['l2'])
    expect(hits[0]?.kind).toBe('location')
  })

  it('dedupes Place and Location hits when merging groups', () => {
    const placeHit = {
      id: 'place:p1',
      kind: 'place' as const,
      title: 'Pho',
      subtitle: 'Han Market',
      position: { lat: 16, lng: 108 },
      placeId: 'p1',
    }
    const osmHit = {
      id: 'geocode:osm:99',
      kind: 'geocode' as const,
      title: 'Da Nang',
      subtitle: 'Vietnam',
      position: { lat: 16.05, lng: 108.2 },
    }

    const merged = mergeMapSearchHits([
      [placeHit],
      [{ ...placeHit, id: 'place:p1-remote' }],
      [osmHit],
    ])

    expect(merged.map((hit) => hit.id)).toEqual(['place:p1', 'geocode:osm:99'])
  })
})
