'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { InputError } from '@/components/common/input-error'
import { TextLink } from '@/components/common/text-link'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

import { useLogin } from '../hooks/use-login'
import { loginSchema, type LoginFormValues } from '../schemas/auth-schemas'
import { AuthSocialLogin } from './auth-social-login'

export function LoginForm() {
  const { t } = useTranslation('common')
  const loginMutation = useLogin()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  return (
    <form
      className="flex flex-col gap-6"
      onSubmit={handleSubmit((values) => loginMutation.mutate(values))}
      noValidate
    >
      <div className="grid gap-5">
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

        <div className="grid gap-2">
          <Label htmlFor="password">{t('auth.fields.password')}</Label>
          <Input
            id="password"
            type="password"
            autoComplete="current-password"
            placeholder="Password"
            aria-invalid={Boolean(errors.password)}
            {...register('password')}
          />
          <InputError message={errors.password?.message} />
        </div>

        <Button
          type="submit"
          className="mt-4 w-full"
          disabled={loginMutation.isPending}
          aria-busy={loginMutation.isPending}
        >
          {loginMutation.isPending ? <Spinner /> : null}
          {loginMutation.isPending
            ? t('auth.actions.signingIn')
            : t('auth.actions.signIn')}
        </Button>

        <AuthSocialLogin />
      </div>

      <div className="text-center text-sm text-muted-foreground">
        {t('auth.login.noAccount')}{' '}
        <TextLink href="/register">{t('auth.actions.register')}</TextLink>
      </div>
    </form>
  )
}
