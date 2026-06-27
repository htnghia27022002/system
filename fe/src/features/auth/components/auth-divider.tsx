'use client'

import { useTranslation } from 'react-i18next'

import { Separator } from '@/components/ui/separator'

export function AuthDivider() {
  const { t } = useTranslation('common')

  return (
    <div className="relative">
      <div className="absolute inset-0 flex items-center">
        <Separator className="w-full" />
      </div>
      <div className="relative flex justify-center text-xs uppercase">
        <span className="bg-card px-2 text-muted-foreground">
          {t('auth.oauth.divider')}
        </span>
      </div>
    </div>
  )
}
