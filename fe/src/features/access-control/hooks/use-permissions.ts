'use client'

import { useCallback, useMemo } from 'react'

import {
  selectIsSuperAdmin,
  selectPermissions,
  useAuthStore,
} from '@/store/auth-store'

import {
  canModifyResource,
  canViewResource,
  isPermissionAllowed,
  type PermissionResource,
} from '../permission-keys'

export function usePermissions() {
  const permissions = useAuthStore(selectPermissions)
  const isSuperAdmin = useAuthStore(selectIsSuperAdmin)
  const sessionSynced = useAuthStore((state) => state.sessionSynced)

  const permissionSet = useMemo(() => new Set(permissions), [permissions])

  const hasPermission = useCallback(
    (key: string) =>
      isSuperAdmin || isPermissionAllowed(permissionSet, key),
    [isSuperAdmin, permissionSet],
  )

  const hasAny = useCallback(
    (...keys: string[]) =>
      isSuperAdmin ||
      keys.some((key) => isPermissionAllowed(permissionSet, key)),
    [isSuperAdmin, permissionSet],
  )

  const canView = useCallback(
    (resource: PermissionResource) =>
      isSuperAdmin || canViewResource(permissionSet, resource),
    [isSuperAdmin, permissionSet],
  )

  const canModify = useCallback(
    (resource: PermissionResource) =>
      isSuperAdmin || canModifyResource(permissionSet, resource),
    [isSuperAdmin, permissionSet],
  )

  return {
    permissions,
    sessionSynced,
    isSuperAdmin,
    hasPermission,
    hasAny,
    canView,
    canModify,
  }
}
