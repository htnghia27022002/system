'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import axios from 'axios'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'

import {
  changePasswordSchema,
  type ChangePasswordFormValues,
} from '../schemas/profile-schemas'
import { profileApi } from '../services/profile-api'

type ChangePasswordFormProps = {
  hasPassword: boolean
}

function apiErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: string; error?: string } | undefined
    return data?.message || data?.error || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function ChangePasswordForm({ hasPassword }: ChangePasswordFormProps) {
  const { t } = useTranslation('admin')
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  })

  if (!hasPassword) {
    return (
      <div className="rounded-lg border border-border bg-muted/30 p-4">
        <p className="text-sm text-muted-foreground">
          {t('profile.password.unavailable')}
        </p>
      </div>
    )
  }

  const onSubmit = handleSubmit(async (values) => {
    try {
      await profileApi.changePassword({
        currentPassword: values.currentPassword,
        newPassword: values.newPassword,
      })
      reset()
      toast.success(t('profile.toasts.passwordChanged'))
    } catch (error) {
      toast.error(apiErrorMessage(error, t('profile.errors.passwordFailed')))
    }
  })

  return (
    <form className="space-y-4" onSubmit={onSubmit} noValidate>
      <div className="grid gap-2">
        <Label htmlFor="current-password">
          {t('profile.fields.currentPassword')}
        </Label>
        <Input
          id="current-password"
          type="password"
          autoComplete="current-password"
          aria-invalid={Boolean(errors.currentPassword)}
          {...register('currentPassword')}
        />
        <InputError message={errors.currentPassword?.message} />
      </div>
      <div className="grid gap-2">
        <Label htmlFor="new-password">{t('profile.fields.newPassword')}</Label>
        <Input
          id="new-password"
          type="password"
          autoComplete="new-password"
          aria-invalid={Boolean(errors.newPassword)}
          {...register('newPassword')}
        />
        <InputError message={errors.newPassword?.message} />
      </div>
      <div className="grid gap-2">
        <Label htmlFor="confirm-password">
          {t('profile.fields.confirmPassword')}
        </Label>
        <Input
          id="confirm-password"
          type="password"
          autoComplete="new-password"
          aria-invalid={Boolean(errors.confirmPassword)}
          {...register('confirmPassword')}
        />
        <InputError message={errors.confirmPassword?.message} />
      </div>
      <div className="flex justify-end">
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? <Spinner className="size-4" /> : null}
          {t('profile.actions.changePassword')}
        </Button>
      </div>
    </form>
  )
}
