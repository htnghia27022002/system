import { env } from '@/config/env'
import type {
  AuthResolvedUser,
  CreateRoleInput,
  CreateUserInput,
  ListUsersParams,
  ManagedUser,
  PaginatedResponse,
  Permission,
  Role,
  UpdateRoleInput,
  UpdateUserInput,
} from '@/features/access-control/types'
import type { UserRole } from '@/types/auth'
import { mockDelay } from '@/utils/mock-delay'

const LEGACY_AUTH_USERS_KEY = 'mock_auth_users'
const PERMISSIONS_KEY = 'mock_rbac_permissions'
const ROLES_KEY = 'mock_rbac_roles'
const USERS_KEY = 'mock_rbac_users'

const SEED_PERMISSIONS: Permission[] = [
  {
    key: 'dashboard:view',
    name: 'View dashboard',
    group: 'dashboard',
    description: 'Access the admin dashboard overview',
  },
  {
    key: 'users:view',
    name: 'View users',
    group: 'users',
    description: 'List and view user accounts',
  },
  {
    key: 'users:modify',
    name: 'Modify users',
    group: 'users',
    description: 'Create, update, and delete user accounts',
  },
  {
    key: 'roles:view',
    name: 'View roles',
    group: 'roles',
    description: 'List and view roles',
  },
  {
    key: 'roles:modify',
    name: 'Modify roles',
    group: 'roles',
    description: 'Create, update, and delete roles',
  },
  {
    key: 'permissions:view',
    name: 'View permissions',
    group: 'permissions',
    description: 'View the permissions catalog',
  },
]

const ALL_PERMISSION_KEYS = SEED_PERMISSIONS.map((p) => p.key)

const SEED_ROLES: Role[] = [
  {
    id: 'role-admin',
    name: 'Administrator',
    slug: 'admin',
    permissionKeys: [...ALL_PERMISSION_KEYS],
  },
  {
    id: 'role-user',
    name: 'Member',
    slug: 'user',
    permissionKeys: ['dashboard:view'],
  },
]

const SEED_USERS: ManagedUser[] = [
  {
    id: 'admin-user',
    email: 'admin@example.com',
    name: 'Admin User',
    password: 'admin1234',
    roleId: 'role-admin',
    status: 'active',
  },
  {
    id: 'demo-user',
    email: 'demo@example.com',
    name: 'Demo User',
    password: 'password123',
    roleId: 'role-user',
    status: 'active',
  },
]

export class MockAccessControlError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'MockAccessControlError'
    this.status = status
  }
}

function readJson<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return null
    return JSON.parse(raw) as T
  } catch {
    return null
  }
}

function writeJson<T>(key: string, value: T) {
  localStorage.setItem(key, JSON.stringify(value))
}

function migrateLegacyAuthUsers() {
  const legacy = readJson<
    Array<{
      id: string
      email: string
      name: string
      password: string
      role?: UserRole
    }>
  >(LEGACY_AUTH_USERS_KEY)

  if (!legacy?.length) return

  const existing = readJson<ManagedUser[]>(USERS_KEY) ?? []
  const byEmail = new Map(
    existing.map((u) => [u.email.toLowerCase(), u]),
  )

  for (const entry of legacy) {
    const key = entry.email.toLowerCase()
    if (byEmail.has(key)) continue

    const roleId =
      entry.role === 'admin' ||
      entry.email.toLowerCase() === 'admin@example.com'
        ? 'role-admin'
        : 'role-user'

    byEmail.set(key, {
      id: entry.id,
      email: entry.email,
      name: entry.name,
      password: entry.password,
      roleId,
      status: 'active',
    })
  }

  writeJson(USERS_KEY, [...byEmail.values()])
  localStorage.removeItem(LEGACY_AUTH_USERS_KEY)
}

