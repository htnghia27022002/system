export const P1_PLACE_CATEGORIES = [
  'roomRental',
  'restaurant',
  'hotel',
  'eatery',
] as const

export type P1PlaceCategory = (typeof P1_PLACE_CATEGORIES)[number]

export type PlaceCategoryKey = P1PlaceCategory | 'uncategorized'

export type PlaceStatus = 'pending' | 'active' | 'hidden'

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH'

export type IngestRunStatus = 'queued' | 'running' | 'completed' | 'failed'

export type IngestSourceStatus =
  | 'pending'
  | 'running'
  | 'completed'
  | 'failed'
  | 'skipped'

export type LocationSummary = {
  id: string
  name: string
  locationKey?: string | null
  countryCode?: string | null
  adminDivisionId?: string | null
  adminPath?: string | null
  street?: string
  postalCode?: string
  formatted?: string
  lat?: number | null
  lng?: number | null
}

export type LocationRecord = LocationSummary & {
  createdAt: string
  updatedAt: string
}

export type LocationWriteInput = {
  name: string
  locationKey?: string
  countryCode?: string
  adminDivisionId?: string
  street?: string
  postalCode?: string
  formatted?: string
  lat?: number | null
  lng?: number | null
}

export type AdminDivision = {
  id: string
  countryCode: string
  parentId?: string | null
  level: number
  code: string
  name: string
  nameEn?: string
  fullName?: string
  type: string
  path: string
  lat?: number | null
  lng?: number | null
}

export type Country = {
  id: string
  code: string
  code3: string
  name: string
  nameLocal?: string
}

export type Category = {
  id: string
  key: PlaceCategoryKey | string
  name: string
  nameLocal?: string
  sortOrder: number
}

export type PlacePin = {
  id: string
  name: string
  categoryId: string
  category: PlaceCategoryKey | string
  status: PlaceStatus
  lat?: number | null
  lng?: number | null
  locationId: string
  locationName: string
  unit?: string
  placeKey?: string
}

export type PlaceWriteInput = {
  locationId: string
  name: string
  category: string
  placeKey?: string
  unit?: string
  status?: PlaceStatus
  lat?: number | null
  lng?: number | null
  details?: PlaceDetails
}

export type PlacePatchInput = Partial<PlaceWriteInput>

export type PlaceDetails = {
  description?: string
  phone?: string
  hours?: string
  website?: string
  priceRange?: string
  [key: string]: string | undefined
}

export type NewsItem = {
  id: string
  title: string
  originalUrl: string
  categoryId: string
  category: PlaceCategoryKey | string
  sourceName: string
  sourceId?: string | null
  createdAt: string
  updatedAt: string
}

export type PlaceDetail = {
  id: string
  name: string
  unit?: string
  placeKey?: string
  categoryId: string
  category: PlaceCategoryKey | string
  status: PlaceStatus
  lat?: number | null
  lng?: number | null
  details?: PlaceDetails
  location: LocationSummary
  news: NewsItem[]
  createdAt: string
  updatedAt: string
}

export type FieldMapping = {
  listPath?: string
  location?: {
    name?: string
    locationKey?: string
    countryCode?: string
    adminCode?: string
    street?: string
    postalCode?: string
    formatted?: string
    lat?: string
    lng?: string
  }
  place?: {
    name?: string
    category?: string
    placeKey?: string
    unit?: string
    lat?: string
    lng?: string
  }
  news?: {
    title?: string
    originalUrl?: string
  }
  details?: {
    description?: string
    phone?: string
    hours?: string
    website?: string
    priceRange?: string
  }
  categoryMap?: Record<string, string>
}

export type SourceRecord = {
  id: string
  name: string
  enabled: boolean
  httpMethod: HttpMethod
  url: string
  headers: Record<string, string>
  queryParams: Record<string, string>
  body: string
  fieldMapping: FieldMapping
  createdBy?: string | null
  updatedBy?: string | null
  createdAt: string
  updatedAt: string
}

export type SourceWriteInput = {
  name: string
  enabled?: boolean
  httpMethod: HttpMethod
  url: string
  headers?: Record<string, string>
  queryParams?: Record<string, string>
  body?: string
  fieldMapping: FieldMapping
}

export type SourcePatchInput = Partial<SourceWriteInput>

export type SourceProbeInput = {
  httpMethod: HttpMethod
  url: string
  headers?: Record<string, string>
  queryParams?: Record<string, string>
  body?: string
}

export type SourceProbeResult = {
  status: number
  body: unknown
}

export type IngestRunSource = {
  id: string
  sourceId: string
  sourceName: string
  status: IngestSourceStatus
  errorMessage?: string
  itemsSuccess: number
  itemsError: number
}

export type IngestRun = {
  id: string
  status: IngestRunStatus
  errorMessage?: string
  sourcesTotal: number
  success: number
  error: number
  itemsSuccess: number
  itemsError: number
  sources?: IngestRunSource[]
  startedAt?: string | null
  finishedAt?: string | null
  createdAt: string
}

export type LatestIngestResponse = {
  run: IngestRun | null
}

export type PaginatedList<T> = {
  items: T[]
  page: number
  limit: number
  total: number
  hasMore: boolean
}

export type MapPlaceFilters = {
  category?: P1PlaceCategory
  q?: string
  adminDivisionId?: string
  countryCode?: string
}

export type ListPlacesParams = MapPlaceFilters & {
  page?: number
  limit?: number
  manage?: boolean
  locationId?: string
  status?: PlaceStatus
}

export type ListLocationsParams = {
  q?: string
  countryCode?: string
  adminDivisionId?: string
  page?: number
  limit?: number
}

export type ListSourcesParams = {
  page?: number
  limit?: number
}

export type ListDivisionsParams = {
  countryCode: string
  parentId?: string
  q?: string
}

export type CatalogList<T> = {
  items: T[]
}

export type MapSearchKind = 'place' | 'location' | 'geocode'

export type MapSearchApiHit = {
  id: string
  kind: MapSearchKind
  title: string
  subtitle: string
  lat: number
  lng: number
  placeId?: string
  locationId?: string
}

export type MapSearchResponse = {
  items: MapSearchApiHit[]
}
