import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

import { env } from '@/config/env'
import type { AuthTokens } from '@/types/auth'

import { authTokenService } from './auth-token-service'
import { setupLoadingInterceptors } from './setup-loading-interceptors'

type RetryConfig = InternalAxiosRequestConfig & {
  _retry?: boolean
}

let refreshPromise: Promise<string | null> | null = null

// The BE keeps the refresh token in an httpOnly cookie by default
// (AUTH_REFRESH_COOKIE). X-Requested-With is its CSRF guard for cookie-backed
// refresh/logout; withCredentials lets the cookie travel to a cross-origin API.
const CSRF_HEADERS = { 'X-Requested-With': 'XMLHttpRequest' } as const

async function refreshAccessToken(): Promise<string | null> {
  // Only set when the BE runs with the cookie disabled and returned it in JSON.
  const refreshToken = authTokenService.getRefreshToken()

  try {
    const response = await axios.post<AuthTokens>(
      `${env.API_BASE_URL}/auth/refresh`,
      refreshToken ? { refreshToken } : {},
      {
        headers: { 'Content-Type': 'application/json', ...CSRF_HEADERS },
        withCredentials: true,
      },
    )

    authTokenService.setTokens({
      accessToken: response.data.accessToken,
      refreshToken: response.data.refreshToken,
      sessionId: response.data.sessionId,
    })
    return response.data.accessToken
  } catch {
    authTokenService.clearTokens()
    return null
  }
}

export const apiClient = axios.create({
  baseURL: env.API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
    ...CSRF_HEADERS,
  },
})

setupLoadingInterceptors(apiClient)

apiClient.interceptors.request.use((config) => {
  const token = authTokenService.getAccessToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as RetryConfig | undefined

    if (
      error.response?.status !== 401 ||
      !originalRequest ||
      originalRequest._retry
    ) {
      return Promise.reject(error)
    }

    originalRequest._retry = true

    if (!refreshPromise) {
      refreshPromise = refreshAccessToken().finally(() => {
        refreshPromise = null
      })
    }

    const newAccessToken = await refreshPromise
    if (!newAccessToken) {
      return Promise.reject(error)
    }

    originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
    return apiClient(originalRequest)
  },
)
