'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { Trash2Icon } from 'lucide-react'
import { useEffect, useMemo, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import {
  Stepper,
  StepperItem,
  StepperList,
} from '@/components/ui/stepper'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Field,
  FieldDescription,
  FieldError,
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
import { Spinner } from '@/components/ui/spinner'
import { JsonEditor } from './json-editor'

import { mapsApiErrorMessage } from '../api-error'
import {
  DETAILS_MAPPING_FIELDS,
  HTTP_METHODS,
  LOCATION_MAPPING_FIELDS,
  NEWS_MAPPING_FIELDS,
  PLACE_MAPPING_FIELDS,
  emptySourceFormValues,
  parseCategoryMap,
  sampleSourceFormValues,
  sourceFormSchema,
  sourceFormToWriteInput,
  sourceRecordToFormValues,
  type SourceFormValues,
} from '../schemas/source-form-schema'
import { mapsApi } from '../services/maps-api'
import {
  autoMapFields,
  flattenFieldPaths,
  resolveList,
  suggestListPaths,
} from '../source-json-preview'
import type { SourceRecord, SourceWriteInput } from '../types'

const STEPS = ['connect', 'path', 'map', 'save'] as const
const NONE = '__none__'

type SourceFormDialogProps = {
  open: boolean
  mode: 'create' | 'edit'
  source?: SourceRecord
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (input: SourceWriteInput) => void
  onUpdate: (input: SourceWriteInput) => void
}

export function SourceFormDialog({
  open,
  mode,
  source,
  isPending,
  onOpenChange,
  onCreate,
  onUpdate,
}: SourceFormDialogProps) {
  const { t } = useTranslation('admin')
  const isEdit = mode === 'edit'
  const [step, setStep] = useState(0)
  const [probeBody, setProbeBody] = useState<unknown>(null)

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    setValue,
    watch,
    getValues,
    trigger,
  } = useForm<SourceFormValues>({
    resolver: zodResolver(sourceFormSchema),
    defaultValues: emptySourceFormValues(),
  })

  const {
    fields: headers,
    append: appendHeader,
    remove: removeHeader,
  } = useFieldArray({ control, name: 'headers' })
  const {
    fields: queryParams,
    append: appendParam,
    remove: removeParam,
  } = useFieldArray({ control, name: 'queryParams' })

  const enabled = watch('enabled')
  const httpMethod = watch('httpMethod')
  const listPath = watch('listPath')
  const locationMap = watch('location')
  const placeMap = watch('place')
  const newsMap = watch('news')
  const detailsMap = watch('details')

  useEffect(() => {
    if (!open) return
    setStep(0)
    setProbeBody(null)
    if (isEdit && source) {
      reset(sourceRecordToFormValues(source))
      return
    }
    reset(emptySourceFormValues())
  }, [open, isEdit, source, reset])

  const listItems = useMemo(
    () => (probeBody == null ? [] : resolveList(probeBody, listPath)),
    [listPath, probeBody],
  )
  const apiFields = useMemo(
    () => (listItems[0] ? flattenFieldPaths(listItems[0]) : []),
    [listItems],
  )
  const suggestedPaths = useMemo(
    () => (probeBody == null ? [] : suggestListPaths(probeBody)),
    [probeBody],
  )
  const mappedCount = [
    ...Object.values(locationMap),
    ...Object.values(placeMap),
    ...Object.values(newsMap),
    ...Object.values(detailsMap),
  ].filter((value) => value.trim()).length
  const targetCount =
    LOCATION_MAPPING_FIELDS.length +
    PLACE_MAPPING_FIELDS.length +
    NEWS_MAPPING_FIELDS.length +
    DETAILS_MAPPING_FIELDS.length

  const probe = useMutation({
    mutationFn: () => {
      const values = getValues()
      return mapsApi.probeSource({
        httpMethod: values.httpMethod,
        url: values.url.trim(),
        headers: Object.fromEntries(
          values.headers
            .filter((row) => row.key.trim())
            .map((row) => [row.key.trim(), row.value]),
        ),
        queryParams: Object.fromEntries(
          values.queryParams
            .filter((row) => row.key.trim())
            .map((row) => [row.key.trim(), row.value]),
        ),
        body: values.body,
      })
    },
    onSuccess: (result) => {
      setProbeBody(result.body)
      const paths = suggestListPaths(result.body)
      if (paths.length === 1) {
        setValue('listPath', paths[0])
      }
      toast.success(t('maps.sources.wizard.probeOk'))
      setStep(1)
    },
    onError: (error) => {
      setProbeBody(null)
      toast.error(mapsApiErrorMessage(error, t('maps.sources.wizard.probeFailed')))
    },
  })

  const applyAutoMap = () => {
    const mapped = autoMapFields(apiFields)
    const assign = (
      group: 'location' | 'place' | 'news' | 'details',
      fields: readonly string[],
    ) => {
      for (const field of fields) {
        setValue(`${group}.${field}` as never, (mapped[`${group}.${field}`] ?? '') as never)
      }
    }
    assign('location', LOCATION_MAPPING_FIELDS)
    assign('place', PLACE_MAPPING_FIELDS)
    assign('news', NEWS_MAPPING_FIELDS)
    assign('details', DETAILS_MAPPING_FIELDS)
  }

  const rawPreview = useMemo(() => {
    if (probeBody == null) return ''
    try {
      const text = JSON.stringify(probeBody, null, 2)
      return text.length > 4000 ? `${text.slice(0, 4000)}\n…` : text
    } catch {
      return ''
    }
  }, [probeBody])

  const goNext = async () => {
    if (step === 0) {
      const ok = await trigger(['name', 'url', 'httpMethod'])
      if (!ok) return
      if (probeBody == null && !isEdit) {
        toast.error(t('maps.sources.wizard.needProbe'))
        return
      }
      setStep(1)
      return
    }
    if (step === 1) {
      if (probeBody != null && listItems.length === 0) {
        toast.error(t('maps.sources.wizard.needList'))
        return
      }
      setStep(2)
      return
    }
    setStep(3)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="flex max-h-[min(46rem,calc(100vh-2rem))] w-full flex-col gap-0 overflow-hidden p-0 sm:max-w-3xl">
        <DialogHeader className="border-b border-border px-6 py-4 pr-12">
          <DialogTitle>
            {isEdit ? t('maps.sources.editTitle') : t('maps.sources.createTitle')}
          </DialogTitle>
          <DialogDescription>
            {t('maps.sources.wizard.description')}
          </DialogDescription>
          <Stepper
            className="mt-3 gap-0"
            value={STEPS[step]}
            onValueChange={(value) => {
              const next = STEPS.indexOf(value as (typeof STEPS)[number])
              if (next >= 0 && next <= step) setStep(next)
            }}
            steps={STEPS.map((value, index) => ({
              value,
              disabled: index > step,
            }))}
          >
            <StepperList>
              {STEPS.map((key, index) => (
                <StepperItem
                  key={key}
                  value={key}
                  completed={index < step}
                  disabled={index > step}
                >
                  {t(`maps.sources.wizard.steps.${key}`)}
                </StepperItem>
              ))}
            </StepperList>
          </Stepper>
        </DialogHeader>

        <form
          className="flex min-h-0 flex-1 flex-col"
          onSubmit={handleSubmit((values) => {
            try {
              parseCategoryMap(values.categoryMap)
              const input = sourceFormToWriteInput(values)
              if (isEdit) onUpdate(input)
              else onCreate(input)
            } catch {
              toast.error(t('maps.sources.fields.categoryMapHint'))
            }
          })}
          noValidate
        >
          <div className="min-h-0 flex-1 overflow-y-auto px-6 py-4">
            {step === 0 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
                  {t('maps.sources.wizard.connectHint')}
                </p>
                {!isEdit ? (
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => reset(sampleSourceFormValues())}
                  >
                    {t('maps.sources.useSample')}
                  </Button>
                ) : null}
                <Field>
                  <FieldLabel htmlFor="maps-src-name">
                    {t('maps.sources.fields.name')}
                  </FieldLabel>
                  <Input id="maps-src-name" {...register('name')} />
                  {errors.name ? (
                    <FieldError>{t('maps.sources.fields.name')}</FieldError>
                  ) : null}
                </Field>
                <div className="grid gap-4 sm:grid-cols-[8rem_minmax(0,1fr)]">
                  <Field>
                    <FieldLabel>{t('maps.sources.fields.httpMethod')}</FieldLabel>
                    <Select
                      value={httpMethod}
                      onValueChange={(value) =>
                        setValue(
                          'httpMethod',
                          value as SourceFormValues['httpMethod'],
                        )
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent className="z-[1500]">
                        {HTTP_METHODS.map((method) => (
                          <SelectItem key={method} value={method}>
                            {method}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </Field>
                  <Field>
                    <FieldLabel htmlFor="maps-src-url">
                      {t('maps.sources.fields.url')}
                    </FieldLabel>
                    <Input id="maps-src-url" {...register('url')} />
                    <FieldDescription>
                      {t('maps.sources.fields.urlHint')}
                    </FieldDescription>
                    {errors.url ? (
                      <FieldError>{t('maps.sources.fields.url')}</FieldError>
                    ) : null}
                  </Field>
                </div>
                <KvEditor
                  title={t('maps.sources.sections.headers')}
                  empty={t('maps.sources.fields.noHeaders')}
                  addLabel={t('maps.sources.fields.addHeader')}
                  keyLabel={t('maps.sources.fields.headerKey')}
                  valueLabel={t('maps.sources.fields.headerValue')}
                  rows={headers}
                  register={register}
                  name="headers"
                  onAdd={() => appendHeader({ key: '', value: '' })}
                  onRemove={removeHeader}
                />
                <KvEditor
                  title={t('maps.sources.sections.queryParams')}
                  empty={t('maps.sources.fields.noParams')}
                  addLabel={t('maps.sources.fields.addParam')}
                  keyLabel={t('maps.sources.fields.paramKey')}
                  valueLabel={t('maps.sources.fields.paramValue')}
                  rows={queryParams}
                  register={register}
                  name="queryParams"
                  onAdd={() => appendParam({ key: '', value: '' })}
                  onRemove={removeParam}
                />
                {httpMethod !== 'GET' ? (
                  <Field>
                    <FieldLabel htmlFor="maps-src-body">
                      {t('maps.sources.fields.body')}
                    </FieldLabel>
                    <FieldDescription>
                      {t('maps.sources.fields.bodyHint')}
                    </FieldDescription>
                    <JsonEditor
                      id="maps-src-body"
                      value={watch('body')}
                      invalid={Boolean(errors.body)}
                      onChange={(next) =>
                        setValue('body', next, {
                          shouldDirty: true,
                          shouldValidate: true,
                        })
                      }
                    />
                  </Field>
                ) : null}
                <Button
                  type="button"
                  onClick={() => void probe.mutate()}
                  disabled={probe.isPending}
                >
                  {probe.isPending ? <Spinner className="size-4" /> : null}
                  {t('maps.sources.wizard.probe')}
                </Button>
              </div>
            ) : null}

            {step === 1 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
                  {t('maps.sources.wizard.pathHint', {
                    count: apiFields.length,
                  })}
                </p>
                {probeBody != null ? (
                  <pre className="max-h-52 overflow-auto rounded-lg border border-border bg-muted/40 p-3 text-xs">
                    {rawPreview}
                  </pre>
                ) : (
                  <p className="text-sm text-muted-foreground">
                    {t('maps.sources.wizard.noPreview')}
                  </p>
                )}
                <Field>
                  <FieldLabel htmlFor="maps-src-list-path">
                    {t('maps.sources.fields.listPath')}
                  </FieldLabel>
                  <Input
                    id="maps-src-list-path"
                    {...register('listPath')}
                    placeholder={suggestedPaths[0] || t('maps.sources.fields.listPathHint')}
                  />
                  {suggestedPaths.length > 0 ? (
                    <div className="flex flex-wrap gap-2 pt-1">
                      {suggestedPaths.map((path) => (
                        <Button
                          key={path || 'root'}
                          type="button"
                          size="sm"
                          variant="outline"
                          onClick={() => setValue('listPath', path)}
                        >
                          {path || t('maps.sources.wizard.rootArray')}
                        </Button>
                      ))}
                    </div>
                  ) : null}
                </Field>
              </div>
            ) : null}

            {step === 2 ? (
              <div className="space-y-4">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <p className="text-sm text-muted-foreground">
                    {t('maps.sources.wizard.mapHint', {
                      mapped: mappedCount,
                      total: targetCount,
                      fields: apiFields.length,
                    })}
                  </p>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={applyAutoMap}
                    disabled={apiFields.length === 0}
                  >
                    {t('maps.sources.wizard.autoMap')}
                  </Button>
                </div>
                <div className="grid gap-4 lg:grid-cols-2">
                  <div className="rounded-lg border border-border p-3">
                    <p className="mb-2 text-sm font-medium">
                      {t('maps.sources.wizard.apiFields')}
                    </p>
                    {apiFields.length === 0 ? (
                      <p className="text-sm text-muted-foreground">
                        {t('maps.sources.wizard.noFields')}
                      </p>
                    ) : (
                      <ul className="space-y-1 text-xs">
                        {apiFields.map((field) => (
                          <li
                            key={field}
                            className="rounded-md bg-muted/50 px-2 py-1 font-mono"
                          >
                            {field}
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                  <div className="space-y-4">
                    <MappingSelects
                      group="location"
                      fields={LOCATION_MAPPING_FIELDS}
                      values={locationMap}
                      options={apiFields}
                      setValue={setValue}
                    />
                    <MappingSelects
                      group="place"
                      fields={PLACE_MAPPING_FIELDS}
                      values={placeMap}
                      options={apiFields}
                      setValue={setValue}
                    />
                    <MappingSelects
                      group="news"
                      fields={NEWS_MAPPING_FIELDS}
                      values={newsMap}
                      options={apiFields}
                      setValue={setValue}
                    />
                    <MappingSelects
                      group="details"
                      fields={DETAILS_MAPPING_FIELDS}
                      values={detailsMap}
                      options={apiFields}
                      setValue={setValue}
                    />
                  </div>
                </div>
              </div>
            ) : null}

            {step === 3 ? (
              <div className="space-y-4">
                <p className="text-sm text-muted-foreground">
                  {t('maps.sources.wizard.saveHint')}
                </p>
                <label className="flex items-start gap-2 text-sm">
                  <Checkbox
                    checked={enabled}
                    onCheckedChange={(checked) =>
                      setValue('enabled', checked === true)
                    }
                  />
                  <span>
                    {t('maps.sources.fields.enabled')}
                    <span className="mt-1 block text-xs text-muted-foreground">
                      {t('maps.sources.fields.enabledHint')}
                    </span>
                  </span>
                </label>
                <p className="text-sm">
                  {watch('httpMethod')} · {watch('url')}
                </p>
                <p className="text-sm text-muted-foreground">
                  {t('maps.sources.wizard.mapHint', {
                    mapped: mappedCount,
                    total: targetCount,
                    fields: apiFields.length,
                  })}
                </p>
              </div>
            ) : null}
          </div>

          <DialogFooter className="mx-0 mb-0 rounded-none">
            {step > 0 ? (
              <Button type="button" variant="outline" onClick={() => setStep(step - 1)}>
                {t('maps.sources.wizard.back')}
              </Button>
            ) : (
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                {t('access.actions.cancel')}
              </Button>
            )}
            {step < 3 ? (
              <Button type="button" onClick={() => void goNext()}>
                {t('maps.sources.wizard.next')}
              </Button>
            ) : (
              <Button type="submit" disabled={isPending}>
                {isPending ? <Spinner className="size-4" /> : null}
                {isEdit ? t('access.actions.save') : t('access.actions.create')}
              </Button>
            )}
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function KvEditor({
  title,
  empty,
  addLabel,
  keyLabel,
  valueLabel,
  rows,
  register,
  name,
  onAdd,
  onRemove,
}: {
  title: string
  empty: string
  addLabel: string
  keyLabel: string
  valueLabel: string
  rows: { id: string }[]
  register: ReturnType<typeof useForm<SourceFormValues>>['register']
  name: 'headers' | 'queryParams'
  onAdd: () => void
  onRemove: (index: number) => void
}) {
  const { t } = useTranslation('admin')
  return (
    <fieldset className="space-y-2">
      <legend className="text-sm font-medium">{title}</legend>
      {rows.length === 0 ? (
        <p className="text-sm text-muted-foreground">{empty}</p>
      ) : (
        rows.map((row, index) => (
          <div key={row.id} className="grid gap-2 sm:grid-cols-[1fr_1fr_auto]">
            <Input
              aria-label={keyLabel}
              {...register(`${name}.${index}.key`)}
              placeholder={keyLabel}
            />
            <Input
              aria-label={valueLabel}
              {...register(`${name}.${index}.value`)}
              placeholder={valueLabel}
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="text-destructive"
              aria-label={t('maps.sources.fields.removeRow')}
              onClick={() => onRemove(index)}
            >
              <Trash2Icon className="size-4" />
            </Button>
          </div>
        ))
      )}
      <Button type="button" variant="outline" size="sm" onClick={onAdd}>
        {addLabel}
      </Button>
    </fieldset>
  )
}

function MappingSelects({
  group,
  fields,
  values,
  options,
  setValue,
}: {
  group: 'location' | 'place' | 'news' | 'details'
  fields: readonly string[]
  values: Record<string, string>
  options: string[]
  setValue: ReturnType<typeof useForm<SourceFormValues>>['setValue']
}) {
  const { t } = useTranslation('admin')
  return (
    <fieldset className="space-y-2 rounded-lg border border-border p-3">
      <legend className="px-1 text-sm font-medium">
        {t(`maps.sources.mappingGroups.${group}`)}
      </legend>
      <div className="grid gap-3">
        {fields.map((field) => {
          const current = values[field] ?? ''
          return (
            <div
              key={`${group}.${field}`}
              className="grid gap-1 sm:grid-cols-[8rem_minmax(0,1fr)] sm:items-center"
            >
              <FieldLabel className="text-xs">
                {t(`maps.sources.mappingFields.${field}`)}
              </FieldLabel>
              <Select
                value={current || NONE}
                onValueChange={(value) =>
                  setValue(
                    `${group}.${field}` as never,
                    (value === NONE ? '' : value) as never,
                  )
                }
              >
                <SelectTrigger className="h-9">
                  <SelectValue placeholder={t('maps.sources.wizard.skipField')} />
                </SelectTrigger>
                <SelectContent className="z-[1500]">
                  <SelectItem value={NONE}>
                    {t('maps.sources.wizard.skipField')}
                  </SelectItem>
                  {options.map((option) => (
                    <SelectItem key={option} value={option}>
                      {option}
                    </SelectItem>
                  ))}
                  {current && !options.includes(current) ? (
                    <SelectItem value={current}>{current}</SelectItem>
                  ) : null}
                </SelectContent>
              </Select>
            </div>
          )
        })}
      </div>
    </fieldset>
  )
}
