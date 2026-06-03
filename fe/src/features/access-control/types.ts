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

export type ManagedUser = {
  id: string
  email: string
  name: string
  roleId: string
  status: ManagedUserStatus
  password: string
}

export type ListUsersParams = {
  page?: number
  pageSize?: number
  search?: string
  roleId?: string
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
}

export type UpdateUserInput = {
  email?: string
  name?: string
  password?: string
  roleId?: string
  status?: ManagedUserStatus
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
