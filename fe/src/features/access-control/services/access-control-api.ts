import { env } from '@/config/env'
import { apiClient } from '@/services/api-client'
import { normalizePaginatedResponse } from '@/lib/normalize-paginated-response'
import { loadAccessControlMock, MockAccessControlError } from '@/services/mock'

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
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.listPermissions(params))
    }
    return apiClient
      .get<PaginatedResponse<Permission>>('/admin/permissions', { params })
      .then((r) => normalizePaginatedResponse<Permission>(r.data, params.pageSize))
  },

  listAllPermissions(): Promise<Permission[]> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock()
        .then((mock) => mock.listPermissions({ page: 1, pageSize: 1000 }))
        .then((r) => r.items)
    }
    return apiClient
      .get<Permission[]>('/admin/permissions/all')
      .then((r) => r.data)
  },

  listRoles(params: ListRolesParams = {}): Promise<PaginatedResponse<Role>> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.listRoles(params))
    }
    return apiClient
      .get<PaginatedResponse<Role>>('/admin/roles', {
        params,
        skipNavLoading: true,
      })
      .then((r) => normalizePaginatedResponse<Role>(r.data, params.pageSize))
  },

  listAllRoles(): Promise<Role[]> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock()
        .then((mock) => mock.listRoles({ page: 1, pageSize: 1000 }))
        .then((r) => r.items)
    }
    return apiClient
      .get<Role[]>('/admin/roles/all', { skipNavLoading: true })
      .then((r) => r.data)
  },

  getRole(id: string): Promise<Role> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.getRole(id))
    }
    return apiClient.get<Role>(`/admin/roles/${id}`).then((r) => r.data)
  },

  createRole(input: CreateRoleInput): Promise<Role> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.createRole(input))
    }
    return apiClient.post<Role>('/admin/roles', input).then((r) => r.data)
  },

  updateRole(id: string, input: UpdateRoleInput): Promise<Role> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.updateRole(id, input))
    }
    return apiClient
      .patch<Role>(`/admin/roles/${id}`, input)
      .then((r) => r.data)
  },

  deleteRole(id: string): Promise<void> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.deleteRole(id))
    }
    return apiClient.delete(`/admin/roles/${id}`).then(() => undefined)
  },

  listUsers(params: ListUsersParams = {}): Promise<PaginatedResponse<ManagedUser>> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.listUsers(params))
    }
    return apiClient
      .get<PaginatedResponse<ManagedUser>>('/admin/users', {
        params,
        skipNavLoading: true,
      })
      .then((r) => normalizePaginatedResponse<ManagedUser>(r.data, params.pageSize))
  },

  getUser(id: string): Promise<ManagedUser> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.getUser(id))
    }
    return apiClient.get<ManagedUser>(`/admin/users/${id}`).then((r) => r.data)
  },

  createUser(input: CreateUserInput): Promise<ManagedUser> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.createUser(input))
    }
    return apiClient.post<ManagedUser>('/admin/users', input).then((r) => r.data)
  },

  updateUser(
    id: string,
    input: UpdateUserInput,
    sessionUserId?: string,
  ): Promise<ManagedUser> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) =>
        mock.updateUser(id, input, sessionUserId),
      )
    }
    return apiClient
      .patch<ManagedUser>(`/admin/users/${id}`, input)
      .then((r) => r.data)
  },

  deleteUser(id: string, sessionUserId?: string): Promise<void> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.deleteUser(id, sessionUserId))
    }
    return apiClient.delete(`/admin/users/${id}`).then(() => undefined)
  },

  uploadUserAvatar(id: string, file: File): Promise<ManagedUser> {
    if (env.USE_MOCK_API) {
      return loadAccessControlMock().then((mock) => mock.uploadUserAvatar(id, file))
    }
    const formData = new FormData()
    formData.append('file', file)
    return apiClient
      .post<ManagedUser>(`/admin/users/${id}/avatar`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then((r) => r.data)
  },
}
