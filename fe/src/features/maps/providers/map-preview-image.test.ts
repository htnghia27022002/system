import { describe, expect, it } from 'vitest'

import { basemapPreviewImageUrl, osmPreviewImageUrl } from './map-preview-image'

describe('osmPreviewImageUrl', () => {
  it('returns a tile URL for Da Nang', () => {
    const url = osmPreviewImageUrl({ lat: 16.047079, lng: 108.20623 }, 16)
    expect(url).toMatch(/^https:\/\/tile\.openstreetmap\.org\/16\/\d+\/\d+\.png$/)
  })

  it('returns a satellite tile for the satellite basemap', () => {
    const url = basemapPreviewImageUrl('satellite')
    expect(url).toContain('World_Imagery')
  })
})