export function ensureRbacSeed() {
  migrateLegacyAuthUsers()

  if (!readJson<Permission[]>(PERMISSIONS_KEY)) {
    writeJson(PERMISSIONS_KEY, SEED_PERMISSIONS)
  }

  const roles = readJson<Role[]>(ROLES_KEY)
  if (!roles?.length) {
    writeJson(ROLES_KEY, SEED_ROLES)
  } else {
    const byId = new Map(roles.map((r) => [r.id, r]))
    for (const seed of SEED_ROLES) {
      const current = byId.get(seed.id)
      if (!current) {
        byId.set(seed.id, seed)
        continue
      }
      if (seed.id === 'role-admin') {
        byId.set(seed.id, {
          ...current,
          permissionKeys: [
            ...new Set([...current.permissionKeys, ...seed.permissionKeys]),
          ],
        })
      }
    }
    writeJson(ROLES_KEY, [...byId.values()])
  }

  const users = readJson<ManagedUser[]>(USERS_KEY)
  if (!users?.length) {
    writeJson(USERS_KEY, SEED_USERS)
    return
  }

  const byEmail = new Map(users.map((u) => [u.email.toLowerCase(), u]))
  for (const seed of SEED_USERS) {
    const key = seed.email.toLowerCase()
    const current = byEmail.get(key)
    if (!current) {
      byEmail.set(key, seed)
      continue
    }
    byEmail.set(key, {
      ...current,
      roleId: current.roleId ?? seed.roleId,
      password:
        key === 'admin@example.com' && current.password === 'admin123'
          ? seed.password
          : current.password,
    })
  }
  writeJson(USERS_KEY, [...byEmail.values()])
}

function readPermissions(): Permission[] {
  ensureRbacSeed()
  return readJson<Permission[]>(PERMISSIONS_KEY) ?? SEED_PERMISSIONS
}

function readRoles(): Role[] {
  ensureRbacSeed()
  return readJson<Role[]>(ROLES_KEY) ?? SEED_ROLES
}

function writeRoles(roles: Role[]) {
  writeJson(ROLES_KEY, roles)
}

function readUsers(): ManagedUser[] {
  ensureRbacSeed()
  return readJson<ManagedUser[]>(USERS_KEY) ?? SEED_USERS
}

function writeUsers(users: ManagedUser[]) {
  writeJson(USERS_KEY, users)
}

export function getRoleById(roleId: string): Role | undefined {
  return readRoles().find((r) => r.id === roleId)
}

export function roleSlugToUserRole(slug: string): UserRole {
  return slug === 'admin' ? 'admin' : 'user'
}

export function getPermissionsForRole(roleId: string): string[] {
  const role = getRoleById(roleId)
  return role?.permissionKeys ?? []
}

export function resolveAuthUserByEmail(
  email: string,
): AuthResolvedUser | undefined {
  const user = readUsers().find(
    (u) => u.email.toLowerCase() === email.toLowerCase(),
  )
  if (!user || user.status !== 'active') return undefined

  const role = getRoleById(user.roleId)
  if (!role) return undefined

  return {
    id: user.id,
    email: user.email,
    name: user.name,
    roleId: user.roleId,
    role: roleSlugToUserRole(role.slug),
    permissions: role.permissionKeys,
  }
}

export function resolveAuthUserByCredentials(
  email: string,
  password: string,
): AuthResolvedUser | undefined {
  const user = readUsers().find(
    (u) => u.email.toLowerCase() === email.toLowerCase(),
  )
  if (!user || user.password !== password || user.status !== 'active') {
    return undefined
  }
  return resolveAuthUserByEmail(user.email)
}

function paginate<T>(
  items: T[],
  page: number,
  pageSize: number,
): PaginatedResponse<T> {
  const start = (page - 1) * pageSize
  return {
    items: items.slice(start, start + pageSize),
    total: items.length,
    page,
    pageSize,
  }
}

