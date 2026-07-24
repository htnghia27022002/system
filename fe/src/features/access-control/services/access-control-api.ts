import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { normalizePaginatedResponse } from '@/lib/normalize-paginated-response'
import {
  mockAccessControlApi,
  MockAccessControlError,
} from '@/services/mock/access-control.mock'

import type {
  CreateRoleInput,
  CreateUserInput,
  ListPermissionsParams,
  ListRolesParams,
  ListUsersParams,
  ManagedUser,
  PaginatedResponse,
  Permission,
  Role,
  UpdateRoleInput,
  UpdateUserInput,
} from '../types'

export { MockAccessControlError }

export const accessControlApi = {
  listPermissions(
    params: ListPermissionsParams = {},
  ): Promise<PaginatedResponse<Permission>> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.listPermissions(params)
    }
    return apiClient
      .get<PaginatedResponse<Permission>>('/admin/permissions', { params })
      .then((r) => normalizePaginatedResponse<Permission>(r.data, params.pageSize))
  },

  listAllPermissions(): Promise<Permission[]> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi
        .listPermissions({ page: 1, pageSize: 1000 })
        .then((r) => r.items)
    }
    return apiClient
      .get<Permission[]>('/admin/permissions/all')
      .then((r) => r.data)
  },

  listRoles(params: ListRolesParams = {}): Promise<PaginatedResponse<Role>> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.listRoles(params)
    }
    return apiClient
      .get<PaginatedResponse<Role>>('/admin/roles', {
        params,
        skipNavLoading: true,
      })
      .then((r) => normalizePaginatedResponse<Role>(r.data, params.pageSize))
  },

  listAllRoles(): Promise<Role[]> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi
        .listRoles({ page: 1, pageSize: 1000 })
        .then((r) => r.items)
    }
    return apiClient
      .get<Role[]>('/admin/roles/all', { skipNavLoading: true })
      .then((r) => r.data)
  },

  getRole(id: string): Promise<Role> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.getRole(id)
    }
    return apiClient.get<Role>(`/admin/roles/${id}`).then((r) => r.data)
  },

  createRole(input: CreateRoleInput): Promise<Role> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.createRole(input)
    }
    return apiClient.post<Role>('/admin/roles', input).then((r) => r.data)
  },

  updateRole(id: string, input: UpdateRoleInput): Promise<Role> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.updateRole(id, input)
    }
    return apiClient
      .patch<Role>(`/admin/roles/${id}`, input)
      .then((r) => r.data)
  },

  deleteRole(id: string): Promise<void> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.deleteRole(id)
    }
    return apiClient.delete(`/admin/roles/${id}`).then(() => undefined)
  },

  listUsers(params: ListUsersParams = {}): Promise<PaginatedResponse<ManagedUser>> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.listUsers(params)
    }
    return apiClient
      .get<PaginatedResponse<ManagedUser>>('/admin/users', {
        params,
        skipNavLoading: true,
      })
      .then((r) => normalizePaginatedResponse<ManagedUser>(r.data, params.pageSize))
  },

  getUser(id: string): Promise<ManagedUser> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.getUser(id)
    }
    return apiClient.get<ManagedUser>(`/admin/users/${id}`).then((r) => r.data)
  },

  createUser(input: CreateUserInput): Promise<ManagedUser> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.createUser(input)
    }
    return apiClient.post<ManagedUser>('/admin/users', input).then((r) => r.data)
  },

  updateUser(
    id: string,
    input: UpdateUserInput,
    sessionUserId?: string,
  ): Promise<ManagedUser> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.updateUser(id, input, sessionUserId)
    }
    return apiClient
      .patch<ManagedUser>(`/admin/users/${id}`, input)
      .then((r) => r.data)
  },

  deleteUser(id: string, sessionUserId?: string): Promise<void> {
    if (env.VITE_USE_MOCK_API) {
      return mockAccessControlApi.deleteUser(id, sessionUserId)
    }
    return apiClient.delete(`/admin/users/${id}`).then(() => undefined)
  },
}
