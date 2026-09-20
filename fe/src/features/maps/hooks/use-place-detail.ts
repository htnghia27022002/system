'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { mapsApiErrorMessage } from '../api-error'
import { mapsApi } from '../services/maps-api'
import type { PlaceStatus } from '../types'
import { mapsKeys } from './maps-keys'

export function usePlaceDetail(id: string | null) {
  return useQuery({
    queryKey: mapsKeys.place(id ?? ''),
    queryFn: () => mapsApi.getPlace(id as string),
    enabled: Boolean(id),
    meta: { skipNavLoading: true },
  })
}

export function usePlaceStatusMutation() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  return useMutation({
    mutationFn: ({
      id,
      status,
    }: {
      id: string
      status: Exclude<PlaceStatus, 'pending'>
    }) => mapsApi.updatePlaceStatus(id, status),
    onSuccess: (place) => {
      void queryClient.invalidateQueries({ queryKey: mapsKeys.all })
      queryClient.setQueryData(mapsKeys.place(place.id), place)
      toast.success(t('maps.detail.statusUpdated'))
    },
    onError: (error) => {
      toast.error(mapsApiErrorMessage(error, t('maps.detail.statusFailed')))
    },
  })
}
