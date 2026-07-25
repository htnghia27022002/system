'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import axios from 'axios'
import { useEffect } from 'react'
import { useFieldArray, useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'
import { Textarea } from '@/components/ui/textarea'
import type { AuthUser } from '@/types/auth'

import { SocialLinksEditor } from './social-links-editor'
import {
  profileFormSchema,
  type ProfileFormValues,
} from '../schemas/profile-schemas'
import { profileApi } from '../services/profile-api'

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
    formState: { errors, isSubmitting },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileFormSchema),
    defaultValues: toFormValues(user),
  })

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'socialLinks',
  })

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
    <form className="space-y-6" onSubmit={onSubmit} noValidate>
      <fieldset className="space-y-4">
        <legend className="text-sm font-medium text-foreground">
          {t('profile.sections.identity')}
        </legend>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="grid gap-2">
            <Label htmlFor="profile-name">{t('profile.fields.name')}</Label>
            <Input
              id="profile-name"
              aria-invalid={Boolean(errors.name)}
              {...register('name')}
            />
            <InputError message={errors.name?.message} />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="profile-email">{t('profile.fields.email')}</Label>
            <Input
              id="profile-email"
              type="email"
              value={user.email}
              disabled
              readOnly
            />
            <p className="text-xs text-muted-foreground">
              {t('profile.fields.emailHint')}
            </p>
          </div>
        </div>
      </fieldset>

      <fieldset className="space-y-4">
        <legend className="text-sm font-medium text-foreground">
          {t('profile.sections.personal')}
        </legend>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="grid gap-2">
            <Label htmlFor="profile-phone">{t('profile.fields.phone')}</Label>
            <Input
              id="profile-phone"
              aria-invalid={Boolean(errors.phone)}
              {...register('phone')}
            />
            <InputError message={errors.phone?.message} />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="profile-birthday">
              {t('profile.fields.birthday')}
            </Label>
            <Input
              id="profile-birthday"
              type="date"
              aria-invalid={Boolean(errors.birthday)}
              {...register('birthday')}
            />
            <InputError message={errors.birthday?.message} />
          </div>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="profile-general">{t('profile.fields.general')}</Label>
          <Textarea
            id="profile-general"
            rows={4}
            aria-invalid={Boolean(errors.general)}
            {...register('general')}
          />
          <InputError message={errors.general?.message} />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="profile-address">{t('profile.fields.address')}</Label>
          <Textarea
            id="profile-address"
            rows={3}
            aria-invalid={Boolean(errors.address)}
            {...register('address')}
          />
          <InputError message={errors.address?.message} />
        </div>
        <SocialLinksEditor
          fields={fields}
          register={register}
          errors={errors}
          append={append}
          remove={remove}
        />
      </fieldset>

      <div className="flex justify-end">
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? <Spinner className="size-4" /> : null}
          {t('profile.actions.saveProfile')}
        </Button>
      </div>
    </form>
  )
}
