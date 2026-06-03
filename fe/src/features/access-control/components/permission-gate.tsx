import type { ReactNode } from 'react'

import { usePermissions } from '../hooks/use-permissions'

type PermissionGateProps = {
  permission: string
  children: ReactNode
  fallback?: ReactNode
}

export function PermissionGate({
  permission,
  children,
  fallback = null,
}: PermissionGateProps) {
  const { hasPermission } = usePermissions()

  if (!hasPermission(permission)) {
    return fallback
  }

  return children
}
