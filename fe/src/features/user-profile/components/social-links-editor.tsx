'use client'

import { PlusIcon, Trash2Icon } from 'lucide-react'
import type {
  Control,
  FieldArrayWithId,
  FieldErrors,
  UseFieldArrayAppend,
  UseFieldArrayRemove,
  UseFormRegister,
} from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

type SocialLinksFormShape = {
  socialLinks: Array<{ label?: string; url: string }>
}

type SocialLinksEditorProps<T extends SocialLinksFormShape> = {
  fields: FieldArrayWithId<T, 'socialLinks', 'id'>[]
  register: UseFormRegister<T>
  errors: FieldErrors<T>
  append: UseFieldArrayAppend<T, 'socialLinks'>
  remove: UseFieldArrayRemove
  control?: Control<T>
  max?: number
}

export function SocialLinksEditor<T extends SocialLinksFormShape>({
  fields,
  register,
  errors,
  append,
  remove,
  max = 5,
}: SocialLinksEditorProps<T>) {
  const { t } = useTranslation('admin')
  const linkErrors = errors.socialLinks

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between gap-2">
        <div>
          <Label>{t('profile.fields.socialLinks')}</Label>
          <p className="text-xs text-muted-foreground">
            {t('profile.fields.socialLinksHint')}
          </p>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={fields.length >= max}
          onClick={() => append({ label: '', url: '' } as never)}
        >
          <PlusIcon className="size-4" />
          {t('profile.actions.addLink')}
        </Button>
      </div>

      {typeof linkErrors?.message === 'string' ? (
        <InputError message={linkErrors.message} />
      ) : null}

      {fields.length === 0 ? (
        <p className="text-sm text-muted-foreground">
          {t('profile.empty.socialLinks')}
        </p>
      ) : (
        <ul className="space-y-3">
          {fields.map((field, index) => {
            const entryError = Array.isArray(linkErrors)
              ? linkErrors[index]
              : undefined
            return (
              <li
                key={field.id}
                className="grid gap-3 rounded-lg border border-border p-3 sm:grid-cols-[1fr_1.4fr_auto]"
              >
                <div className="grid gap-2">
                  <Label htmlFor={`social-label-${field.id}`}>
                    {t('profile.fields.socialLabel')}
                  </Label>
                  <Input
                    id={`social-label-${field.id}`}
                    placeholder={t('profile.fields.socialLabelPlaceholder')}
                    aria-invalid={Boolean(entryError?.label)}
                    {...register(`socialLinks.${index}.label` as never)}
                  />
                  <InputError message={entryError?.label?.message as string} />
                </div>
                <div className="grid gap-2">
                  <Label htmlFor={`social-url-${field.id}`}>
                    {t('profile.fields.socialUrl')}
                  </Label>
                  <Input
                    id={`social-url-${field.id}`}
                    type="url"
                    placeholder="https://"
                    aria-invalid={Boolean(entryError?.url)}
                    {...register(`socialLinks.${index}.url` as never)}
                  />
                  <InputError message={entryError?.url?.message as string} />
                </div>
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
    </div>
  )
}
