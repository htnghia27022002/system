'use client'

import { useRef, useState } from 'react'
import type {
  Control,
  FieldErrors,
  UseFieldArrayAppend,
  UseFieldArrayRemove,
  UseFormRegister,
  UseFormSetValue,
  FieldArrayWithId,
} from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import axios from 'axios'

import { InputError } from '@/components/common/input-error'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
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
import { Textarea } from '@/components/ui/textarea'
import { SocialLinksEditor, validateAvatarFile } from '@/features/user-profile'
import { useInitials } from '@/hooks/use-initials'
import { resolveMediaUrl } from '@/utils/resolve-media-url'

import type {
  CreateUserFormValues,
  UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import { accessControlApi } from '../services/access-control-api'
import type { Role } from '../types'

type UserFormValues = CreateUserFormValues | UpdateUserFormValues

type UserFormFieldsProps = {
  isEdit: boolean
  userId?: string
  userName: string
  avatarUrl?: string | null
  roles: Role[]
  roleId: string
  status: string
  errors: FieldErrors<UserFormValues>
  register: UseFormRegister<UserFormValues>
  setValue: UseFormSetValue<UserFormValues>
  control: Control<UserFormValues>
  socialFields: FieldArrayWithId<UserFormValues, 'socialLinks', 'id'>[]
  appendSocial: UseFieldArrayAppend<UserFormValues, 'socialLinks'>
  removeSocial: UseFieldArrayRemove
  onAvatarUpdated?: (avatarUrl: string) => void
}

export function UserFormFields({
  isEdit,
  userId,
  userName,
  avatarUrl,
  roles,
  roleId,
  status,
  errors,
  register,
  setValue,
  socialFields,
  appendSocial,
  removeSocial,
  onAvatarUpdated,
}: UserFormFieldsProps) {
  const { t } = useTranslation('admin')
  const getInitials = useInitials()
  const inputRef = useRef<HTMLInputElement>(null)
  const [avatarPending, setAvatarPending] = useState(false)
  const [avatarError, setAvatarError] = useState<string | null>(null)
  const avatarSrc = resolveMediaUrl(avatarUrl ?? undefined)

  const onAvatarPick = async (file: File | undefined) => {
    if (!file || !userId || !isEdit) return
    const clientError = validateAvatarFile(file)
    if (clientError) {
      setAvatarError(clientError)
      return
    }
    setAvatarError(null)
    setAvatarPending(true)
    try {
      const updated = await accessControlApi.uploadUserAvatar(userId, file)
      onAvatarUpdated?.(updated.avatarUrl ?? '')
      toast.success(t('access.users.toasts.avatarUpdated'))
    } catch (error) {
      const message = axios.isAxiosError(error)
        ? (error.response?.data as { message?: string } | undefined)?.message ||
          t('access.users.errors.avatarUpload')
        : t('access.users.errors.avatarUpload')
      setAvatarError(message)
      toast.error(message)
    } finally {
      setAvatarPending(false)
      if (inputRef.current) inputRef.current.value = ''
    }
  }

  return (
    <div className="min-w-0 space-y-6">
      <fieldset className="min-w-0 space-y-4">
        <legend className="text-sm font-medium text-foreground">
          {t('access.users.sections.account')}
        </legend>

        {isEdit ? (
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <Avatar className="h-16 w-16 overflow-hidden rounded-full">
              <AvatarImage src={avatarSrc} alt={userName} />
              <AvatarFallback className="rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                {getInitials(userName)}
              </AvatarFallback>
            </Avatar>
            <div className="min-w-0 space-y-2">
              <p className="text-xs text-muted-foreground">
                {t('access.users.fields.avatarHint')}
              </p>
              <Button
                type="button"
                variant="outline"
                size="sm"
                disabled={avatarPending || !userId}
                onClick={() => inputRef.current?.click()}
              >
                {avatarPending ? <Spinner className="size-4" /> : null}
                {t('access.users.actions.changeAvatar')}
              </Button>
              <input
                ref={inputRef}
                type="file"
                accept="image/jpeg,image/png,image/webp"
                className="sr-only"
                aria-label={t('access.users.actions.changeAvatar')}
                onChange={(event) => void onAvatarPick(event.target.files?.[0])}
              />
              <InputError message={avatarError ?? undefined} />
            </div>
          </div>
        ) : null}

        <div className="grid min-w-0 gap-4 sm:grid-cols-2">
          <div className="grid min-w-0 gap-2">
            <Label htmlFor="user-name">{t('access.users.fields.name')}</Label>
            <Input
              id="user-name"
              aria-invalid={Boolean(errors.name)}
              {...register('name')}
            />
            <InputError message={errors.name?.message} />
          </div>

          <div className="grid min-w-0 gap-2">
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
        </div>

        <div className="grid min-w-0 gap-2">
          <Label htmlFor="user-password">{t('access.users.fields.password')}</Label>
          {isEdit ? (
            <p className="text-xs text-muted-foreground">
              {t('access.users.fields.passwordOptionalHint')}
            </p>
          ) : null}
          <Input
            id="user-password"
            type="password"
            autoComplete="new-password"
            aria-invalid={Boolean(errors.password)}
            {...register('password')}
          />
          <InputError message={errors.password?.message} />
        </div>
      </fieldset>

      <fieldset className="min-w-0 space-y-4">
        <legend className="text-sm font-medium text-foreground">
          {t('access.users.sections.personal')}
        </legend>
        <div className="grid min-w-0 gap-4 sm:grid-cols-2">
          <div className="grid min-w-0 gap-2">
            <Label htmlFor="user-phone">{t('access.users.fields.phone')}</Label>
            <Input
              id="user-phone"
              aria-invalid={Boolean(errors.phone)}
              {...register('phone')}
            />
            <InputError message={errors.phone?.message} />
          </div>
          <div className="grid min-w-0 gap-2">
            <Label htmlFor="user-birthday">
              {t('access.users.fields.birthday')}
            </Label>
            <Input
              id="user-birthday"
              type="date"
              aria-invalid={Boolean(errors.birthday)}
              {...register('birthday')}
            />
            <InputError message={errors.birthday?.message} />
          </div>
        </div>
        <div className="grid min-w-0 gap-2">
          <Label htmlFor="user-general">{t('access.users.fields.general')}</Label>
          <Textarea
            id="user-general"
            rows={3}
            aria-invalid={Boolean(errors.general)}
            {...register('general')}
          />
          <InputError message={errors.general?.message} />
        </div>
        <div className="grid min-w-0 gap-2">
          <Label htmlFor="user-address">{t('access.users.fields.address')}</Label>
          <Textarea
            id="user-address"
            rows={2}
            aria-invalid={Boolean(errors.address)}
            {...register('address')}
          />
          <InputError message={errors.address?.message} />
        </div>
        <SocialLinksEditor
          fields={socialFields}
          register={register}
          errors={errors}
          append={appendSocial}
          remove={removeSocial}
        />
      </fieldset>

      <fieldset className="min-w-0 space-y-4">
        <legend className="text-sm font-medium text-foreground">
          {t('access.users.sections.access')}
        </legend>

        <div className="grid min-w-0 gap-4 sm:grid-cols-2">
          <div className="grid min-w-0 gap-2">
            <Label>{t('access.users.fields.role')}</Label>
            <Select
              value={roleId}
              onValueChange={(value) => setValue('roleId', value)}
            >
              <SelectTrigger className="w-full min-w-0">
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

          <div className="grid min-w-0 gap-2">
            <Label>{t('access.users.fields.status')}</Label>
            <Select
              value={status}
              onValueChange={(value) =>
                setValue('status', value as 'active' | 'inactive')
              }
            >
              <SelectTrigger className="w-full min-w-0">
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
      </fieldset>
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