function countAdminUsers(users: ManagedUser[], roles: Role[]): number {
  const adminRoleIds = new Set(
    roles.filter((r) => r.slug === 'admin').map((r) => r.id),
  )
  return users.filter((u) => adminRoleIds.has(u.roleId)).length
}

async function withDelay<T>(run: () => T | Promise<T>): Promise<T> {
  await mockDelay(env.VITE_MOCK_API_DELAY_MS)
  return run()
}

export const mockAccessControlApi = {
  async listPermissions(): Promise<Permission[]> {
    return withDelay(() => readPermissions())
  },

  async listRoles(): Promise<Role[]> {
    return withDelay(() => readRoles())
  },

  async getRole(id: string): Promise<Role> {
    return withDelay(() => {
      const role = getRoleById(id)
      if (!role) throw new MockAccessControlError('Role not found', 404)
      return role
    })
  },

  async createRole(input: CreateRoleInput): Promise<Role> {
    return withDelay(() => {
      const roles = readRoles()
      const slugTaken = roles.some(
        (r) => r.slug.toLowerCase() === input.slug.toLowerCase(),
      )
      if (slugTaken) {
        throw new MockAccessControlError('Role slug already exists', 409)
      }

      const catalogKeys = new Set(readPermissions().map((p) => p.key))
      const invalid = input.permissionKeys.filter((k) => !catalogKeys.has(k))
      if (invalid.length) {
        throw new MockAccessControlError('Invalid permission keys', 400)
      }

      const role: Role = {
        id: crypto.randomUUID(),
        name: input.name,
        slug: input.slug,
        permissionKeys: input.permissionKeys,
      }
      writeRoles([...roles, role])
      return role
    })
  },

  async updateRole(id: string, input: UpdateRoleInput): Promise<Role> {
    return withDelay(() => {
      const roles = readRoles()
      const index = roles.findIndex((r) => r.id === id)
      if (index < 0) throw new MockAccessControlError('Role not found', 404)

      const current = roles[index]
      if (input.slug && input.slug !== current.slug) {
        const slugTaken = roles.some(
          (r) =>
            r.id !== id && r.slug.toLowerCase() === input.slug!.toLowerCase(),
        )
        if (slugTaken) {
          throw new MockAccessControlError('Role slug already exists', 409)
        }
      }

      if (input.permissionKeys) {
        const catalogKeys = new Set(readPermissions().map((p) => p.key))
        const invalid = input.permissionKeys.filter((k) => !catalogKeys.has(k))
        if (invalid.length) {
          throw new MockAccessControlError('Invalid permission keys', 400)
        }
      }

      const updated: Role = {
        ...current,
        name: input.name ?? current.name,
        slug: input.slug ?? current.slug,
        permissionKeys: input.permissionKeys ?? current.permissionKeys,
      }
      const next = [...roles]
      next[index] = updated
      writeRoles(next)
      return updated
    })
  },

  async deleteRole(id: string): Promise<void> {
    return withDelay(() => {
      const roles = readRoles()
      const role = roles.find((r) => r.id === id)
      if (!role) throw new MockAccessControlError('Role not found', 404)

      const users = readUsers()
      if (users.some((u) => u.roleId === id)) {
        throw new MockAccessControlError(
          'Cannot delete a role assigned to users',
          409,
        )
      }

      writeRoles(roles.filter((r) => r.id !== id))
    })
  },

  async listUsers(
    params: ListUsersParams = {},
  ): Promise<PaginatedResponse<ManagedUser>> {
    return withDelay(() => {
      const page = params.page ?? 1
      const pageSize = params.pageSize ?? 10
      let users = readUsers()

      if (params.roleId) {
        users = users.filter((u) => u.roleId === params.roleId)
      }

      if (params.search?.trim()) {
        const q = params.search.trim().toLowerCase()
        users = users.filter(
          (u) =>
            u.name.toLowerCase().includes(q) ||
            u.email.toLowerCase().includes(q),
        )
      }

      return paginate(users, page, pageSize)
    })
  },

  async getUser(id: string): Promise<ManagedUser> {
    return withDelay(() => {
      const user = readUsers().find((u) => u.id === id)
      if (!user) throw new MockAccessControlError('User not found', 404)
      return user
    })
  },

  async createUser(input: CreateUserInput): Promise<ManagedUser> {
    return withDelay(() => {
      const users = readUsers()
      const exists = users.some(
        (u) => u.email.toLowerCase() === input.email.toLowerCase(),
      )
      if (exists) {
        throw new MockAccessControlError('Email is already registered', 409)
      }

      if (!getRoleById(input.roleId)) {
        throw new MockAccessControlError('Role not found', 404)
      }

      const user: ManagedUser = {
        id: crypto.randomUUID(),
        email: input.email,
        name: input.name,
        password: input.password,
        roleId: input.roleId,
        status: input.status ?? 'active',
      }
      writeUsers([...users, user])
      return user
    })
  },

  async updateUser(
    id: string,
    input: UpdateUserInput,
    currentSessionUserId?: string,
  ): Promise<ManagedUser> {
    return withDelay(() => {
      const users = readUsers()
      const roles = readRoles()
      const index = users.findIndex((u) => u.id === id)
      if (index < 0) throw new MockAccessControlError('User not found', 404)

      const current = users[index]

      if (input.email && input.email.toLowerCase() !== current.email.toLowerCase()) {
        const emailTaken = users.some(
          (u) =>
            u.id !== id &&
            u.email.toLowerCase() === input.email!.toLowerCase(),
        )
        if (emailTaken) {
          throw new MockAccessControlError('Email is already registered', 409)
        }
      }

      if (input.roleId && !getRoleById(input.roleId)) {
        throw new MockAccessControlError('Role not found', 404)
      }

      const nextRoleId = input.roleId ?? current.roleId
      const adminCount = countAdminUsers(users, roles)
      const isCurrentAdmin = roles.some(
        (r) => r.id === current.roleId && r.slug === 'admin',
      )
      const nextRole = getRoleById(nextRoleId)
      const willBeAdmin = nextRole?.slug === 'admin'

      if (
        currentSessionUserId === id &&
        isCurrentAdmin &&
        !willBeAdmin &&
        adminCount <= 1
      ) {
        throw new MockAccessControlError(
          'Cannot remove admin access from the last administrator',
          400,
        )
      }

      if (isCurrentAdmin && !willBeAdmin && adminCount <= 1) {
        throw new MockAccessControlError(
          'Cannot downgrade the last administrator',
          400,
        )
      }

      const updated: ManagedUser = {
        ...current,
        name: input.name ?? current.name,
        email: input.email ?? current.email,
        password:
          input.password && input.password.length > 0
            ? input.password
            : current.password,
        roleId: nextRoleId,
        status: input.status ?? current.status,
      }
      const next = [...users]
      next[index] = updated
      writeUsers(next)
      return updated
    })
  },

  async deleteUser(
    id: string,
    currentSessionUserId?: string,
  ): Promise<void> {
    return withDelay(() => {
      const users = readUsers()
      const roles = readRoles()
      const user = users.find((u) => u.id === id)
      if (!user) throw new MockAccessControlError('User not found', 404)

      if (currentSessionUserId === id) {
        throw new MockAccessControlError('Cannot delete your own account', 400)
      }

      const role = getRoleById(user.roleId)
      if (role?.slug === 'admin') {
        const adminCount = countAdminUsers(users, roles)
        if (adminCount <= 1) {
          throw new MockAccessControlError(
            'Cannot delete the last administrator',
            400,
          )
        }
      }

      writeUsers(users.filter((u) => u.id !== id))
    })
  },

  async getPermissionsForRole(roleId: string): Promise<string[]> {
    return withDelay(() => getPermissionsForRole(roleId))
  },
}
