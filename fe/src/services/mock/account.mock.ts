import { env } from '@/config/env'
import { mockDelay } from '@/utils/mock-delay'

import { MockAuthError } from './mock-errors'

/** Any token works in mock mode except this one, which simulates an expired link. */
export const MOCK_INVALID_TOKEN = 'expired'

async function check(token: string) {
  await mockDelay(env.MOCK_API_DELAY_MS)
  if (!token || token === MOCK_INVALID_TOKEN) {
    throw new MockAuthError('This link is invalid or has expired', 400)
  }
}

export const mockAccountApi = {
  async forgotPassword(_email: string): Promise<void> {
    await mockDelay(env.MOCK_API_DELAY_MS)
  },
  async resetPassword(token: string, _password: string): Promise<void> {
    await check(token)
  },
  async verifyEmail(token: string): Promise<void> {
    await check(token)
  },
  async inspectInvite(token: string) {
    await check(token)
    return { email: 'invitee@example.com', name: 'Invited User' }
  },
  async acceptInvite(token: string, _password: string, _name?: string): Promise<void> {
    await check(token)
  },
  async resendVerification(): Promise<void> {
    await mockDelay(env.MOCK_API_DELAY_MS)
  },
}
