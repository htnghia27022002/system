import { apiClient } from '@/services/api-client'

import type { AdminDivision, CatalogList, Country } from '../types'

export const addressApi = {
  listCountries(): Promise<CatalogList<Country>> {
    return apiClient
      .get<CatalogList<Country>>('/address/countries')
      .then((response) => response.data)
  },

  listDivisions(params: {
    countryCode: string
    parentId?: string
    q?: string
  }): Promise<CatalogList<AdminDivision>> {
    return apiClient
      .get<CatalogList<AdminDivision>>('/address/divisions', {
        params: {
          countryCode: params.countryCode,
          parentId: params.parentId,
          q: params.q,
        },
      })
      .then((response) => response.data)
  },
}
