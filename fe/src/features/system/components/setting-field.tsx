'use client'

import { useTranslation } from 'react-i18next'

import { Checkbox } from '@/components/ui/checkbox'
import { Field, FieldContent, FieldDescription, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'

import type { AppSetting } from '../types'

type SettingFieldProps = {
  setting: AppSetting
  value: string | boolean | number
  disabled: boolean
  onChange: (value: string | boolean | number) => void
}

/** Renders one BE-declared setting by type. Module settings fall back to the BE label. */
export function SettingField({ setting, value, disabled, onChange }: SettingFieldProps) {
  const { t } = useTranslation('admin')
  const id = `setting-${setting.key}`
  const label = t(`system.settings.items.${setting.key}.label`, { defaultValue: setting.label })
  const description = t(`system.settings.items.${setting.key}.description`, {
    defaultValue: setting.description ?? '',
  })

  if (setting.type === 'bool') {
    return (
      <Field orientation="horizontal">
        <Checkbox
          id={id}
          checked={value === true}
          disabled={disabled}
          onCheckedChange={(checked) => onChange(checked === true)}
        />
        <FieldContent>
          <FieldLabel htmlFor={id}>{label}</FieldLabel>
          {description ? <FieldDescription>{description}</FieldDescription> : null}
        </FieldContent>
      </Field>
    )
  }

  return (
    <Field>
      <FieldLabel htmlFor={id}>{label}</FieldLabel>
      {setting.type === 'int' ? (
        <Input
          id={id}
          type="number"
          inputMode="numeric"
          className="max-w-40"
          min={setting.min}
          max={setting.max}
          disabled={disabled}
          value={String(value)}
          onChange={(event) => onChange(event.target.value === '' ? '' : Number(event.target.value))}
        />
      ) : (
        <Input
          id={id}
          maxLength={setting.maxLength || 500}
          disabled={disabled}
          value={String(value ?? '')}
          onChange={(event) => onChange(event.target.value)}
        />
      )}
      {description ? <FieldDescription>{description}</FieldDescription> : null}
    </Field>
  )
}
