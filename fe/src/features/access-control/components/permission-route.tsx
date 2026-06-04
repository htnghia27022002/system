'use client'

import { useEffect, type ReactNode } from 'react'
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { usePermissions } from '../hooks/use-permissions'

type PermissionGuardProps = {
  permission: string
  children: ReactNode
}

export function PermissionGuard({ permission, children }: PermissionGuardProps) {
  const { hasPermission } = usePermissions()
  const { t } = useTranslation('admin')
  const router = useRouter()
  const allowed = hasPermission(permission)

  useEffect(() => {
    if (!allowed) {
      toast.error(t('access.insufficientPermissions'))
      router.replace('/admin')
    }
  }, [allowed, t, router])

  if (!allowed) return null

  return <>{children}</>
}

/** @deprecated Use PermissionGuard. Kept for backward compatibility. */
export function PermissionRoute(_props: { permission: string }) {
  return null
}
