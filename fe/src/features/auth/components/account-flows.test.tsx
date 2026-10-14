import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { AxiosError, AxiosHeaders } from 'axios'
import type { ReactNode } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import i18n from '@/config/i18n'

const { api, nav } = vi.hoisted(() => ({
  api: {
    forgotPassword: vi.fn(),
    resetPassword: vi.fn(),
    verifyEmail: vi.fn(),
    inspectInvite: vi.fn(),
    acceptInvite: vi.fn(),
  },
  nav: { token: 'raw-token' as string | null, replace: vi.fn() },
}))

vi.mock('../services/account-api', () => ({ accountApi: api }))
vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: nav.replace, push: vi.fn() }),
  useSearchParams: () => new URLSearchParams(nav.token ? { token: nav.token } : {}),
}))

import { AcceptInvitePage } from './accept-invite-page'
import { ForgotPasswordPage } from './forgot-password-page'
import { ResetPasswordPage } from './reset-password-page'
import { VerifyEmailPage } from './verify-email-page'

function wrap(ui: ReactNode) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>)
}

function badRequest() {
  return new AxiosError('bad', 'ERR', undefined, undefined, {
    status: 400,
    statusText: '',
    headers: {},
    config: { headers: new AxiosHeaders() },
    data: { error: 'invalid' },
  })
}

describe('account flow pages', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('en')
    Object.values(api).forEach((fn) => fn.mockReset())
    nav.token = 'raw-token'
    nav.replace.mockReset()
  })

  it('forgot password shows the same confirmation for any address', async () => {
    api.forgotPassword.mockResolvedValue(undefined)
    wrap(<ForgotPasswordPage />)
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'a@example.com' } })
    fireEvent.click(screen.getByRole('button', { name: 'Send reset link' }))
    expect(await screen.findByRole('status')).toHaveTextContent('If an account exists for a@example.com')
    expect(api.forgotPassword).toHaveBeenCalledWith('a@example.com')
  })

  it('reset password validates, submits the token and returns to sign in', async () => {
    api.resetPassword.mockResolvedValue(undefined)
    wrap(<ResetPasswordPage />)
    fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'new-password-1' } })
    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'different-123' } })
    fireEvent.click(screen.getByRole('button', { name: 'Update password' }))
    expect(await screen.findByText('Passwords do not match')).toBeInTheDocument()
    expect(api.resetPassword).not.toHaveBeenCalled()

    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'new-password-1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Update password' }))
    await waitFor(() => expect(nav.replace).toHaveBeenCalledWith('/login'))
    expect(api.resetPassword).toHaveBeenCalledWith('raw-token', 'new-password-1')
  })

  it('reset password shows the invalid-link notice without a token or on 400', async () => {
    nav.token = null
    const { unmount } = wrap(<ResetPasswordPage />)
    expect(screen.getByRole('alert')).toHaveTextContent('invalid or has expired')
    unmount()

    nav.token = 'used-token'
    api.resetPassword.mockRejectedValue(badRequest())
    wrap(<ResetPasswordPage />)
    fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'new-password-1' } })
    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'new-password-1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Update password' }))
    expect(await screen.findByRole('link', { name: 'Request a new link' })).toHaveAttribute('href', '/forgot-password')
  })

  it('verify email confirms the token once', async () => {
    api.verifyEmail.mockResolvedValue(undefined)
    wrap(<VerifyEmailPage />)
    expect(await screen.findByText('Your email address is verified.')).toBeInTheDocument()
    expect(api.verifyEmail).toHaveBeenCalledTimes(1)
    expect(api.verifyEmail).toHaveBeenCalledWith('raw-token')
  })

  it('accept invite pre-fills the invitation and activates the account', async () => {
    api.inspectInvite.mockResolvedValue({ email: 'invitee@example.com', name: 'Invited Person' })
    api.acceptInvite.mockResolvedValue(undefined)
    wrap(<AcceptInvitePage />)
    expect(await screen.findByDisplayValue('Invited Person')).toBeInTheDocument()
    expect(screen.getByDisplayValue('invitee@example.com')).toBeDisabled()

    fireEvent.change(screen.getByLabelText('New password'), { target: { value: 'chosen-pass-1' } })
    fireEvent.change(screen.getByLabelText('Confirm password'), { target: { value: 'chosen-pass-1' } })
    fireEvent.click(screen.getByRole('button', { name: 'Activate account' }))
    await waitFor(() => expect(nav.replace).toHaveBeenCalledWith('/login'))
    expect(api.acceptInvite).toHaveBeenCalledWith('raw-token', 'chosen-pass-1', 'Invited Person')
  })

  it('accept invite rejects an invalid link', async () => {
    api.inspectInvite.mockRejectedValue(badRequest())
    wrap(<AcceptInvitePage />)
    expect(await screen.findByRole('alert')).toHaveTextContent('invalid or has expired')
  })
})
