import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios'

import { env } from '@/config/env'
import { useGlobalLoadingStore } from '@/store/global-loading-store'

const TRACKED_METHODS = new Set(['get', 'head'])

function shouldTrackRequest(config: InternalAxiosRequestConfig): boolean {
  if (env.VITE_USE_MOCK_API) {
    return false
  }

  const method = (config.method ?? 'get').toLowerCase()
  return TRACKED_METHODS.has(method)
}

export function setupLoadingInterceptors(client: AxiosInstance): void {
  client.interceptors.request.use((config) => {
    if (shouldTrackRequest(config)) {
      useGlobalLoadingStore.getState().incrementHttp()
    }
    return config
  })

  client.interceptors.response.use(
    (response) => {
      if (shouldTrackRequest(response.config)) {
        useGlobalLoadingStore.getState().decrementHttp()
      }
      return response
    },
    (error) => {
      const config = error.config as InternalAxiosRequestConfig | undefined
      if (config && shouldTrackRequest(config)) {
        useGlobalLoadingStore.getState().decrementHttp()
      }
      return Promise.reject(error)
    },
  )
}
