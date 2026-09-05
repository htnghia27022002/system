'use client'

import { useRef, useState } from 'react'
import type {
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
import { EyeIcon, EyeOffIcon, WandSparklesIcon } from 'lucide-react'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from '@/components/ui/input-group'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import {
  SocialLinksEditor,
  generateSecurePassword,
  validateAvatarFile,
} from '@/features/user-profile'
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
  superAdmin: boolean
  canGrantSuperAdmin: boolean
  errors: FieldErrors<UserFormValues>
  register: UseFormRegister<UserFormValues>
  setValue: UseFormSetValue<UserFormValues>
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
  superAdmin,
  canGrantSuperAdmin,
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
  const [revealPassword, setRevealPassword] = useState(false)
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

  const fillGeneratedPassword = () => {
    const password = generateSecurePassword()
    setValue('password', password, { shouldDirty: true, shouldValidate: true })
    toast.success(t('profile.toasts.passwordGenerated'))
  }

  return (
    <Tabs defaultValue="profile" className="w-full min-w-0 gap-6">
      <TabsList variant="line" className="h-auto w-fit justify-start pb-1">
        <TabsTrigger value="profile" className="flex-none px-3">
          {t('access.users.tabs.profile')}
        </TabsTrigger>
        <TabsTrigger value="password" className="flex-none px-3">
          {t('access.users.tabs.password')}
        </TabsTrigger>
        <TabsTrigger value="account" className="flex-none px-3">
          {t('access.users.tabs.account')}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="profile" className="min-w-0">
        <FieldGroup>
          {isEdit ? (
            <Field>
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                <Avatar className="h-16 w-16 overflow-hidden rounded-full">
                  <AvatarImage src={avatarSrc} alt={userName} />
                  <AvatarFallback>{getInitials(userName)}</AvatarFallback>
                </Avatar>
                <FieldContent>
                  <FieldDescription>
                    {t('access.users.fields.avatarHint')}
                  </FieldDescription>
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
                    onChange={(event) =>
                      void onAvatarPick(event.target.files?.[0])
                    }
                  />
                  <FieldError>{avatarError}</FieldError>
                </FieldContent>
              </div>
            </Field>
          ) : null}

          <div className="grid min-w-0 gap-4 sm:grid-cols-2">
            <Field data-invalid={Boolean(errors.name)}>
              <FieldLabel htmlFor="user-name">
                {t('access.users.fields.name')}
              </FieldLabel>
              <Input
                id="user-name"
                aria-invalid={Boolean(errors.name)}
                {...register('name')}
              />
              <FieldError errors={[errors.name]} />
            </Field>

            <Field data-invalid={Boolean(errors.email)}>
              <FieldLabel htmlFor="user-email">
                {t('access.users.fields.email')}
              </FieldLabel>
              <Input
                id="user-email"
                type="email"
                autoComplete="email"
                aria-invalid={Boolean(errors.email)}
                {...register('email')}
              />
              <FieldError errors={[errors.email]} />
            </Field>
          </div>

          <div className="grid min-w-0 gap-4 sm:grid-cols-2">
            <Field data-invalid={Boolean(errors.phone)}>
              <FieldLabel htmlFor="user-phone">
                {t('access.users.fields.phone')}
              </FieldLabel>
              <Input
                id="user-phone"
                aria-invalid={Boolean(errors.phone)}
                {...register('phone')}
              />
              <FieldError errors={[errors.phone]} />
            </Field>
            <Field data-invalid={Boolean(errors.birthday)}>
              <FieldLabel htmlFor="user-birthday">
                {t('access.users.fields.birthday')}
              </FieldLabel>
              <Input
                id="user-birthday"
                type="date"
                aria-invalid={Boolean(errors.birthday)}
                {...register('birthday')}
              />
              <FieldError errors={[errors.birthday]} />
            </Field>
          </div>

          <Field data-invalid={Boolean(errors.general)}>
            <FieldLabel htmlFor="user-general">
              {t('access.users.fields.general')}
            </FieldLabel>
            <Textarea
              id="user-general"
              rows={3}
              aria-invalid={Boolean(errors.general)}
              {...register('general')}
            />
            <FieldError errors={[errors.general]} />
          </Field>

          <Field data-invalid={Boolean(errors.address)}>
            <FieldLabel htmlFor="user-address">
              {t('access.users.fields.address')}
            </FieldLabel>
            <Textarea
              id="user-address"
              rows={2}
              aria-invalid={Boolean(errors.address)}
              {...register('address')}
            />
            <FieldError errors={[errors.address]} />
          </Field>

          <SocialLinksEditor
            fields={socialFields}
            register={register}
            errors={errors}
            append={appendSocial}
            remove={removeSocial}
          />
        </FieldGroup>
      </TabsContent>

      <TabsContent value="password" className="min-w-0">
        <FieldGroup>
          <Field data-invalid={Boolean(errors.password)}>
            <FieldLabel htmlFor="user-password">
              {t('access.users.fields.password')}
            </FieldLabel>
            <FieldDescription>
              {isEdit
                ? t('access.users.fields.passwordOptionalHint')
                : t('profile.password.generateHint')}
            </FieldDescription>
            <InputGroup>
              <InputGroupInput
                id="user-password"
                type={revealPassword ? 'text' : 'password'}
                autoComplete="new-password"
                spellCheck={false}
                aria-invalid={Boolean(errors.password)}
                {...register('password')}
              />
              <InputGroupAddon align="inline-end">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <InputGroupButton
                      size="icon-xs"
                      aria-label={
                        revealPassword
                          ? t('access.users.actions.hidePassword')
                          : t('access.users.actions.showPassword')
                      }
                      onClick={() => setRevealPassword((open) => !open)}
                    >
                      {revealPassword ? (
                        <EyeOffIcon className="size-3.5" />
                      ) : (
                        <EyeIcon className="size-3.5" />
                      )}
                    </InputGroupButton>
                  </TooltipTrigger>
                  <TooltipContent>
                    {revealPassword
                      ? t('access.users.actions.hidePassword')
                      : t('access.users.actions.showPassword')}
                  </TooltipContent>
                </Tooltip>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <InputGroupButton
                      size="icon-xs"
                      aria-label={t('profile.actions.generatePassword')}
                      onClick={fillGeneratedPassword}
                    >
                      <WandSparklesIcon className="size-3.5" />
                    </InputGroupButton>
                  </TooltipTrigger>
                  <TooltipContent>
                    {t('profile.actions.generatePassword')}
                  </TooltipContent>
                </Tooltip>
              </InputGroupAddon>
            </InputGroup>
            <FieldError errors={[errors.password]} />
          </Field>
        </FieldGroup>
      </TabsContent>

      <TabsContent value="account" className="min-w-0">
        <FieldGroup>
          <div className="grid min-w-0 gap-4 sm:grid-cols-2">
            <Field data-invalid={Boolean(errors.roleId)}>
              <FieldLabel>{t('access.users.fields.role')}</FieldLabel>
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
              <FieldError errors={[errors.roleId]} />
            </Field>

            <Field>
              <FieldLabel>{t('access.users.fields.status')}</FieldLabel>
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
            </Field>
          </div>

          {canGrantSuperAdmin ? (
            <Field orientation="horizontal">
              <Checkbox
                id="user-super-admin"
                checked={superAdmin}
                onCheckedChange={(checked) =>
                  setValue('superAdmin', checked === true, {
                    shouldDirty: true,
                  })
                }
              />
              <FieldContent>
                <FieldLabel htmlFor="user-super-admin">
                  {t('access.users.fields.superAdmin')}
                </FieldLabel>
                <FieldDescription>
                  {t('access.users.fields.superAdminHint')}
                </FieldDescription>
              </FieldContent>
            </Field>
          ) : null}
        </FieldGroup>
      </TabsContent>
    </Tabs>
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
    <>
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
    </>
  )
}
