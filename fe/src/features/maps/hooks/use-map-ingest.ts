'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { isIngestConflict, mapsApiErrorMessage } from '../api-error'
import { mapsApi } from '../services/maps-api'
import { mapsKeys } from './maps-keys'

const IN_PROGRESS = new Set(['queued', 'running'])

export function useLatestIngest() {
  const queryClient = useQueryClient()
  const previousStatus = useRef<string | undefined>(undefined)

  const query = useQuery({
    queryKey: mapsKeys.ingest,
    queryFn: () => mapsApi.getLatestIngest(),
    refetchInterval: (current) => {
      const status = current.state.data?.run?.status
      return status && IN_PROGRESS.has(status) ? 2000 : false
    },
    meta: { skipNavLoading: true },
  })

  const status = query.data?.run?.status

  useEffect(() => {
    const wasInProgress =
      previousStatus.current !== undefined &&
      IN_PROGRESS.has(previousStatus.current)
    const finished = status === 'completed' || status === 'failed'
    if (wasInProgress && finished) {
      void queryClient.invalidateQueries({ queryKey: mapsKeys.places({}) })
      void queryClient.invalidateQueries({
        predicate: (entry) =>
          entry.queryKey[0] === 'admin' &&
          entry.queryKey[1] === 'maps' &&
          entry.queryKey[2] === 'places',
      })
    }
    previousStatus.current = status
  }, [queryClient, status])

  return query
}

export function useStartIngest() {
  const queryClient = useQueryClient()
  const { t } = useTranslation('admin')

  return useMutation({
    mutationFn: () => mapsApi.startIngest(),
    onSuccess: (run) => {
      queryClient.setQueryData(mapsKeys.ingest, { run })
      toast.success(t('maps.ingest.started'))
      if (run.status === 'completed' || run.status === 'failed') {
        void queryClient.invalidateQueries({
          predicate: (entry) =>
            entry.queryKey[0] === 'admin' &&
            entry.queryKey[1] === 'maps' &&
            entry.queryKey[2] === 'places',
        })
      }
    },
    onError: (error) => {
      toast.error(
        mapsApiErrorMessage(
          error,
          isIngestConflict(error)
            ? t('maps.ingest.conflict')
            : t('maps.ingest.startFailed'),
        ),
      )
    },
  })
}

export function isIngestInProgress(status?: string | null): boolean {
  return Boolean(status && IN_PROGRESS.has(status))
}
