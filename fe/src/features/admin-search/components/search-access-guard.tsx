'use client'

import { useEffect, type ReactNode } from 'react'
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { useSearchNavigation } from '../hooks/use-search-navigation'

type SearchAccessGuardProps = {
  children: ReactNode
}

export function SearchAccessGuard({ children }: SearchAccessGuardProps) {
  const { canUseSearch, sessionSynced } = useSearchNavigation()
  const { t } = useTranslation('admin')
  const router = useRouter()

  useEffect(() => {
    if (!sessionSynced || canUseSearch) {
      return
    }

    toast.error(t('access.insufficientPermissions'))
    router.replace('/admin')
  }, [canUseSearch, sessionSynced, t, router])

  if (!sessionSynced) {
    return null
  }

  if (!canUseSearch) {
    return null
  }

  return <>{children}</>
}
