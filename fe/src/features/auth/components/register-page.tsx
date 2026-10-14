'use client'

import { useTranslation } from 'react-i18next'

import { AuthLayout } from '@/components/common/auth-layout'
import { TextLink } from '@/components/common/text-link'
import { useRegistrationEnabled } from '@/hooks/use-app-info'

import { RegisterForm } from './register-form'

export function RegisterPage() {
  const { t } = useTranslation('common')
  const registrationEnabled = useRegistrationEnabled()

  return (
    <AuthLayout
      title={t('auth.register.title')}
      description={t('auth.register.description')}
    >
      {registrationEnabled ? (
        <RegisterForm />
      ) : (
        <div className="flex flex-col gap-4 text-center" role="status">
          <p className="text-sm text-muted-foreground">{t('auth.register.closed')}</p>
          <TextLink href="/login" className="text-sm">
            {t('auth.links.backToSignIn')}
          </TextLink>
        </div>
      )}
    </AuthLayout>
  )
}
