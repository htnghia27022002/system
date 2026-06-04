'use client'

import { useTranslation } from 'react-i18next'

import { AuthLayout } from '@/components/common/auth-layout'

import { LoginForm } from './login-form'

export function LoginPage() {
  const { t } = useTranslation('common')

  return (
    <AuthLayout
      title={t('auth.login.title')}
      description={t('auth.login.description')}
    >
      <LoginForm />
    </AuthLayout>
  )
}
