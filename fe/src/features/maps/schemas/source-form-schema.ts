import { z } from 'zod'

export const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH'] as const

export function isValidJsonBody(value: string) {
  const text = value.trim()
  if (!text) return true
  try {
    JSON.parse(text)
    return true
  } catch {
    return false
  }
}

export function formatJsonBody(value: string) {
  const text = value.trim()
  if (!text) return ''
  return JSON.stringify(JSON.parse(text), null, 2)
}

const kvRowSchema = z.object({
  key: z.string(),
  value: z.string(),
})

export const sourceFormSchema = z.object({
  name: z.string().trim().min(1).max(200),
  enabled: z.boolean(),
  httpMethod: z.enum(HTTP_METHODS),
  url: z
    .string()
    .trim()
    .min(1)
    .refine(
      (value) => value.startsWith('http://') || value.startsWith('https://'),
      'URL must be an absolute http:// or https:// address',
    ),
  headers: z.array(kvRowSchema),
  queryParams: z.array(kvRowSchema),
  body: z.string().refine(isValidJsonBody, 'Body must be valid JSON'),
  listPath: z.string(),
  location: z.object({
    name: z.string(),
    locationKey: z.string(),
    countryCode: z.string(),
    adminCode: z.string(),
    street: z.string(),
    postalCode: z.string(),
    formatted: z.string(),
    lat: z.string(),
    lng: z.string(),
  }),
  place: z.object({
    name: z.string(),
    category: z.string(),
    placeKey: z.string(),
    unit: z.string(),
    lat: z.string(),
    lng: z.string(),
  }),
  news: z.object({
    title: z.string(),
    originalUrl: z.string(),
  }),
  details: z.object({
    description: z.string(),
    phone: z.string(),
    hours: z.string(),
    website: z.string(),
    priceRange: z.string(),
  }),
  categoryMap: z.string(),
})

export type SourceFormValues = z.infer<typeof sourceFormSchema>

export const LOCATION_MAPPING_FIELDS = [
  'name',
  'locationKey',
  'countryCode',
  'adminCode',
  'street',
  'postalCode',
  'formatted',
  'lat',
  'lng',
] as const

export const PLACE_MAPPING_FIELDS = [
  'name',
  'category',
  'placeKey',
  'unit',
  'lat',
  'lng',
] as const

export const NEWS_MAPPING_FIELDS = ['title', 'originalUrl'] as const

export const DETAILS_MAPPING_FIELDS = [
  'description',
  'phone',
  'hours',
  'website',
  'priceRange',
] as const

export function emptySourceFormValues(): SourceFormValues {
  return {
    name: '',
    enabled: true,
    httpMethod: 'GET',
    url: '',
    headers: [],
    queryParams: [],
    body: '',
    listPath: '',
    location: {
      name: '',
      locationKey: '',
      countryCode: '',
      adminCode: '',
      street: '',
      postalCode: '',
      formatted: '',
      lat: '',
      lng: '',
    },
    place: {
      name: '',
      category: '',
      placeKey: '',
      unit: '',
      lat: '',
      lng: '',
    },
    news: {
      title: '',
      originalUrl: '',
    },
    details: {
      description: '',
      phone: '',
      hours: '',
      website: '',
      priceRange: '',
    },
    categoryMap: '',
  }
}

export const SAMPLE_FEED_URL = 'http://fe:3000/maps-sample-feed.json'

export function sampleSourceFormValues(): SourceFormValues {
  const empty = emptySourceFormValues()
  return {
    ...empty,
    name: 'Sample Da Nang feed',
    enabled: true,
    httpMethod: 'GET',
    url: SAMPLE_FEED_URL,
    location: {
      ...empty.location,
      name: 'locationName',
      countryCode: 'countryCode',
      lat: 'lat',
      lng: 'lng',
    },
    place: {
      ...empty.place,
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
  }
}

function compactRecord(
  values: Record<string, string>,
): Record<string, string> | undefined {
  const next: Record<string, string> = {}
  for (const [key, value] of Object.entries(values)) {
    const trimmed = value.trim()
    if (trimmed) next[key] = trimmed
  }
  return Object.keys(next).length > 0 ? next : undefined
}

function rowsToRecord(
  rows: { key: string; value: string }[],
): Record<string, string> {
  const next: Record<string, string> = {}
  for (const row of rows) {
    const key = row.key.trim()
    if (key) next[key] = row.value
  }
  return next
}

export function parseCategoryMap(raw: string): Record<string, string> | undefined {
  const trimmed = raw.trim()
  if (!trimmed) return undefined
  const parsed = JSON.parse(trimmed) as unknown
  if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
    throw new Error('categoryMap must be a JSON object')
  }
  const next: Record<string, string> = {}
  for (const [key, value] of Object.entries(parsed)) {
    if (typeof value === 'string' && value.trim()) {
      next[key] = value.trim()
    }
  }
  return Object.keys(next).length > 0 ? next : undefined
}

export function sourceFormToWriteInput(values: SourceFormValues) {
  return {
    name: values.name.trim(),
    enabled: values.enabled,
    httpMethod: values.httpMethod,
    url: values.url.trim(),
    headers: rowsToRecord(values.headers),
    queryParams: rowsToRecord(values.queryParams),
    body: values.body,
    fieldMapping: {
      listPath: values.listPath.trim() || undefined,
      location: compactRecord(values.location),
      place: compactRecord(values.place),
      news: compactRecord(values.news),
      details: compactRecord(values.details),
      categoryMap: parseCategoryMap(values.categoryMap),
    },
  }
}

function recordToRows(record?: Record<string, string> | null) {
  if (!record) return []
  return Object.entries(record).map(([key, value]) => ({ key, value }))
}

function fillGroup<T extends Record<string, string>>(
  empty: T,
  stored?: Partial<T> | null,
): T {
  const next = { ...empty }
  if (!stored) return next
  for (const key of Object.keys(empty) as Array<keyof T>) {
    const value = stored[key]
    if (typeof value === 'string') next[key] = value
  }
  return next
}

export function sourceRecordToFormValues(
  source: {
    name: string
    enabled: boolean
    httpMethod: SourceFormValues['httpMethod']
    url: string
    headers?: Record<string, string>
    queryParams?: Record<string, string>
    body?: string
    fieldMapping?: {
      listPath?: string
      location?: Record<string, string>
      place?: Record<string, string>
      news?: Record<string, string>
      details?: Record<string, string>
      categoryMap?: Record<string, string>
    }
  },
): SourceFormValues {
  const empty = emptySourceFormValues()
  return {
    ...empty,
    name: source.name,
    enabled: source.enabled,
    httpMethod: source.httpMethod,
    url: source.url,
    headers: recordToRows(source.headers),
    queryParams: recordToRows(source.queryParams),
    body: source.body ?? '',
    listPath: source.fieldMapping?.listPath ?? '',
    location: fillGroup(empty.location, source.fieldMapping?.location),
    place: fillGroup(empty.place, source.fieldMapping?.place),
    news: fillGroup(empty.news, source.fieldMapping?.news),
    details: fillGroup(empty.details, source.fieldMapping?.details),
    categoryMap: source.fieldMapping?.categoryMap
      ? JSON.stringify(source.fieldMapping.categoryMap, null, 2)
      : '',
  }
}
