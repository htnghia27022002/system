'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { FormOverlay } from '@/components/common/form-overlay'
import { Button } from '@/components/ui/button'
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Spinner } from '@/components/ui/spinner'

import { PermissionPicker } from './permission-picker'
import { usePermissionsCatalog } from '../hooks/use-permissions-catalog'
import {
  createRoleSchema,
  updateRoleSchema,
  type CreateRoleFormValues,
  type UpdateRoleFormValues,
} from '../schemas/access-control-schemas'
import type { Role } from '../types'

type RoleFormDialogProps = {
  open: boolean
  mode: 'create' | 'edit'
  role?: Role
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (values: CreateRoleFormValues) => void
  onUpdate: (values: UpdateRoleFormValues) => void
}

export function RoleFormDialog({
  open,
  mode,
  role,
  isPending,
  onOpenChange,
  onCreate,
  onUpdate,
}: RoleFormDialogProps) {
  const { t } = useTranslation('admin')
  const permissionsQuery = usePermissionsCatalog()
  const isEdit = mode === 'edit'

  const schema = useMemo(
    () => (isEdit ? updateRoleSchema : createRoleSchema),
    [isEdit],
  )

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
    watch,
  } = useForm<CreateRoleFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      slug: '',
      permissionKeys: [],
    },
  })

  const selectedKeys = watch('permissionKeys') ?? []
  const permissions = permissionsQuery.data ?? []

  useEffect(() => {
    if (!open) return
    if (isEdit && role) {
      reset({
        name: role.name,
        slug: role.slug,
        permissionKeys: role.permissionKeys,
      })
      return
    }
    reset({ name: '', slug: '', permissionKeys: [] })
  }, [open, isEdit, role, reset])

  return (
    <FormOverlay
      open={open}
      onOpenChange={onOpenChange}
      title={
        isEdit ? t('access.roles.editTitle') : t('access.roles.createTitle')
      }
      description={t('access.roles.permissionsSheetHint')}
      className="sm:max-w-2xl lg:max-w-3xl"
      onSubmit={handleSubmit((values) => {
        if (isEdit) {
          onUpdate(values)
        } else {
          onCreate(values)
        }
      })}
      footer={
        <>
          <Button
            type="button"
            variant="outline"
            className="w-full sm:w-auto"
            onClick={() => onOpenChange(false)}
          >
            {t('access.actions.cancel')}
          </Button>
          <Button
            type="submit"
            disabled={isPending}
            className="w-full sm:w-auto"
          >
            {isPending ? <Spinner className="size-4" /> : null}
            {isEdit ? t('access.actions.save') : t('access.actions.create')}
          </Button>
        </>
      }
    >
      <FieldGroup>
        <div className="grid gap-4 sm:grid-cols-2">
          <Field data-invalid={Boolean(errors.name)}>
            <FieldLabel htmlFor="role-name">
              {t('access.roles.fields.name')}
            </FieldLabel>
            <Input
              id="role-name"
              aria-invalid={Boolean(errors.name)}
              {...register('name')}
            />
            <FieldError errors={[errors.name]} />
          </Field>

          <Field data-invalid={Boolean(errors.slug)}>
            <FieldLabel htmlFor="role-slug">
              {t('access.roles.fields.slug')}
            </FieldLabel>
            <Input
              id="role-slug"
              aria-invalid={Boolean(errors.slug)}
              {...register('slug')}
            />
            <FieldError errors={[errors.slug]} />
          </Field>
        </div>

        <PermissionPicker
          permissions={permissions}
          selectedKeys={selectedKeys}
          onChange={(keys) =>
            setValue('permissionKeys', keys, { shouldValidate: true })
          }
          error={errors.permissionKeys?.message}
        />
      </FieldGroup>
    </FormOverlay>
  )
}
