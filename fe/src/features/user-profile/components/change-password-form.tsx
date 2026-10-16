'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import axios from 'axios'
import {
  CheckIcon,
  CopyIcon,
  KeyRoundIcon,
  SparklesIcon,
  XIcon,
} from 'lucide-react'
import { useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { InputGroupButton } from '@/components/ui/input-group'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Spinner } from '@/components/ui/spinner'

import { PasswordInput } from './password-input'
import { PasswordStrengthMeter } from './password-strength-meter'
import { SettingsSection } from './settings-section'
import {
  changePasswordSchema,
  type ChangePasswordFormValues,
} from '../schemas/profile-schemas'
import { profileApi } from '../services/profile-api'
import { generateSecurePassword } from '../utils/generate-secure-password'
import { evaluatePassword } from '../utils/password-strength'

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
  const [revealCurrent, setRevealCurrent] = useState(false)
  const [revealNew, setRevealNew] = useState(false)
  const {
    register,
    control,
    handleSubmit,
    reset,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    },
  })

  const [newPassword, confirmPassword] = useWatch({
    control,
    name: ['newPassword', 'confirmPassword'],
  })
  const strength = evaluatePassword(newPassword ?? '')
  const showMatch = Boolean(confirmPassword) && !errors.confirmPassword
  const passwordsMatch = newPassword === confirmPassword

  if (!hasPassword) {
    return (
      <SettingsSection
        title={t('profile.sections.password')}
        description={t('profile.sections.passwordDescription')}
      >
        <Card>
          <CardContent className="flex flex-col items-center gap-3 py-6 text-center">
            <span className="flex size-10 items-center justify-center rounded-full bg-muted">
              <KeyRoundIcon className="size-5 text-muted-foreground" aria-hidden />
            </span>
            <div className="space-y-1">
              <p className="font-medium">{t('profile.password.unavailableTitle')}</p>
              <p className="max-w-md text-sm text-muted-foreground">
                {t('profile.password.unavailable')}
              </p>
            </div>
          </CardContent>
        </Card>
      </SettingsSection>
    )
  }

  const fillGeneratedPassword = () => {
    const password = generateSecurePassword()
    setValue('newPassword', password, { shouldDirty: true, shouldValidate: true })
    setValue('confirmPassword', password, { shouldDirty: true, shouldValidate: true })
    setRevealNew(true)
    toast.success(t('profile.toasts.passwordGenerated'))
  }

  const copyNewPassword = async () => {
    if (!newPassword) return
    try {
      await navigator.clipboard.writeText(newPassword)
      toast.success(t('profile.toasts.passwordCopied'))
    } catch {
      toast.error(t('profile.errors.copyFailed'))
    }
  }

  const onSubmit = handleSubmit(async (values) => {
    try {
      await profileApi.changePassword({
        currentPassword: values.currentPassword,
        newPassword: values.newPassword,
      })
      reset()
      setRevealCurrent(false)
      setRevealNew(false)
      toast.success(t('profile.toasts.passwordChanged'))
    } catch (error) {
      toast.error(apiErrorMessage(error, t('profile.errors.passwordFailed')))
    }
  })

  return (
    <SettingsSection
      title={t('profile.sections.password')}
      description={t('profile.sections.passwordDescription')}
    >
      <Card className="pb-0">
        <form onSubmit={onSubmit} noValidate className="flex flex-col gap-4">
          <CardContent className="space-y-5">
            <div className="grid gap-2">
              <Label htmlFor="current-password">{t('profile.fields.currentPassword')}</Label>
              <PasswordInput
                id="current-password"
                autoComplete="current-password"
                aria-invalid={Boolean(errors.currentPassword)}
                revealed={revealCurrent}
                onToggleReveal={() => setRevealCurrent((v) => !v)}
                {...register('currentPassword')}
              />
              <InputError message={errors.currentPassword?.message} />
            </div>

            <Separator />

            <div className="grid gap-2">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <Label htmlFor="new-password">{t('profile.fields.newPassword')}</Label>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="text-primary hover:text-primary"
                  onClick={fillGeneratedPassword}
                >
                  <SparklesIcon />
                  {t('profile.actions.generatePassword')}
                </Button>
              </div>
              <PasswordInput
                id="new-password"
                autoComplete="new-password"
                aria-invalid={Boolean(errors.newPassword)}
                aria-describedby="new-password-strength"
                revealed={revealNew}
                onToggleReveal={() => setRevealNew((v) => !v)}
                extraActions={
                  newPassword ? (
                    <InputGroupButton
                      size="icon-xs"
                      aria-label={t('profile.actions.copyPassword')}
                      onClick={() => void copyNewPassword()}
                    >
                      <CopyIcon />
                    </InputGroupButton>
                  ) : null
                }
                {...register('newPassword')}
              />
              <InputError message={errors.newPassword?.message} />
              <PasswordStrengthMeter id="new-password-strength" strength={strength} />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="confirm-password">{t('profile.fields.confirmPassword')}</Label>
              <PasswordInput
                id="confirm-password"
                autoComplete="new-password"
                aria-invalid={Boolean(errors.confirmPassword)}
                revealed={revealNew}
                onToggleReveal={() => setRevealNew((v) => !v)}
                {...register('confirmPassword')}
              />
              {errors.confirmPassword ? (
                <InputError message={errors.confirmPassword.message} />
              ) : showMatch ? (
                <p
                  className={
                    passwordsMatch
                      ? 'flex items-center gap-1.5 text-xs text-success'
                      : 'flex items-center gap-1.5 text-xs text-muted-foreground'
                  }
                >
                  {passwordsMatch ? (
                    <CheckIcon className="size-3.5" aria-hidden />
                  ) : (
                    <XIcon className="size-3.5" aria-hidden />
                  )}
                  {passwordsMatch
                    ? t('profile.password.match')
                    : t('profile.password.mismatch')}
                </p>
              ) : null}
            </div>
          </CardContent>
          <CardFooter className="flex-col-reverse items-stretch gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-xs text-muted-foreground">{t('profile.password.footerHint')}</p>
            <Button type="submit" disabled={isSubmitting} className="sm:w-auto">
              {isSubmitting ? <Spinner className="size-4" /> : null}
              {t('profile.actions.changePassword')}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </SettingsSection>
  )
}
