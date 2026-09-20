'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { mapsApiErrorMessage } from '../api-error'
import { mapsApi } from '../services/maps-api'
import type { PlacePatchInput, PlaceWriteInput } from '../types'
import { mapsKeys } from './maps-keys'

export function useManagedPlaces(locationId?: string) {
  return useQuery({
    queryKey: mapsKeys.managedPlaces(locationId),
    queryFn: () =>
      mapsApi.listPlaces({
        manage: true,
        locationId,
        page: 1,
        limit: 100,
      }),
    meta: { skipNavLoading: true },
  })
}

export function usePlaceMutations() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: mapsKeys.all })
  }

  return {
    createPlace: useMutation({
      mutationFn: (input: PlaceWriteInput) => mapsApi.createPlace(input),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.places.created'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.places.createFailed')))
      },
    }),
    updatePlace: useMutation({
      mutationFn: ({ id, input }: { id: string; input: PlacePatchInput }) =>
        mapsApi.updatePlace(id, input),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.places.updated'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.places.updateFailed')))
      },
    }),
    deletePlace: useMutation({
      mutationFn: (id: string) => mapsApi.deletePlace(id),
      onSuccess: () => {
        invalidate()
        toast.success(t('maps.places.deleted'))
      },
      onError: (error) => {
        toast.error(mapsApiErrorMessage(error, t('maps.places.deleteFailed')))
      },
    }),
  }
}
