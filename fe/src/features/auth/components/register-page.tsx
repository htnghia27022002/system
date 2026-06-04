'use client'

import { useTranslation } from 'react-i18next'

import { AuthLayout } from '@/components/common/auth-layout'

import { RegisterForm } from './register-form'

export function RegisterPage() {
  const { t } = useTranslation('common')

  return (
    <AuthLayout
      title={t('auth.register.title')}
      description={t('auth.register.description')}
    >
      <RegisterForm />
    </AuthLayout>
  )
}
