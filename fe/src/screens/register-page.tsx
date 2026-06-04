'use client'

import { useTranslation } from 'react-i18next'

import { RegisterForm } from '@/features/auth'
import { AuthLayout } from '@/layouts/auth-layout'

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
