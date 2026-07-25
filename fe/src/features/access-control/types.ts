export type Permission = {
  key: string
  name: string
  group: string
  description: string
}

export type Role = {
  id: string
  name: string
  slug: string
  permissionKeys: string[]
}

export type ManagedUserStatus = 'active' | 'inactive'

export type SocialLink = {
  label?: string
  url: string
}

export type ManagedUser = {
  id: string
  email: string
  name: string
  roleId: string
  status: ManagedUserStatus
  password: string
  phone?: string
  avatarUrl?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
  /** Linked OAuth provider ids, e.g. `google`. */
  oauthProviders?: string[]
  createdAt?: string
}

export type ListUsersParams = {
  page?: number
  pageSize?: number
  search?: string
  roleId?: string
  id?: string
  status?: ManagedUserStatus
}

export type ListRolesParams = {
  page?: number
  pageSize?: number
  search?: string
  id?: string
  permissionKey?: string
}

export type ListPermissionsParams = {
  page?: number
  pageSize?: number
  search?: string
  group?: string
}

export type PaginatedResponse<T> = {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export type CreateUserInput = {
  email: string
  name: string
  password: string
  roleId: string
  status?: ManagedUserStatus
  phone?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
}

export type UpdateUserInput = {
  email?: string
  name?: string
  password?: string
  roleId?: string
  status?: ManagedUserStatus
  phone?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
}

export type CreateRoleInput = {
  name: string
  slug: string
  permissionKeys: string[]
}

export type UpdateRoleInput = {
  name?: string
  slug?: string
  permissionKeys?: string[]
}

export type AuthResolvedUser = {
  id: string
  email: string
  name: string
  roleId: string
  role: 'admin' | 'user'
  permissions: string[]
}
