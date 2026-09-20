import { describe, expect, it } from 'vitest'

import {
  hasActivePlaceFilters,
  isP1FilterCategory,
  P1_PLACE_CATEGORIES,
  toMapPlaceFilters,
  toPlacesCategoryParam,
} from './category-filters'

describe('category filter helpers', () => {
  it('recognizes P1 filter categories and excludes uncategorized', () => {
    expect(P1_PLACE_CATEGORIES).toEqual([
      'roomRental',
      'restaurant',
      'hotel',
      'eatery',
    ])
    expect(isP1FilterCategory('roomRental')).toBe(true)
    expect(isP1FilterCategory('uncategorized')).toBe(false)
    expect(isP1FilterCategory('unknown')).toBe(false)
  })

  it('maps UI category values to the places query param', () => {
    expect(toPlacesCategoryParam('all')).toBeUndefined()
    expect(toPlacesCategoryParam('')).toBeUndefined()
    expect(toPlacesCategoryParam('hotel')).toBe('hotel')
    expect(toPlacesCategoryParam('uncategorized')).toBeUndefined()
  })

  it('builds place filters and detects an active filter state', () => {
    expect(hasActivePlaceFilters({})).toBe(false)
    expect(
      hasActivePlaceFilters(
        toMapPlaceFilters({
          category: 'eatery',
          q: '  pho  ',
          countryCode: 'VN',
        }),
      ),
    ).toBe(true)
    expect(toMapPlaceFilters({ q: '   ' })).toEqual({
      category: undefined,
      q: undefined,
      adminDivisionId: undefined,
      countryCode: undefined,
    })
  })
})
