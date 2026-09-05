'use client'

import { PlusIcon, Trash2Icon } from 'lucide-react'
import type {
  FieldErrors,
  UseFieldArrayRemove,
  UseFormRegister,
} from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'

type SocialLinkRow = { id: string; label?: string; url: string }

type SocialLinksFormShape = {
  socialLinks: Array<{ label?: string; url: string }>
}

type SocialLinksEditorProps = {
  fields: SocialLinkRow[]
  register: UseFormRegister<SocialLinksFormShape>
  errors: FieldErrors<SocialLinksFormShape>
  append: (value: { label?: string; url: string }) => void
  remove: UseFieldArrayRemove
  max?: number
}

export function SocialLinksEditor({
  fields,
  register,
  errors,
  append,
  remove,
  max = 5,
}: SocialLinksEditorProps) {
  const { t } = useTranslation('admin')
  const linkErrors = errors.socialLinks

  return (
    <FieldGroup>
      <Field>
        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0">
            <FieldLabel>{t('profile.fields.socialLinks')}</FieldLabel>
            <FieldDescription>
              {t('profile.fields.socialLinksHint')}
            </FieldDescription>
          </div>
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={fields.length >= max}
            onClick={() => append({ label: '', url: '' })}
          >
            <PlusIcon className="size-4" />
            {t('profile.actions.addLink')}
          </Button>
        </div>
        {typeof linkErrors?.message === 'string' ? (
          <FieldError>{linkErrors.message}</FieldError>
        ) : null}
      </Field>

      {fields.length === 0 ? (
        <FieldDescription>{t('profile.empty.socialLinks')}</FieldDescription>
      ) : (
        <ul className="space-y-3">
          {fields.map((field, index) => {
            const entryError = Array.isArray(linkErrors)
              ? linkErrors[index]
              : undefined
            return (
              <li
                key={field.id}
                className="grid gap-3 rounded-lg border p-3 sm:grid-cols-[1fr_1.4fr_auto]"
              >
                <Field data-invalid={Boolean(entryError?.label)}>
                  <FieldLabel htmlFor={`social-label-${field.id}`}>
                    {t('profile.fields.socialLabel')}
                  </FieldLabel>
                  <Input
                    id={`social-label-${field.id}`}
                    placeholder={t('profile.fields.socialLabelPlaceholder')}
                    aria-invalid={Boolean(entryError?.label)}
                    {...register(`socialLinks.${index}.label`)}
                  />
                  <FieldError errors={[entryError?.label]} />
                </Field>
                <Field data-invalid={Boolean(entryError?.url)}>
                  <FieldLabel htmlFor={`social-url-${field.id}`}>
                    {t('profile.fields.socialUrl')}
                  </FieldLabel>
                  <Input
                    id={`social-url-${field.id}`}
                    type="url"
                    placeholder="https://"
                    aria-invalid={Boolean(entryError?.url)}
                    {...register(`socialLinks.${index}.url`)}
                  />
                  <FieldError errors={[entryError?.url]} />
                </Field>
                <div className="flex items-end">
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="text-destructive"
                    aria-label={t('profile.actions.removeLink')}
                    onClick={() => remove(index)}
                  >
                    <Trash2Icon className="size-4" />
                  </Button>
                </div>
              </li>
            )
          })}
        </ul>
      )}
    </FieldGroup>
  )
}
