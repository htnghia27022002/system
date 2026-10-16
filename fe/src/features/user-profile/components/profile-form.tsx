'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import axios from 'axios'
import { BadgeCheckIcon, CircleAlertIcon, LockIcon } from 'lucide-react'
import { useEffect } from 'react'
import { useFieldArray, useForm, useWatch } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from '@/components/ui/input-group'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'
import { Textarea } from '@/components/ui/textarea'
import { cn } from '@/lib/utils'
import type { AuthUser } from '@/types/auth'

import { SettingsSection } from './settings-section'
import { SocialLinksEditor } from './social-links-editor'
import {
  profileFormSchema,
  type ProfileFormValues,
} from '../schemas/profile-schemas'
import { profileApi } from '../services/profile-api'

const GENERAL_MAX = 1000

type ProfileFormProps = {
  user: AuthUser
  onSaved: (user: AuthUser) => void
}

function toFormValues(user: AuthUser): ProfileFormValues {
  return {
    name: user.name ?? '',
    phone: user.phone ?? '',
    general: user.general ?? '',
    birthday: user.birthday ?? '',
    address: user.address ?? '',
    socialLinks: (user.socialLinks ?? []).map((link) => ({
      label: link.label ?? '',
      url: link.url,
    })),
  }
}

function apiErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: string; error?: string } | undefined
    return data?.message || data?.error || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function ProfileForm({ user, onSaved }: ProfileFormProps) {
  const { t } = useTranslation('admin')
  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting, isDirty },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileFormSchema),
    defaultValues: toFormValues(user),
  })

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'socialLinks',
  })

  const generalLength = useWatch({ control, name: 'general' })?.length ?? 0

  useEffect(() => {
    reset(toFormValues(user))
  }, [user, reset])

  const onSubmit = handleSubmit(async (values) => {
    try {
      const phone = values.phone?.trim() ?? ''
      const updated = await profileApi.updateProfile({
        name: values.name.trim(),
        phone,
        general: values.general?.trim() ?? '',
        birthday: values.birthday?.trim() ? values.birthday.trim() : null,
        address: values.address?.trim() ?? '',
        socialLinks: values.socialLinks.map((link) => ({
          label: link.label?.trim() || undefined,
          url: link.url.trim(),
        })),
      })
      onSaved(updated)
      toast.success(t('profile.toasts.profileSaved'))
    } catch (error) {
      toast.error(apiErrorMessage(error, t('profile.errors.saveFailed')))
    }
  })

  return (
    <form className="flex flex-col gap-8" onSubmit={onSubmit} noValidate>
      <SettingsSection
        title={t('profile.sections.personalInfo')}
        description={t('profile.sections.personalInfoDescription')}
      >
        <Card>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="grid content-start gap-2">
                <Label htmlFor="profile-name">{t('profile.fields.name')}</Label>
                <Input
                  id="profile-name"
                  autoComplete="name"
                  aria-invalid={Boolean(errors.name)}
                  {...register('name')}
                />
                <InputError message={errors.name?.message} />
              </div>
              <div className="grid content-start gap-2">
                <div className="flex items-center justify-between gap-2">
                  <Label htmlFor="profile-email">{t('profile.fields.email')}</Label>
                  {user.emailVerified === false ? (
                    <span className="inline-flex items-center gap-1 text-xs font-medium text-warning">
                      <CircleAlertIcon className="size-3.5" aria-hidden />
                      {t('profile.fields.emailUnverified')}
                    </span>
                  ) : (
                    <span className="inline-flex items-center gap-1 text-xs font-medium text-success">
                      <BadgeCheckIcon className="size-3.5" aria-hidden />
                      {t('profile.fields.emailVerified')}
                    </span>
                  )}
                </div>
                <InputGroup>
                  <InputGroupInput
                    id="profile-email"
                    type="email"
                    value={user.email}
                    disabled
                    readOnly
                    aria-describedby="profile-email-hint"
                  />
                  <InputGroupAddon align="inline-end">
                    <LockIcon aria-hidden />
                  </InputGroupAddon>
                </InputGroup>
                <p id="profile-email-hint" className="text-xs text-muted-foreground">
                  {t('profile.fields.emailHint')}
                </p>
              </div>
              <div className="grid content-start gap-2">
                <Label htmlFor="profile-phone">{t('profile.fields.phone')}</Label>
                <Input
                  id="profile-phone"
                  type="tel"
                  autoComplete="tel"
                  aria-invalid={Boolean(errors.phone)}
                  {...register('phone')}
                />
                <InputError message={errors.phone?.message} />
              </div>
              <div className="grid content-start gap-2">
                <Label htmlFor="profile-birthday">{t('profile.fields.birthday')}</Label>
                <Input
                  id="profile-birthday"
                  type="date"
                  autoComplete="bday"
                  aria-invalid={Boolean(errors.birthday)}
                  {...register('birthday')}
                />
                <InputError message={errors.birthday?.message} />
              </div>
            </div>

            <div className="grid gap-2">
              <div className="flex items-baseline justify-between gap-2">
                <Label htmlFor="profile-general">{t('profile.fields.general')}</Label>
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
                id="profile-general"
                rows={4}
                placeholder={t('profile.fields.generalPlaceholder')}
                aria-invalid={Boolean(errors.general)}
                {...register('general')}
              />
              <InputError message={errors.general?.message} />
            </div>

            <div className="grid gap-2">
              <Label htmlFor="profile-address">{t('profile.fields.address')}</Label>
              <Textarea
                id="profile-address"
                rows={2}
                autoComplete="street-address"
                aria-invalid={Boolean(errors.address)}
                {...register('address')}
              />
              <InputError message={errors.address?.message} />
            </div>
          </CardContent>
        </Card>
      </SettingsSection>

      <SettingsSection
        title={t('profile.fields.socialLinks')}
        description={t('profile.fields.socialLinksHint')}
      >
        <Card>
          <CardContent>
            <SocialLinksEditor
              fields={fields}
              register={register}
              errors={errors}
              append={append}
              remove={remove}
              hideHeader
            />
          </CardContent>
        </Card>
      </SettingsSection>

      <div
        className={cn(
          'grid lg:grid-cols-[minmax(0,16rem)_minmax(0,1fr)] lg:gap-10',
          isDirty && 'pointer-events-none sticky bottom-4 z-10',
        )}
      >
        <div
          className={cn(
            'flex flex-col-reverse gap-3 sm:flex-row sm:items-center sm:justify-end lg:col-start-2',
            isDirty &&
              'pointer-events-auto rounded-xl border border-border bg-card p-3 shadow-lg sm:justify-between',
          )}
        >
          {isDirty ? (
            <p className="text-sm text-muted-foreground" aria-live="polite">
              {t('profile.unsavedChanges')}
            </p>
          ) : null}
          <div className="flex flex-col-reverse gap-2 sm:flex-row">
            {isDirty ? (
              <Button
                type="button"
                variant="outline"
                disabled={isSubmitting}
                onClick={() => reset(toFormValues(user))}
              >
                {t('profile.actions.discard')}
              </Button>
            ) : null}
            <Button type="submit" disabled={isSubmitting || !isDirty}>
              {isSubmitting ? <Spinner className="size-4" /> : null}
              {t('profile.actions.saveProfile')}
            </Button>
          </div>
        </div>
      </div>
    </form>
  )
}
