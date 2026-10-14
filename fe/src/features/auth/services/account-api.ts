import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { loadAccountMock } from '@/services/mock'

export type Invitation = {
  email: string
  name: string
}

/**
 * Email-driven account flows. Tokens come from links in emails and are always
 * sent in the JSON body, never in API URLs.
 */
export const accountApi = {
  async forgotPassword(email: string): Promise<void> {
    if (env.USE_MOCK_API) return (await loadAccountMock()).forgotPassword(email)
    await apiClient.post('/auth/password/forgot', { email })
  },

  async resetPassword(token: string, password: string): Promise<void> {
    if (env.USE_MOCK_API) return (await loadAccountMock()).resetPassword(token, password)
    await apiClient.post('/auth/password/reset', { token, password })
  },

  async verifyEmail(token: string): Promise<void> {
    if (env.USE_MOCK_API) return (await loadAccountMock()).verifyEmail(token)
    await apiClient.post('/auth/email/verify', { token })
  },

  async inspectInvite(token: string): Promise<Invitation> {
    if (env.USE_MOCK_API) return (await loadAccountMock()).inspectInvite(token)
    const { data } = await apiClient.post<Invitation>('/auth/invite/inspect', { token })
    return data
  },

  async acceptInvite(token: string, password: string, name?: string): Promise<void> {
    if (env.USE_MOCK_API) return (await loadAccountMock()).acceptInvite(token, password, name)
    await apiClient.post('/auth/invite/accept', { token, password, ...(name ? { name } : {}) })
  },
}
