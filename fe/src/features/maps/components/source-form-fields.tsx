'use client'

import type {
  FieldErrors,
  UseFieldArrayAppend,
  UseFieldArrayRemove,
  UseFormRegister,
  UseFormSetValue,
} from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

import {
  DETAILS_MAPPING_FIELDS,
  HTTP_METHODS,
  LOCATION_MAPPING_FIELDS,
  NEWS_MAPPING_FIELDS,
  PLACE_MAPPING_FIELDS,
  type SourceFormValues,
} from '../schemas/source-form-schema'

type KvRow = { key: string; value: string }

type SourceFormFieldsProps = {
  errors: FieldErrors<SourceFormValues>
  register: UseFormRegister<SourceFormValues>
  setValue: UseFormSetValue<SourceFormValues>
  enabled: boolean
  httpMethod: SourceFormValues['httpMethod']
  headers: KvRow[]
  queryParams: KvRow[]
  appendHeader: UseFieldArrayAppend<SourceFormValues, 'headers'>
  removeHeader: UseFieldArrayRemove
  appendParam: UseFieldArrayAppend<SourceFormValues, 'queryParams'>
  removeParam: UseFieldArrayRemove
}

function MappingInputs({
  group,
  fields,
  register,
}: {
  group: 'location' | 'place' | 'news' | 'details'
  fields: readonly string[]
  register: UseFormRegister<SourceFormValues>
}) {
  const { t } = useTranslation('admin')
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {fields.map((field) => (
        <Field key={`${group}.${field}`}>
          <FieldLabel htmlFor={`maps-map-${group}-${field}`}>
            {t(`maps.sources.mappingFields.${field}`)}
          </FieldLabel>
          <Input
            id={`maps-map-${group}-${field}`}
            {...register(`${group}.${field}` as never)}
            placeholder={t('maps.sources.fields.mappingPathHint')}
          />
        </Field>
      ))}
    </div>
  )
}

