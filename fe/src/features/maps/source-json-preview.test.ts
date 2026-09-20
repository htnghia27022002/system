import { describe, expect, it } from 'vitest'

import {
  autoMapFields,
  flattenFieldPaths,
  resolveList,
  suggestListPaths,
} from './source-json-preview'

describe('source-json-preview', () => {
  it('reads a nested list path', () => {
    const root = { data: { items: [{ name: 'A' }] } }
    expect(resolveList(root, 'data.items')).toEqual([{ name: 'A' }])
    expect(suggestListPaths(root)).toEqual(['data.items'])
  })

  it('flattens item fields and auto-maps sample names', () => {
    const item = {
      locationName: 'Cho Han',
      countryCode: 'VN',
      placeName: 'Cafe',
      category: 'eatery',
      lat: 16.04,
      lng: 108.2,
      title: 'News',
      originalUrl: 'https://example.com/a',
      price: '4tr',
    }
    const fields = flattenFieldPaths(item)
    const mapped = autoMapFields(fields)
    expect(mapped['location.name']).toBe('locationName')
    expect(mapped['place.name']).toBe('placeName')
    expect(mapped['details.priceRange']).toBe('price')
    expect(mapped['news.originalUrl']).toBe('originalUrl')
  })
})
