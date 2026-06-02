import { env } from '@/config/env'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
} from '@/features/auth/types'
import type { UserRole } from '@/types/auth'
import { mockDelay } from '@/utils/mock-delay'
import { createMockAuthTokens } from '@/utils/mock-jwt'

const MOCK_USERS_KEY = 'mock_auth_users'

type MockUser = {
  id: string
  email: string
  name: string
  password: string
  role: UserRole
}

function readUsers(): MockUser[] {
  try {
    const raw = localStorage.getItem(MOCK_USERS_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as Array<MockUser & { role?: UserRole }>
    return parsed.map((entry) => ({
      ...entry,
      role:
        entry.role ??
        (entry.email.toLowerCase() === 'admin@example.com' ? 'admin' : 'user'),
    }))
  } catch {
    return []
  }
}

function writeUsers(users: MockUser[]) {
  localStorage.setItem(MOCK_USERS_KEY, JSON.stringify(users))
}

const SEED_USERS: MockUser[] = [
  {
    id: 'admin-user',
    email: 'admin@example.com',
    name: 'Admin User',
    password: 'admin1234',
    role: 'admin',
  },
  {
    id: 'demo-user',
    email: 'demo@example.com',
    name: 'Demo User',
    password: 'password123',
    role: 'user',
  },
]

function ensureSeedUsers() {
  const existing = readUsers()
  const byEmail = new Map(
    existing.map((user) => [user.email.toLowerCase(), user]),
  )

  for (const seed of SEED_USERS) {
    const key = seed.email.toLowerCase()
    const current = byEmail.get(key)
    if (!current) {
      byEmail.set(key, seed)
      continue
    }

    // Upgrade legacy seeds (e.g. missing role or old short admin password).
    byEmail.set(key, {
      ...current,
      role: current.role ?? seed.role,
      password:
        key === 'admin@example.com' && current.password === 'admin123'
          ? seed.password
          : current.password,
    })
  }

  writeUsers([...byEmail.values()])
}

function buildAuthResponse(user: MockUser): AuthResponse {
  const tokens = createMockAuthTokens(user.id, user.email, user.name, user.role)
  return {
    ...tokens,
    user: {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role,
    },
  }
}

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

export const mockAuthApi = {
  async login(payload: LoginRequest): Promise<AuthResponse> {
    return withDelay(() => {
      ensureSeedUsers()
      const user = readUsers().find(
        (entry) => entry.email.toLowerCase() === payload.email.toLowerCase(),
      )

      if (!user || user.password !== payload.password) {
        throw new MockAuthError('Invalid email or password', 401)
      }

      return buildAuthResponse(user)
    })
  },

  async register(payload: RegisterRequest): Promise<AuthResponse> {
    return withDelay(() => {
      ensureSeedUsers()
      const users = readUsers()
      const exists = users.some(
        (entry) => entry.email.toLowerCase() === payload.email.toLowerCase(),
      )

      if (exists) {
        throw new MockAuthError('Email is already registered', 409)
      }

      const user: MockUser = {
        id: crypto.randomUUID(),
        email: payload.email,
        name: payload.name,
        password: payload.password,
        role: 'user',
      }

      writeUsers([...users, user])
      return buildAuthResponse(user)
    })
  },
}
