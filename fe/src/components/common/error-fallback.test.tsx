import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import i18n from '@/config/i18n'

import { ErrorFallback } from './error-fallback'

describe('ErrorFallback', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
    vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  it('shows a generic message and the digest, never the raw error', () => {
    const error = Object.assign(new Error('pq: password authentication failed'), {
      digest: 'abc123',
    })
    render(<ErrorFallback error={error} reset={() => {}} homeHref="/admin" />)

    expect(screen.getByRole('alert')).toHaveTextContent('Something went wrong')
    expect(screen.getByText('Reference: abc123')).toBeInTheDocument()
    expect(screen.queryByText(/password authentication/)).not.toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Go back' })).toHaveAttribute('href', '/admin')
  })

  it('retries via reset', () => {
    const reset = vi.fn()
    render(<ErrorFallback error={new Error('boom')} reset={reset} />)
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }))
    expect(reset).toHaveBeenCalledOnce()
  })
})
