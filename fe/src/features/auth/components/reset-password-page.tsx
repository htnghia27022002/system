'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useRouter, useSearchParams } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { AuthLayout } from '@/components/common/auth-layout'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { apiErrorStatus } from '@/lib/api-error'

import { useResetPassword } from '../hooks/use-account-flows'
import {
  resetPasswordSchema,
  type ResetPasswordFormValues,
} from '../schemas/auth-schemas'
import { InvalidLinkNotice } from './invalid-link-notice'
import { PasswordFields } from './password-fields'

export function ResetPasswordPage() {
  const { t } = useTranslation('common')
  const router = useRouter()
  const token = useSearchParams().get('token') ?? ''
  const reset = useResetPassword()
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { password: '', confirmPassword: '' },
  })

  const linkInvalid = !token || apiErrorStatus(reset.error) === 400

  const onSubmit = handleSubmit((values) =>
    reset.mutate(
      { token, password: values.password },
      {
        onSuccess: () => {
          toast.success(t('auth.reset.success'))
          router.replace('/login')
        },
        onError: (error) => {
          if (apiErrorStatus(error) !== 400) toast.error(t('auth.reset.failed'))
        },
      },
    ),
  )

  return (
    <AuthLayout title={t('auth.reset.title')} description={t('auth.reset.description')}>
      {linkInvalid ? (
        <InvalidLinkNotice retryHref="/forgot-password" retryLabel={t('auth.reset.requestNew')} />
      ) : (
        <form className="flex flex-col gap-6" onSubmit={onSubmit} noValidate>
          <div className="grid gap-5">
            <PasswordFields register={register} errors={errors} autoFocus />
          </div>
          <Button type="submit" className="w-full" disabled={reset.isPending} aria-busy={reset.isPending}>
            {reset.isPending ? <Spinner /> : null}
            {t('auth.reset.submit')}
          </Button>
        </form>
      )}
    </AuthLayout>
  )
}
