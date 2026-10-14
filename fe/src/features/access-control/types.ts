export type { PaginatedResponse } from '@/types/pagination'

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

/** `invited`: created by an admin invitation, waiting for the user to set a password. */
export type ManagedUserStatus = 'active' | 'inactive' | 'invited'

export type SocialLink = {
  label?: string
  url: string
}

export type ManagedUser = {
  id: string
  email: string
  name: string
  roleId: string
  superAdmin?: boolean
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
  /** False until the user confirms the email address. */
  emailVerified?: boolean
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
  superAdmin?: boolean
}

export type InviteUserInput = {
  email: string
  name: string
  roleId: string
}

export type UpdateUserInput = {
  email?: string
  name?: string
  password?: string
  roleId?: string
  /** The BE only accepts active / inactive; `invited` ends by accepting the invitation. */
  status?: Exclude<ManagedUserStatus, 'invited'>
  phone?: string
  general?: string
  birthday?: string | null
  address?: string
  socialLinks?: SocialLink[]
  superAdmin?: boolean
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
  superAdmin?: boolean
}
