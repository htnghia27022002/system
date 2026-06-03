import { env } from '@/config/env'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
} from '@/features/auth/types'
import {
  ensureRbacSeed,
  getRoleById,
  mockAccessControlApi,
  MockAccessControlError,
  resolveAuthUserByCredentials,
  resolveAuthUserByEmail,
} from '@/services/mock/access-control.mock'
import { mockDelay } from '@/utils/mock-delay'
import { createMockAuthTokens } from '@/utils/mock-jwt'

export class MockAuthError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'MockAuthError'
    this.status = status
  }
}

async function withDelay<T>(run: () => T | Promise<T>): Promise<T> {
  await mockDelay(env.VITE_MOCK_API_DELAY_MS)
  return run()
}

function buildAuthResponse(
  resolved: NonNullable<ReturnType<typeof resolveAuthUserByEmail>>,
): AuthResponse {
  const tokens = createMockAuthTokens(
    resolved.id,
    resolved.email,
    resolved.name,
    resolved.role,
    resolved.roleId,
    resolved.permissions,
  )
  return {
    ...tokens,
    user: {
      id: resolved.id,
      email: resolved.email,
      name: resolved.name,
      role: resolved.role,
      roleId: resolved.roleId,
      permissions: resolved.permissions,
    },
  }
}

export const mockAuthApi = {
  async login(payload: LoginRequest): Promise<AuthResponse> {
    return withDelay(() => {
      ensureRbacSeed()
      const resolved = resolveAuthUserByCredentials(
        payload.email,
        payload.password,
      )

      if (!resolved) {
        throw new MockAuthError('Invalid email or password', 401)
      }

      return buildAuthResponse(resolved)
    })
  },

  async register(payload: RegisterRequest): Promise<AuthResponse> {
    return withDelay(async () => {
      ensureRbacSeed()
      const memberRole = getRoleById('role-user')
      if (!memberRole) {
        throw new MockAuthError('Default role is not configured', 500)
      }

      try {
        await mockAccessControlApi.createUser({
          email: payload.email,
          name: payload.name,
          password: payload.password,
          roleId: memberRole.id,
          status: 'active',
        })
      } catch (error) {
        if (error instanceof MockAccessControlError) {
          throw new MockAuthError(error.message, error.status)
        }
        throw error
      }

      const resolved = resolveAuthUserByEmail(payload.email)
      if (!resolved) {
        throw new MockAuthError('Registration failed', 500)
      }

      return buildAuthResponse(resolved)
    })
  },
}
