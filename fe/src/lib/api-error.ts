import axios from 'axios'

type ApiErrorBody = { error?: string; message?: string; code?: string }

function body(error: unknown): ApiErrorBody | undefined {
  if (axios.isAxiosError(error)) {
    return error.response?.data as ApiErrorBody | undefined
  }
  return undefined
}

/** HTTP status of a failed API call, or undefined for non-HTTP errors. */
export function apiErrorStatus(error: unknown): number | undefined {
  if (axios.isAxiosError(error)) return error.response?.status
  if (error && typeof error === 'object' && 'status' in error) {
    const status = (error as { status?: unknown }).status
    return typeof status === 'number' ? status : undefined
  }
  return undefined
}

/**
 * Machine-readable error code from the BE (`{"code": "email_not_verified"}`).
 * Branch on this — never on the English message.
 */
export function apiErrorCode(error: unknown): string | undefined {
  const code = body(error)?.code
  if (code) return code
  if (error && typeof error === 'object' && 'code' in error) {
    const value = (error as { code?: unknown }).code
    return typeof value === 'string' ? value : undefined
  }
  return undefined
}

/** Server message when present, else `fallback`. */
export function apiErrorMessage(error: unknown, fallback: string): string {
  const data = body(error)
  return data?.message || data?.error || fallback
}
