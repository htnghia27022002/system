import { useEffect } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { usePermissions } from '../hooks/use-permissions'

type PermissionRouteProps = {
  permission: string
}

export function PermissionRoute({ permission }: PermissionRouteProps) {
  const { hasPermission } = usePermissions()
  const { t } = useTranslation('admin')
  const allowed = hasPermission(permission)

  useEffect(() => {
    if (!allowed) {
      toast.error(t('access.insufficientPermissions'))
    }
  }, [allowed, t])

  if (!allowed) {
    return <Navigate to="/admin" replace />
  }

  return <Outlet />
}
