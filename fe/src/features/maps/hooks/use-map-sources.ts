'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { mapsApiErrorMessage } from '../api-error'
import { mapsApi } from '../services/maps-api'
import type { SourcePatchInput, SourceWriteInput } from '../types'
import { mapsKeys } from './maps-keys'

export function useMapSources(page = 1) {
  return useQuery({
    queryKey: mapsKeys.sources(page),
    queryFn: () => mapsApi.listSources({ page, limit: 50 }),
    meta: { skipNavLoading: true },
  })
}

export function useMapSourceMutations() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: mapsKeys.all })
  }

  const createSource = useMutation({
    mutationFn: (input: SourceWriteInput) => mapsApi.createSource(input),
    onSuccess: () => {
      invalidate()
      toast.success(t('maps.sources.created'))
    },
    onError: (error) => {
      toast.error(mapsApiErrorMessage(error, t('maps.sources.createFailed')))
    },
  })

  const updateSource = useMutation({
    mutationFn: ({ id, input }: { id: string; input: SourcePatchInput }) =>
      mapsApi.updateSource(id, input),
    onSuccess: () => {
      invalidate()
      toast.success(t('maps.sources.updated'))
    },
    onError: (error) => {
      toast.error(mapsApiErrorMessage(error, t('maps.sources.updateFailed')))
    },
  })

  const deleteSource = useMutation({
    mutationFn: (id: string) => mapsApi.deleteSource(id),
    onSuccess: () => {
      invalidate()
      toast.success(t('maps.sources.deleted'))
    },
    onError: (error) => {
      toast.error(mapsApiErrorMessage(error, t('maps.sources.deleteFailed')))
    },
  })

  return { createSource, updateSource, deleteSource }
}
