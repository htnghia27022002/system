import axios from 'axios'

export function mapsApiErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as
      | { message?: string; error?: string }
      | undefined
    return data?.error || data?.message || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function isIngestConflict(error: unknown): boolean {
  return axios.isAxiosError(error) && error.response?.status === 409
}
