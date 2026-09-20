'use client'

import { useQuery } from '@tanstack/react-query'

import { addressApi } from '../services/address-api'

export const addressKeys = {
  all: ['address'] as const,
  countries: ['address', 'countries'] as const,
  divisions: (countryCode: string) =>
    ['address', 'divisions', countryCode] as const,
}

export function useAddressCountries() {
  return useQuery({
    queryKey: addressKeys.countries,
    queryFn: () => addressApi.listCountries(),
    meta: { skipNavLoading: true },
  })
}

export function useAddressDivisions(countryCode: string | undefined) {
  return useQuery({
    queryKey: addressKeys.divisions(countryCode ?? ''),
    queryFn: () =>
      addressApi.listDivisions({ countryCode: countryCode as string }),
    enabled: Boolean(countryCode),
    meta: { skipNavLoading: true },
  })
}
