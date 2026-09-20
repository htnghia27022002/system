import { describe, expect, it } from 'vitest'

import { createMapProvider } from './create-map-provider'

describe('createMapProvider', () => {
  it('returns an OSM adapter that implements the map port', () => {
    const provider = createMapProvider({ provider: 'osm' })
    expect(typeof provider.mount).toBe('function')
    expect(typeof provider.unmount).toBe('function')
    expect(typeof provider.addPin).toBe('function')
    expect(typeof provider.clearPins).toBe('function')
    expect(typeof provider.onPinClick).toBe('function')
    expect(typeof provider.onMapClick).toBe('function')
    expect(typeof provider.setBasemap).toBe('function')
    expect(typeof provider.fitBounds).toBe('function')
  })

  it('returns a Google stub that does not throw when the key is missing', () => {
    const container = document.createElement('div')
    const provider = createMapProvider({
      provider: 'google',
      googleMapsApiKey: '',
    })
    expect(() => provider.mount(container)).not.toThrow()
    expect(container.textContent).toContain('Google Maps is not configured')
  })

  it('keeps the Google stub configured-not-ready even when a key is present', () => {
    const container = document.createElement('div')
    const provider = createMapProvider({
      provider: 'google',
      googleMapsApiKey: 'test-key',
    })
    expect(() => {
      provider.mount(container)
      provider.addPin({
        id: 'place-a',
        position: { lat: 10.77, lng: 106.7 },
      })
      provider.addPin({
        id: 'place-b',
        position: { lat: 10.77, lng: 106.7 },
      })
    }).not.toThrow()
    expect(container.textContent).toContain('Google Maps is not configured')
  })
})
