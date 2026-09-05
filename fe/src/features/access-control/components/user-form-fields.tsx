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

import { InputError } from '@/components/common/input-error'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from '@/components/ui/input-group'
import { Label } from '@/components/ui/label'
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
    <Tabs defaultValue="profile" className="min-w-0 gap-4">
      <TabsList
        variant="line"
        className="h-auto w-full max-w-full flex-wrap justify-start"
      >
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

      <TabsContent value="profile" className="min-w-0 space-y-6">
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
      </TabsContent>

      <TabsContent value="password" className="min-w-0 space-y-4">
        <div className="grid min-w-0 gap-2">
          <Label htmlFor="user-password">{t('access.users.fields.password')}</Label>
          {isEdit ? (
            <p className="text-xs text-muted-foreground">
              {t('access.users.fields.passwordOptionalHint')}
            </p>
          ) : (
            <p className="text-xs text-muted-foreground">
              {t('profile.password.generateHint')}
            </p>
          )}
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
          <InputError message={errors.password?.message} />
        </div>
      </TabsContent>

      <TabsContent value="account" className="min-w-0 space-y-4">
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

        {canGrantSuperAdmin ? (
          <div className="flex items-start gap-3 rounded-lg border border-border p-3">
            <Checkbox
              id="user-super-admin"
              checked={superAdmin}
              onCheckedChange={(checked) =>
                setValue('superAdmin', checked === true, {
                  shouldDirty: true,
                })
              }
            />
            <div className="min-w-0 space-y-1">
              <Label htmlFor="user-super-admin" className="font-medium">
                {t('access.users.fields.superAdmin')}
              </Label>
              <p className="text-xs text-muted-foreground">
                {t('access.users.fields.superAdminHint')}
              </p>
            </div>
          </div>
        ) : null}
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
