import { AxiosError, AxiosHeaders } from 'axios'
import { describe, expect, it } from 'vitest'

import { apiErrorCode, apiErrorMessage, apiErrorStatus } from './api-error'

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
