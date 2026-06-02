import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2Icon } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import {
  AuthFormField,
  AuthFormFooterText,
  AuthFormShell,
  authFieldDescribedBy,
  authFieldInputClassName,
} from './auth-form-shell'
import { useLogin } from '../hooks/use-login'
import { loginSchema, type LoginFormValues } from '../schemas/auth-schemas'

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
    <AuthFormShell
      title={t('auth.login.title')}
      description={t('auth.login.description')}
      footer={
        <AuthFormFooterText>
          {t('auth.login.noAccount')}{' '}
          <Button variant="link" className="h-auto p-0 text-[17px]" asChild>
            <Link to="/register">{t('auth.actions.register')}</Link>
          </Button>
        </AuthFormFooterText>
      }
    >
      <form
        className="space-y-5"
        onSubmit={handleSubmit((values) => loginMutation.mutate(values))}
        noValidate
      >
        <AuthFormField
          id="email"
          label={t('auth.fields.email')}
          error={errors.email?.message}
        >
          <Input
            id="email"
            type="email"
            autoComplete="email"
            className={authFieldInputClassName}
            aria-invalid={Boolean(errors.email)}
            aria-describedby={authFieldDescribedBy('email', Boolean(errors.email))}
            {...register('email')}
          />
        </AuthFormField>

        <AuthFormField
          id="password"
          label={t('auth.fields.password')}
          error={errors.password?.message}
        >
          <Input
            id="password"
            type="password"
            autoComplete="current-password"
            className={authFieldInputClassName}
            aria-invalid={Boolean(errors.password)}
            aria-describedby={authFieldDescribedBy(
              'password',
              Boolean(errors.password),
            )}
            {...register('password')}
          />
        </AuthFormField>

        <Button
          type="submit"
          variant="pill"
          size="cta"
          className="mt-2 w-full"
          disabled={loginMutation.isPending}
          aria-busy={loginMutation.isPending}
        >
          {loginMutation.isPending ? (
            <Loader2Icon className="size-4 animate-spin" aria-hidden />
          ) : null}
          {loginMutation.isPending
            ? t('auth.actions.signingIn')
            : t('auth.actions.signIn')}
        </Button>
      </form>
    </AuthFormShell>
  )
}