export function SourceFormFields({
  errors,
  register,
  setValue,
  enabled,
  httpMethod,
  headers,
  queryParams,
  appendHeader,
  removeHeader,
  appendParam,
  removeParam,
}: SourceFormFieldsProps) {
  const { t } = useTranslation('admin')

  return (
    <div className="space-y-5">
      <fieldset className="space-y-4 rounded-xl border border-border bg-card p-4">
        <legend className="px-1 text-sm font-medium">
          {t('maps.sources.sections.identity')}
        </legend>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field>
            <FieldLabel htmlFor="maps-source-name">
              {t('maps.sources.fields.name')}
            </FieldLabel>
            <Input id="maps-source-name" {...register('name')} />
            {errors.name ? <FieldError>{errors.name.message}</FieldError> : null}
          </Field>
          <Field>
            <FieldLabel htmlFor="maps-source-enabled">
              {t('maps.sources.fields.enabled')}
            </FieldLabel>
            <FieldDescription>
              {t('maps.sources.fields.enabledHint')}
            </FieldDescription>
            <label
              htmlFor="maps-source-enabled"
              className="flex items-center gap-2 text-sm"
            >
              <Checkbox
                id="maps-source-enabled"
                checked={enabled}
                onCheckedChange={(checked) =>
                  setValue('enabled', checked === true)
                }
              />
              {enabled
                ? t('maps.sources.enabled')
                : t('maps.sources.disabled')}
            </label>
          </Field>
        </div>
      </fieldset>

      <fieldset className="space-y-4 rounded-xl border border-border bg-card p-4">
        <legend className="px-1 text-sm font-medium">
          {t('maps.sources.sections.request')}
        </legend>
        <div className="grid gap-4 sm:grid-cols-[10rem_1fr]">
          <Field>
            <FieldLabel>{t('maps.sources.fields.httpMethod')}</FieldLabel>
            <Select
              value={httpMethod}
              onValueChange={(value) =>
                setValue('httpMethod', value as SourceFormValues['httpMethod'])
              }
            >
              <SelectTrigger aria-label={t('maps.sources.fields.httpMethod')}>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {HTTP_METHODS.map((method) => (
                  <SelectItem key={method} value={method}>
                    {method}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel htmlFor="maps-source-url">
              {t('maps.sources.fields.url')}
            </FieldLabel>
            <Input id="maps-source-url" {...register('url')} />
            {errors.url ? <FieldError>{errors.url.message}</FieldError> : null}
            <FieldDescription>{t('maps.sources.fields.urlHint')}</FieldDescription>
          </Field>
        </div>
        <Field>
          <FieldLabel htmlFor="maps-source-body">
            {t('maps.sources.fields.body')}
          </FieldLabel>
          <FieldDescription>{t('maps.sources.fields.bodyHint')}</FieldDescription>
          <Textarea id="maps-source-body" rows={4} {...register('body')} />
        </Field>
      </fieldset>

      <fieldset className="space-y-4 rounded-xl border border-border bg-card p-4">
        <legend className="px-1 text-sm font-medium">
          {t('maps.sources.sections.headers')}
        </legend>
        {headers.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {t('maps.sources.fields.noHeaders')}
          </p>
        ) : null}
        <FieldGroup>
          {headers.map((_row, index) => (
            <div key={`header-${index}`} className="grid gap-4 sm:grid-cols-2">
              <Field>
                <FieldLabel htmlFor={`maps-header-key-${index}`}>
                  {t('maps.sources.fields.headerKey')}
                </FieldLabel>
                <Input
                  id={`maps-header-key-${index}`}
                  {...register(`headers.${index}.key`)}
                />
              </Field>
              <Field>
                <FieldLabel htmlFor={`maps-header-value-${index}`}>
                  {t('maps.sources.fields.headerValue')}
                </FieldLabel>
                <div className="flex gap-2">
                  <Input
                    id={`maps-header-value-${index}`}
                    {...register(`headers.${index}.value`)}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => removeHeader(index)}
                  >
                    {t('maps.sources.fields.removeRow')}
                  </Button>
                </div>
              </Field>
            </div>
          ))}
        </FieldGroup>
        <Button
          type="button"
          variant="outline"
          onClick={() => appendHeader({ key: '', value: '' })}
        >
          {t('maps.sources.fields.addHeader')}
        </Button>
      </fieldset>

      <fieldset className="space-y-4 rounded-xl border border-border bg-card p-4">
        <legend className="px-1 text-sm font-medium">
          {t('maps.sources.sections.queryParams')}
        </legend>
        {queryParams.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {t('maps.sources.fields.noParams')}
          </p>
        ) : null}
        <FieldGroup>
          {queryParams.map((_row, index) => (
            <div key={`param-${index}`} className="grid gap-4 sm:grid-cols-2">
              <Field>
                <FieldLabel htmlFor={`maps-param-key-${index}`}>
                  {t('maps.sources.fields.paramKey')}
                </FieldLabel>
                <Input
                  id={`maps-param-key-${index}`}
                  {...register(`queryParams.${index}.key`)}
                />
              </Field>
              <Field>
                <FieldLabel htmlFor={`maps-param-value-${index}`}>
                  {t('maps.sources.fields.paramValue')}
                </FieldLabel>
                <div className="flex gap-2">
                  <Input
                    id={`maps-param-value-${index}`}
                    {...register(`queryParams.${index}.value`)}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => removeParam(index)}
                  >
                    {t('maps.sources.fields.removeRow')}
                  </Button>
                </div>
              </Field>
            </div>
          ))}
        </FieldGroup>
        <Button
          type="button"
          variant="outline"
          onClick={() => appendParam({ key: '', value: '' })}
        >
          {t('maps.sources.fields.addParam')}
        </Button>
      </fieldset>

      <fieldset className="space-y-4 rounded-xl border border-border bg-card p-4">
        <legend className="px-1 text-sm font-medium">
          {t('maps.sources.sections.mapping')}
        </legend>
        <Field>
          <FieldLabel htmlFor="maps-list-path">
            {t('maps.sources.fields.listPath')}
          </FieldLabel>
          <FieldDescription>
            {t('maps.sources.fields.listPathHint')}
          </FieldDescription>
          <Input id="maps-list-path" {...register('listPath')} />
        </Field>

        <div className="space-y-3">
          <h3 className="text-sm font-medium">
            {t('maps.sources.mappingGroups.location')}
          </h3>
          <MappingInputs
            group="location"
            fields={LOCATION_MAPPING_FIELDS}
            register={register}
          />
        </div>
        <div className="space-y-3">
          <h3 className="text-sm font-medium">
            {t('maps.sources.mappingGroups.place')}
          </h3>
          <MappingInputs
            group="place"
            fields={PLACE_MAPPING_FIELDS}
            register={register}
          />
        </div>
        <div className="space-y-3">
          <h3 className="text-sm font-medium">
            {t('maps.sources.mappingGroups.news')}
          </h3>
          <MappingInputs
            group="news"
            fields={NEWS_MAPPING_FIELDS}
            register={register}
          />
        </div>
        <div className="space-y-3">
          <h3 className="text-sm font-medium">
            {t('maps.sources.mappingGroups.details')}
          </h3>
          <MappingInputs
            group="details"
            fields={DETAILS_MAPPING_FIELDS}
            register={register}
          />
        </div>
        <Field>
          <FieldLabel htmlFor="maps-category-map">
            {t('maps.sources.fields.categoryMap')}
          </FieldLabel>
          <FieldDescription>
            {t('maps.sources.fields.categoryMapHint')}
          </FieldDescription>
          <Textarea
            id="maps-category-map"
            rows={4}
            {...register('categoryMap')}
          />
        </Field>
      </fieldset>
    </div>
  )
}
