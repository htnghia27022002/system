import { useCallback, useMemo } from 'react'

import { selectPermissions, useAuthStore } from '@/store/auth-store'

export function usePermissions() {
  const permissions = useAuthStore(selectPermissions)

  const permissionSet = useMemo(() => new Set(permissions), [permissions])

  const hasPermission = useCallback(
    (key: string) => permissionSet.has(key),
    [permissionSet],
  )

  const hasAny = useCallback(
    (...keys: string[]) => keys.some((key) => permissionSet.has(key)),
    [permissionSet],
  )

  return { permissions, hasPermission, hasAny }
}
