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
import { useRegister } from '../hooks/use-register'
import {
  registerSchema,
  type RegisterFormValues,
} from '../schemas/auth-schemas'

export function RegisterForm() {
  const { t } = useTranslation('common')
  const registerMutation = useRegister()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  })

  return (
    <AuthFormShell
      title={t('auth.register.title')}
      description={t('auth.register.description')}
      footer={
        <AuthFormFooterText>
          {t('auth.register.hasAccount')}{' '}
          <Button variant="link" className="h-auto p-0 text-[17px]" asChild>
            <Link to="/login">{t('auth.actions.signIn')}</Link>
          </Button>
        </AuthFormFooterText>
      }
    >
      <form
        className="space-y-5"
        onSubmit={handleSubmit((values) => registerMutation.mutate(values))}
        noValidate
      >
        <AuthFormField
          id="name"
          label={t('auth.fields.name')}
          error={errors.name?.message}
        >
          <Input
            id="name"
            type="text"
            autoComplete="name"
            className={authFieldInputClassName}
            aria-invalid={Boolean(errors.name)}
            aria-describedby={authFieldDescribedBy('name', Boolean(errors.name))}
            {...register('name')}
          />
        </AuthFormField>

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
            autoComplete="new-password"
            className={authFieldInputClassName}
            aria-invalid={Boolean(errors.password)}
            aria-describedby={authFieldDescribedBy(
              'password',
              Boolean(errors.password),
            )}
            {...register('password')}
          />
        </AuthFormField>

        <AuthFormField
          id="confirmPassword"
          label={t('auth.fields.confirmPassword')}
          error={errors.confirmPassword?.message}
        >
          <Input
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            className={authFieldInputClassName}
            aria-invalid={Boolean(errors.confirmPassword)}
            aria-describedby={authFieldDescribedBy(
              'confirmPassword',
              Boolean(errors.confirmPassword),
            )}
            {...register('confirmPassword')}
          />
        </AuthFormField>

        <Button
          type="submit"
          variant="pill"
          size="cta"
          className="mt-2 w-full"
          disabled={registerMutation.isPending}
          aria-busy={registerMutation.isPending}
        >
          {registerMutation.isPending ? (
            <Loader2Icon className="size-4 animate-spin" aria-hidden />
          ) : null}
          {registerMutation.isPending
            ? t('auth.actions.creatingAccount')
            : t('auth.actions.createAccount')}
        </Button>
      </form>
    </AuthFormShell>
  )
}
