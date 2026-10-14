'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { AuthLayout } from '@/components/common/auth-layout'
import { InputError } from '@/components/common/input-error'
import { TextLink } from '@/components/common/text-link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

import { useForgotPassword } from '../hooks/use-account-flows'
import {
  forgotPasswordSchema,
  type ForgotPasswordFormValues,
} from '../schemas/auth-schemas'

export function ForgotPasswordPage() {
  const { t } = useTranslation('common')
  const forgot = useForgotPassword()
  const {
    register,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: '' },
  })

  const onSubmit = handleSubmit((values) =>
    forgot.mutate(values.email, {
      onError: () => toast.error(t('auth.forgot.failed')),
    }),
  )

  return (
    <AuthLayout title={t('auth.forgot.title')} description={t('auth.forgot.description')}>
      {forgot.isSuccess ? (
        <div className="flex flex-col gap-4 text-center" role="status">
          <p className="text-sm text-muted-foreground">
            {t('auth.forgot.sent', { email: getValues('email') })}
          </p>
          <TextLink href="/login" className="text-sm">
            {t('auth.links.backToSignIn')}
          </TextLink>
        </div>
      ) : (
        <form className="flex flex-col gap-6" onSubmit={onSubmit} noValidate>
          <div className="grid gap-2">
            <Label htmlFor="email">{t('auth.fields.email')}</Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              autoFocus
              placeholder="email@example.com"
              aria-invalid={Boolean(errors.email)}
              {...register('email')}
            />
            <InputError message={errors.email?.message} />
          </div>
          <Button type="submit" className="w-full" disabled={forgot.isPending} aria-busy={forgot.isPending}>
            {forgot.isPending ? <Spinner /> : null}
            {t('auth.forgot.submit')}
          </Button>
          <TextLink href="/login" className="text-center text-sm">
            {t('auth.links.backToSignIn')}
          </TextLink>
        </form>
      )}
    </AuthLayout>
  )
}
