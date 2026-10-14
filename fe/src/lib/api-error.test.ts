import { AxiosError, AxiosHeaders } from 'axios'
import { describe, expect, it } from 'vitest'

import { apiErrorCode, apiErrorMessage, apiErrorStatus, translateApiError } from './api-error'

function axiosError(status: number, data: unknown) {
  return new AxiosError('fail', 'ERR', undefined, undefined, {
    status,
    statusText: '',
    headers: {},
    config: { headers: new AxiosHeaders() },
    data,
  })
}

describe('api-error', () => {
  it('reads status, code and message from an API error', () => {
    const error = axiosError(403, { error: 'verify first', code: 'email_not_verified' })
    expect(apiErrorStatus(error)).toBe(403)
    expect(apiErrorCode(error)).toBe('email_not_verified')
    expect(apiErrorMessage(error, 'fallback')).toBe('verify first')
  })

  it('falls back for non-HTTP errors', () => {
    expect(apiErrorStatus(new Error('x'))).toBeUndefined()
    expect(apiErrorCode(new Error('x'))).toBeUndefined()
    expect(apiErrorMessage(new Error('x'), 'fallback')).toBe('fallback')
  })

  it('reads status/code carried by mock errors', () => {
    const mock = Object.assign(new Error('locked'), { status: 429, code: 'account_locked' })
    expect(apiErrorStatus(mock)).toBe(429)
    expect(apiErrorCode(mock)).toBe('account_locked')
  })
})

describe('translateApiError', () => {
  const t = (key: string, options?: { defaultValue?: string }) =>
    key === 'errors.forbidden' ? 'Bạn không có quyền' : (options?.defaultValue ?? key)

  it('translates known codes and falls back to the server message', () => {
    expect(translateApiError(axiosError(403, { error: 'forbidden', code: 'forbidden' }), t, 'x')).toBe('Bạn không có quyền')
    expect(translateApiError(axiosError(409, { error: 'email already exists', code: 'conflict' }), t, 'x')).toBe('email already exists')
    expect(translateApiError(new Error('boom'), t, 'fallback')).toBe('fallback')
  })
})
