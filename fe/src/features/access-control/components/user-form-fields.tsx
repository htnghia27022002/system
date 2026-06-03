import type { FieldErrors, UseFormRegister, UseFormSetValue } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'

import type {
  CreateUserFormValues,
  UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import type { Role } from '../types'

type UserFormValues = CreateUserFormValues | UpdateUserFormValues

type UserFormFieldsProps = {
  isEdit: boolean
  roles: Role[]
  roleId: string
  status: string
  errors: FieldErrors<UserFormValues>
  register: UseFormRegister<UserFormValues>
  setValue: UseFormSetValue<UserFormValues>
}

export function UserFormFields({
  isEdit,
  roles,
  roleId,
  status,
  errors,
  register,
  setValue,
}: UserFormFieldsProps) {
  const { t } = useTranslation('admin')

  return (
    <div className="grid gap-4">
      <div className="grid gap-2">
        <Label htmlFor="user-name">{t('access.users.fields.name')}</Label>
        <Input
          id="user-name"
          aria-invalid={Boolean(errors.name)}
          {...register('name')}
        />
        <InputError message={errors.name?.message} />
      </div>

      <div className="grid gap-2">
        <Label htmlFor="user-email">{t('access.users.fields.email')}</Label>
        <Input
          id="user-email"
          type="email"
          autoComplete="email"
          aria-invalid={Boolean(errors.email)}
          {...register('email')}
        />
        <InputError message={errors.email?.message} />
      </div>

      <div className="grid gap-2">
        <Label htmlFor="user-password">
          {t('access.users.fields.password')}
          {isEdit ? ` (${t('access.users.fields.passwordOptional')})` : ''}
        </Label>
        <Input
          id="user-password"
          type="password"
          autoComplete={isEdit ? 'new-password' : 'new-password'}
          aria-invalid={Boolean(errors.password)}
          {...register('password')}
        />
        <InputError message={errors.password?.message} />
      </div>

      <div className="grid gap-2">
        <Label>{t('access.users.fields.role')}</Label>
        <Select
          value={roleId}
          onValueChange={(value) => setValue('roleId', value)}
        >
          <SelectTrigger className="w-full">
            <SelectValue placeholder={t('access.users.fields.role')} />
          </SelectTrigger>
          <SelectContent>
            {roles.map((role) => (
              <SelectItem key={role.id} value={role.id}>
                {role.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <InputError message={errors.roleId?.message} />
      </div>

      <div className="grid gap-2">
        <Label>{t('access.users.fields.status')}</Label>
        <Select
          value={status}
          onValueChange={(value) =>
            setValue('status', value as 'active' | 'inactive')
          }
        >
          <SelectTrigger className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="active">
              {t('access.users.status.active')}
            </SelectItem>
            <SelectItem value="inactive">
              {t('access.users.status.inactive')}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}

type UserFormFooterProps = {
  isEdit: boolean
  isPending: boolean
  onCancel: () => void
}

export function UserFormFooter({
  isEdit,
  isPending,
  onCancel,
}: UserFormFooterProps) {
  const { t } = useTranslation('admin')

  return (
    <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
      <Button
        type="button"
        variant="outline"
        className="w-full sm:w-auto"
        onClick={onCancel}
      >
        {t('access.actions.cancel')}
      </Button>
      <Button type="submit" disabled={isPending} className="w-full sm:w-auto">
        {isPending ? <Spinner className="size-4" /> : null}
        {isEdit ? t('access.actions.save') : t('access.actions.create')}
      </Button>
    </div>
  )
}
