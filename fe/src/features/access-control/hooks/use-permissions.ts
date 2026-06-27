'use client'

import { useCallback, useMemo } from 'react'

import { selectPermissions, useAuthStore } from '@/store/auth-store'

import {
  canModifyResource,
  canViewResource,
  isPermissionAllowed,
  type PermissionResource,
} from '../permission-keys'

export function usePermissions() {
  const permissions = useAuthStore(selectPermissions)
  const sessionSynced = useAuthStore((state) => state.sessionSynced)

  const permissionSet = useMemo(() => new Set(permissions), [permissions])

  const hasPermission = useCallback(
    (key: string) => isPermissionAllowed(permissionSet, key),
    [permissionSet],
  )

  const hasAny = useCallback(
    (...keys: string[]) => keys.some((key) => isPermissionAllowed(permissionSet, key)),
    [permissionSet],
  )

  const canView = useCallback(
    (resource: PermissionResource) => canViewResource(permissionSet, resource),
    [permissionSet],
  )

  const canModify = useCallback(
    (resource: PermissionResource) => canModifyResource(permissionSet, resource),
    [permissionSet],
  )

  return {
    permissions,
    sessionSynced,
    hasPermission,
    hasAny,
    canView,
    canModify,
  }
}
