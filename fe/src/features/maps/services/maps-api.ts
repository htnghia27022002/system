import { apiClient } from '@/services/api-client'

import type {
  CatalogList,
  Category,
  IngestRun,
  LatestIngestResponse,
  ListLocationsParams,
  ListPlacesParams,
  ListSourcesParams,
  MapSearchResponse,
  LocationRecord,
  LocationWriteInput,
  PaginatedList,
  PlaceDetail,
  PlacePatchInput,
  PlacePin,
  PlaceStatus,
  PlaceWriteInput,
  SourcePatchInput,
  SourceRecord,
  SourceProbeInput,
  SourceProbeResult,
  SourceWriteInput,
} from '../types'

const BASE = '/admin/maps'

export const mapsApi = {
  listPlaces(params: ListPlacesParams = {}): Promise<PaginatedList<PlacePin>> {
    return apiClient
      .get<PaginatedList<PlacePin>>(`${BASE}/places`, {
        params: {
          category: params.category,
          q: params.q,
          adminDivisionId: params.adminDivisionId,
          countryCode: params.countryCode,
          locationId: params.locationId,
          status: params.status,
          manage: params.manage || undefined,
          page: params.page ?? 1,
          limit: params.limit ?? 200,
        },
      })
      .then((r) => r.data)
  },

  getPlace(id: string): Promise<PlaceDetail> {
    return apiClient.get<PlaceDetail>(`${BASE}/places/${id}`).then((r) => r.data)
  },

  updatePlaceStatus(id: string, status: Exclude<PlaceStatus, 'pending'>) {
    return apiClient
      .patch<PlaceDetail>(`${BASE}/places/${id}`, { status })
      .then((r) => r.data)
  },

  createPlace(input: PlaceWriteInput) {
    return apiClient
      .post<PlaceDetail>(`${BASE}/places`, input)
      .then((r) => r.data)
  },

  updatePlace(id: string, input: PlacePatchInput) {
    return apiClient
      .patch<PlaceDetail>(`${BASE}/places/${id}`, input)
      .then((r) => r.data)
  },

  deletePlace(id: string) {
    return apiClient.delete(`${BASE}/places/${id}`).then(() => undefined)
  },

  listLocations(
    params: ListLocationsParams = {},
  ): Promise<PaginatedList<LocationRecord>> {
    return apiClient
      .get<PaginatedList<LocationRecord>>(`${BASE}/locations`, {
        params: {
          q: params.q,
          countryCode: params.countryCode,
          adminDivisionId: params.adminDivisionId,
          page: params.page ?? 1,
          limit: params.limit ?? 50,
        },
      })
      .then((r) => r.data)
  },

  createLocation(input: LocationWriteInput) {
    return apiClient
      .post<LocationRecord>(`${BASE}/locations`, input)
      .then((r) => r.data)
  },

  updateLocation(id: string, input: Partial<LocationWriteInput>) {
    return apiClient
      .patch<LocationRecord>(`${BASE}/locations/${id}`, input)
      .then((r) => r.data)
  },

  deleteLocation(id: string) {
    return apiClient.delete(`${BASE}/locations/${id}`).then(() => undefined)
  },

  listCategories(): Promise<CatalogList<Category>> {
    return apiClient
      .get<CatalogList<Category>>(`${BASE}/categories`)
      .then((r) => r.data)
  },

  search(params: { q: string; countryCode?: string }): Promise<MapSearchResponse> {
    return apiClient
      .get<MapSearchResponse>(`${BASE}/search`, {
        params: {
          q: params.q,
          countryCode: params.countryCode,
        },
      })
      .then((r) => r.data)
  },

  listSources(
    params: ListSourcesParams = {},
  ): Promise<PaginatedList<SourceRecord>> {
    return apiClient
      .get<PaginatedList<SourceRecord>>(`${BASE}/sources`, {
        params: {
          page: params.page ?? 1,
          limit: params.limit ?? 50,
        },
      })
      .then((r) => r.data)
  },

  getSource(id: string): Promise<SourceRecord> {
    return apiClient
      .get<SourceRecord>(`${BASE}/sources/${id}`)
      .then((r) => r.data)
  },

  probeSource(input: SourceProbeInput): Promise<SourceProbeResult> {
    return apiClient
      .post<SourceProbeResult>(`${BASE}/sources/probe`, input)
      .then((r) => r.data)
  },

  createSource(input: SourceWriteInput): Promise<SourceRecord> {
    return apiClient
      .post<SourceRecord>(`${BASE}/sources`, input)
      .then((r) => r.data)
  },

  updateSource(id: string, input: SourcePatchInput): Promise<SourceRecord> {
    return apiClient
      .patch<SourceRecord>(`${BASE}/sources/${id}`, input)
      .then((r) => r.data)
  },

  deleteSource(id: string): Promise<void> {
    return apiClient.delete(`${BASE}/sources/${id}`).then(() => undefined)
  },

  startIngest(): Promise<IngestRun> {
    return apiClient.post<IngestRun>(`${BASE}/ingest`).then((r) => r.data)
  },

  getLatestIngest(): Promise<LatestIngestResponse> {
    return apiClient
      .get<LatestIngestResponse | IngestRun>(`${BASE}/ingest`)
      .then((r) => {
        const data = r.data
        if (data && typeof data === 'object' && 'run' in data) {
          return data as LatestIngestResponse
        }
        if (data && typeof data === 'object' && 'id' in data) {
          return { run: data as IngestRun }
        }
        return { run: null }
      })
  },

  getIngest(id: string): Promise<IngestRun> {
    return apiClient.get<IngestRun>(`${BASE}/ingest/${id}`).then((r) => r.data)
  },
}
