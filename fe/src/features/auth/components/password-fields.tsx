'use client'

import type { FieldErrors, UseFormRegister } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { InputError } from '@/components/common/input-error'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

type PasswordValues = { password: string; confirmPassword: string }

type PasswordFieldsProps<T extends PasswordValues> = {
  register: UseFormRegister<T>
  errors: FieldErrors<T>
  autoFocus?: boolean
}

/** New password + confirmation, shared by reset-password and accept-invite. */
export function PasswordFields<T extends PasswordValues>({
  register,
  errors,
  autoFocus,
}: PasswordFieldsProps<T>) {
  const { t } = useTranslation('common')
  const reg = register as unknown as UseFormRegister<PasswordValues>
  const errs = errors as FieldErrors<PasswordValues>

  return (
    <>
      <div className="grid gap-2">
        <Label htmlFor="password">{t('auth.fields.newPassword')}</Label>
        <Input
          id="password"
          type="password"
          autoComplete="new-password"
          autoFocus={autoFocus}
          aria-invalid={Boolean(errs.password)}
          {...reg('password')}
        />
        <InputError message={errs.password?.message} />
      </div>
      <div className="grid gap-2">
        <Label htmlFor="confirmPassword">{t('auth.fields.confirmPassword')}</Label>
        <Input
          id="confirmPassword"
          type="password"
          autoComplete="new-password"
          aria-invalid={Boolean(errs.confirmPassword)}
          {...reg('confirmPassword')}
        />
        <InputError message={errs.confirmPassword?.message} />
      </div>
    </>
  )
}
