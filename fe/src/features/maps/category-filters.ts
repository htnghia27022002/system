import {
  P1_PLACE_CATEGORIES,
  type MapPlaceFilters,
  type P1PlaceCategory,
} from './types'

export { P1_PLACE_CATEGORIES }
export type { P1PlaceCategory }

export function isP1FilterCategory(
  value: string | undefined | null,
): value is P1PlaceCategory {
  return (
    typeof value === 'string' &&
    (P1_PLACE_CATEGORIES as readonly string[]).includes(value)
  )
}

export function toPlacesCategoryParam(
  value: string | undefined | null,
): P1PlaceCategory | undefined {
  if (!value || value === 'all') return undefined
  return isP1FilterCategory(value) ? value : undefined
}

export function hasActivePlaceFilters(filters: MapPlaceFilters): boolean {
  return Boolean(
    filters.category ||
      filters.q?.trim() ||
      filters.adminDivisionId ||
      filters.countryCode,
  )
}

export function toMapPlaceFilters(input: {
  category?: string
  q?: string
  adminDivisionId?: string
  countryCode?: string
}): MapPlaceFilters {
  const q = input.q?.trim()
  return {
    category: toPlacesCategoryParam(input.category),
    q: q || undefined,
    adminDivisionId: input.adminDivisionId || undefined,
    countryCode: input.countryCode || undefined,
  }
}
