import { describe, expect, it } from 'vitest'

import {
  SAMPLE_FEED_URL,
  formatJsonBody,
  isValidJsonBody,
  sampleSourceFormValues,
  sourceFormToWriteInput,
} from './source-form-schema'

describe('sampleSourceFormValues', () => {
  it('maps the Da Nang sample feed onto Location, Place, and news paths', () => {
    const input = sourceFormToWriteInput(sampleSourceFormValues())

    expect(input.name).toBe('Sample Da Nang feed')
    expect(input.enabled).toBe(true)
    expect(input.url).toBe(SAMPLE_FEED_URL)
    expect(input.fieldMapping).toEqual({
      location: {
        name: 'locationName',
        countryCode: 'countryCode',
        lat: 'lat',
        lng: 'lng',
      },
      place: {
        name: 'placeName',
        category: 'category',
        placeKey: 'placeKey',
        lat: 'lat',
        lng: 'lng',
      },
      news: {
        title: 'title',
        originalUrl: 'originalUrl',
      },
    })
  })
})

describe('source body JSON', () => {
  it('accepts empty or valid JSON and rejects invalid JSON', () => {
    expect(isValidJsonBody('')).toBe(true)
    expect(isValidJsonBody('{"limit":10}')).toBe(true)
    expect(isValidJsonBody('{')).toBe(false)
  })

  it('pretty-prints a valid JSON body', () => {
    expect(formatJsonBody('{"limit":10}')).toBe('{\n  "limit": 10\n}')
  })
})
