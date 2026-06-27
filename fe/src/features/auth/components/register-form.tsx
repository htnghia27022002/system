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

import { useRegister } from '../hooks/use-register'
import {
  registerSchema,
  type RegisterFormValues,
} from '../schemas/auth-schemas'
import { AuthSocialLogin } from './auth-social-login'

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
    <form
      className="flex flex-col gap-6"
      onSubmit={handleSubmit((values) => registerMutation.mutate(values))}
      noValidate
    >
      <div className="grid gap-6">
        <div className="grid gap-2">
          <Label htmlFor="name">{t('auth.fields.name')}</Label>
          <Input
            id="name"
            type="text"
            autoComplete="name"
            autoFocus
            placeholder="Full name"
            aria-invalid={Boolean(errors.name)}
            {...register('name')}
          />
          <InputError message={errors.name?.message} />
        </div>

        <div className="grid gap-2">
          <Label htmlFor="email">{t('auth.fields.email')}</Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
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
            autoComplete="new-password"
            placeholder="Password"
            aria-invalid={Boolean(errors.password)}
            {...register('password')}
          />
          <InputError message={errors.password?.message} />
        </div>

        <div className="grid gap-2">
          <Label htmlFor="confirmPassword">
            {t('auth.fields.confirmPassword')}
          </Label>
          <Input
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            placeholder="Confirm password"
            aria-invalid={Boolean(errors.confirmPassword)}
            {...register('confirmPassword')}
          />
          <InputError message={errors.confirmPassword?.message} />
        </div>

        <Button
          type="submit"
          className="mt-2 w-full"
          disabled={registerMutation.isPending}
          aria-busy={registerMutation.isPending}
        >
          {registerMutation.isPending ? <Spinner /> : null}
          {registerMutation.isPending
            ? t('auth.actions.creatingAccount')
            : t('auth.actions.createAccount')}
        </Button>

        <AuthSocialLogin />
      </div>

      <div className="text-center text-sm text-muted-foreground">
        {t('auth.register.hasAccount')}{' '}
        <TextLink href="/login">{t('auth.actions.signIn')}</TextLink>
      </div>
    </form>
  )
}
