'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { mapsApiErrorMessage } from '../api-error'
import { mapsApi } from '../services/maps-api'
import type { LocationWriteInput } from '../types'
import { mapsKeys } from './maps-keys'

export function useMapLocations() {
  return useQuery({
    queryKey: mapsKeys.locations(),
    queryFn: () => mapsApi.listLocations({ page: 1, limit: 100 }),
    meta: { skipNavLoading: true },
  })
}

export function useMapLocationMutations() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: mapsKeys.all })
  }

  return {
    createLocation: useMutation({
      mutationFn: (input: LocationWriteInput) => mapsApi.createLocation(input),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.locations.created'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.locations.createFailed')))
      },
    }),
    updateLocation: useMutation({
      mutationFn: ({
        id,
        input,
      }: {
        id: string
        input: Partial<LocationWriteInput>
      }) => mapsApi.updateLocation(id, input),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.locations.updated'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.locations.updateFailed')))
      },
    }),
    deleteLocation: useMutation({
      mutationFn: (id: string) => mapsApi.deleteLocation(id),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.locations.deleted'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.locations.deleteFailed')))
      },
    }),
  }
}
