'use client'

import { useState } from 'react'
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
import {
  CopyIcon,
  KeyRoundIcon,
  ShieldCheckIcon,
  SparklesIcon,
  UserRoundIcon,
} from 'lucide-react'

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
import { InputGroupButton } from '@/components/ui/input-group'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Spinner } from '@/components/ui/spinner'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import {
  AvatarUpload,
  PasswordInput,
  PasswordStrengthMeter,
  SocialLinksEditor,
  evaluatePassword,
  generateSecurePassword,
} from '@/features/user-profile'
import { cn } from '@/lib/utils'

import type {
  CreateUserFormValues,
  UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import { accessControlApi } from '../services/access-control-api'
import type { Role } from '../types'

type UserFormValues = CreateUserFormValues | UpdateUserFormValues

const GENERAL_MAX = 1000

const TAB_ICONS = {
  profile: UserRoundIcon,
  password: KeyRoundIcon,
  account: ShieldCheckIcon,
} as const

type UserFormFieldsProps = {
  isEdit: boolean
  userId?: string
  userName: string
  userEmail?: string
  avatarUrl?: string | null
  roles: Role[]
  roleId: string
  status: string
  superAdmin: boolean
  canGrantSuperAdmin: boolean
  password: string
  generalLength: number
  errors: FieldErrors<UserFormValues>
  register: UseFormRegister<UserFormValues>
  setValue: UseFormSetValue<UserFormValues>
  socialFields: FieldArrayWithId<UserFormValues, 'socialLinks', 'id'>[]
  appendSocial: UseFieldArrayAppend<UserFormValues, 'socialLinks'>
  removeSocial: UseFieldArrayRemove
  onAvatarUpdated?: (avatarUrl: string) => void
}

function SectionHeading({ title, description }: { title: string; description?: string }) {
  return (
    <div className="space-y-0.5">
      <h3 className="text-sm font-semibold">{title}</h3>
      {description ? <p className="text-xs text-muted-foreground">{description}</p> : null}
    </div>
  )
}

export function UserFormFields({
  isEdit,
  userId,
  userName,
  userEmail,
  avatarUrl,
  roles,
  roleId,
  status,
  superAdmin,
  canGrantSuperAdmin,
  password,
  generalLength,
  errors,
  register,
  setValue,
  socialFields,
  appendSocial,
  removeSocial,
  onAvatarUpdated,
}: UserFormFieldsProps) {
  const { t } = useTranslation('admin')
  const [revealPassword, setRevealPassword] = useState(false)

  const fillGeneratedPassword = () => {
    setValue('password', generateSecurePassword(), { shouldDirty: true, shouldValidate: true })
    setRevealPassword(true)
    toast.success(t('profile.toasts.passwordGenerated'))
  }

  const copyPassword = async () => {
    if (!password) return
    try {
      await navigator.clipboard.writeText(password)
      toast.success(t('profile.toasts.passwordCopied'))
    } catch {
      toast.error(t('profile.errors.copyFailed'))
    }
  }

  return (
    <Tabs defaultValue="profile" className="w-full min-w-0 gap-5">
      <div className="-mx-4 overflow-x-auto border-b border-border px-4">
        <TabsList variant="line" className="h-10 w-max justify-start gap-4 p-0">
          {(['profile', 'password', 'account'] as const).map((value) => {
            const Icon = TAB_ICONS[value]
            return (
              <TabsTrigger key={value} value={value} className="flex-none px-1 pb-2">
                <Icon aria-hidden />
                {t(`access.users.tabs.${value}`)}
              </TabsTrigger>
            )
          })}
        </TabsList>
      </div>

      <TabsContent value="profile" className="min-w-0">
        <FieldGroup className="gap-5">
          {isEdit && userId ? (
            <>
              <AvatarUpload
                name={userName}
                subtitle={userEmail}
                avatarUrl={avatarUrl}
                upload={(file) => accessControlApi.uploadUserAvatar(userId, file)}
                onUploaded={(url) => onAvatarUpdated?.(url)}
              />
              <Separator />
            </>
          ) : null}

          <SectionHeading title={t('access.users.sections.personal')} />

          <div className="grid min-w-0 gap-4 sm:grid-cols-2">
            <Field data-invalid={Boolean(errors.name)}>
              <FieldLabel htmlFor="user-name">{t('access.users.fields.name')}</FieldLabel>
              <Input
                id="user-name"
                autoComplete="off"
                aria-invalid={Boolean(errors.name)}
                {...register('name')}
              />
              <FieldError errors={[errors.name]} />
            </Field>
            <Field data-invalid={Boolean(errors.email)}>
              <FieldLabel htmlFor="user-email">{t('access.users.fields.email')}</FieldLabel>
              <Input
                id="user-email"
                type="email"
                autoComplete="off"
                aria-invalid={Boolean(errors.email)}
                {...register('email')}
              />
              <FieldError errors={[errors.email]} />
            </Field>
            <Field data-invalid={Boolean(errors.phone)}>
              <FieldLabel htmlFor="user-phone">{t('access.users.fields.phone')}</FieldLabel>
              <Input
                id="user-phone"
                type="tel"
                aria-invalid={Boolean(errors.phone)}
                {...register('phone')}
              />
              <FieldError errors={[errors.phone]} />
            </Field>
            <Field data-invalid={Boolean(errors.birthday)}>
              <FieldLabel htmlFor="user-birthday">{t('access.users.fields.birthday')}</FieldLabel>
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
            <div className="flex items-baseline justify-between gap-2">
              <FieldLabel htmlFor="user-general">{t('access.users.fields.general')}</FieldLabel>
              <span
                className={cn(
                  'text-xs tabular-nums',
                  generalLength > GENERAL_MAX ? 'text-destructive' : 'text-muted-foreground',
                )}
              >
                {generalLength} / {GENERAL_MAX}
              </span>
            </div>
            <Textarea
              id="user-general"
              rows={3}
              aria-invalid={Boolean(errors.general)}
              {...register('general')}
            />
            <FieldError errors={[errors.general]} />
          </Field>

          <Field data-invalid={Boolean(errors.address)}>
            <FieldLabel htmlFor="user-address">{t('access.users.fields.address')}</FieldLabel>
            <Textarea
              id="user-address"
              rows={2}
              aria-invalid={Boolean(errors.address)}
              {...register('address')}
            />
            <FieldError errors={[errors.address]} />
          </Field>

          <Separator />

          <SectionHeading
            title={t('profile.fields.socialLinks')}
            description={t('profile.fields.socialLinksHint')}
          />
          <SocialLinksEditor
            fields={socialFields}
            register={register}
            errors={errors}
            append={appendSocial}
            remove={removeSocial}
            hideHeader
          />
        </FieldGroup>
      </TabsContent>

      <TabsContent value="password" className="min-w-0">
        <FieldGroup className="gap-5">
          <SectionHeading
            title={isEdit ? t('access.users.password.resetTitle') : t('access.users.password.setTitle')}
            description={
              isEdit
                ? t('access.users.fields.passwordOptionalHint')
                : t('profile.sections.passwordDescription')
            }
          />
          <Field data-invalid={Boolean(errors.password)}>
            <div className="flex flex-wrap items-center justify-between gap-2">
              <FieldLabel htmlFor="user-password">{t('access.users.fields.password')}</FieldLabel>
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
              id="user-password"
              autoComplete="new-password"
              aria-invalid={Boolean(errors.password)}
              aria-describedby="user-password-strength"
              revealed={revealPassword}
              onToggleReveal={() => setRevealPassword((open) => !open)}
              extraActions={
                password ? (
                  <InputGroupButton
                    size="icon-xs"
                    aria-label={t('profile.actions.copyPassword')}
                    onClick={() => void copyPassword()}
                  >
                    <CopyIcon />
                  </InputGroupButton>
                ) : null
              }
              {...register('password')}
            />
            <FieldError errors={[errors.password]} />
            {!isEdit || password ? (
              <PasswordStrengthMeter
                id="user-password-strength"
                strength={evaluatePassword(password)}
              />
            ) : null}
          </Field>
        </FieldGroup>
      </TabsContent>

      <TabsContent value="account" className="min-w-0">
        <FieldGroup className="gap-5">
          <SectionHeading
            title={t('access.users.sections.access')}
            description={t('access.users.sections.accessDescription')}
          />
          <div className="grid min-w-0 gap-4 sm:grid-cols-2">
            <Field data-invalid={Boolean(errors.roleId)}>
              <FieldLabel>{t('access.users.fields.role')}</FieldLabel>
              <Select value={roleId} onValueChange={(value) => setValue('roleId', value)}>
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
                onValueChange={(value) => setValue('status', value as 'active' | 'inactive')}
                disabled={status === 'invited'}
              >
                <SelectTrigger className="w-full min-w-0">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {status === 'invited' ? (
                    <SelectItem value="invited" disabled>
                      {t('access.users.status.invited')}
                    </SelectItem>
                  ) : null}
                  <SelectItem value="active">{t('access.users.status.active')}</SelectItem>
                  <SelectItem value="inactive">{t('access.users.status.inactive')}</SelectItem>
                </SelectContent>
              </Select>
            </Field>
          </div>

          {canGrantSuperAdmin ? (
            <label
              htmlFor="user-super-admin"
              className={cn(
                'flex cursor-pointer items-start gap-3 rounded-lg border p-3 transition-colors hover:bg-muted/50 motion-reduce:transition-none',
                superAdmin ? 'border-primary bg-primary/5 hover:bg-primary/5' : 'border-border',
              )}
            >
              <Checkbox
                id="user-super-admin"
                className="mt-0.5"
                checked={superAdmin}
                onCheckedChange={(checked) =>
                  setValue('superAdmin', checked === true, { shouldDirty: true })
                }
              />
              <FieldContent>
                <span className="flex items-center gap-1.5 text-sm font-medium">
                  <ShieldCheckIcon className="size-4 text-primary" aria-hidden />
                  {t('access.users.fields.superAdmin')}
                </span>
                <FieldDescription>{t('access.users.fields.superAdminHint')}</FieldDescription>
              </FieldContent>
            </label>
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
