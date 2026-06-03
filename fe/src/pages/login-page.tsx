import { useTranslation } from 'react-i18next'

import { LoginForm } from '@/features/auth'
import { AuthLayout } from '@/layouts/auth-layout'

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
