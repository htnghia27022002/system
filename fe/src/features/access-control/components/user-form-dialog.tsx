'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { FormOverlay } from '@/components/common/form-overlay'
import { emptyPersonalFields } from '@/features/user-profile'

import { usePermissions } from '../hooks/use-permissions'
import { UserFormFields, UserFormFooter } from './user-form-fields'
import {
  createUserSchema,
  updateUserSchema,
  type CreateUserFormValues,
  type UpdateUserFormValues,
} from '../schemas/access-control-schemas'
import type { ManagedUser, Role } from '../types'

type UserFormDialogProps = {
  open: boolean
  mode: 'create' | 'edit'
  user?: ManagedUser
  roles: Role[]
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (values: CreateUserFormValues) => void
  onUpdate: (values: UpdateUserFormValues) => void
}

const defaultPersonal = emptyPersonalFields()

export function UserFormDialog({
  open,
  mode,
  user,
  roles,
  isPending,
  onOpenChange,
  onCreate,
  onUpdate,
}: UserFormDialogProps) {
  const { t } = useTranslation('admin')
  const { isSuperAdmin } = usePermissions()
  const isEdit = mode === 'edit'
  const schema = useMemo(
    () => (isEdit ? updateUserSchema : createUserSchema),
    [isEdit],
  )
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>()

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    setValue,
    watch,
  } = useForm<CreateUserFormValues | UpdateUserFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
      roleId: roles[0]?.id ?? '',
      status: 'active',
      superAdmin: false,
      ...defaultPersonal,
    },
  })

  const {
    fields: socialFields,
    append: appendSocial,
    remove: removeSocial,
  } = useFieldArray({
    control,
    name: 'socialLinks',
  })

  const roleId = watch('roleId')
  const status = watch('status')
  const name = watch('name')
  const superAdmin = Boolean(watch('superAdmin'))
  const title = isEdit
    ? t('access.users.editTitle')
    : t('access.users.createTitle')
  const description = isEdit
    ? t('access.users.editDescription')
    : t('access.users.createDescription')

  useEffect(() => {
    if (!open) return
    if (isEdit && user) {
      setAvatarUrl(user.avatarUrl)
      reset({
        name: user.name,
        email: user.email,
        password: '',
        roleId: user.roleId,
        status: user.status,
        superAdmin: Boolean(user.superAdmin),
        phone: user.phone ?? '',
        general: user.general ?? '',
        birthday: user.birthday ?? '',
        address: user.address ?? '',
        socialLinks: (user.socialLinks ?? []).map((link) => ({
          label: link.label ?? '',
          url: link.url,
        })),
      })
      return
    }
    setAvatarUrl(undefined)
    reset({
      name: '',
      email: '',
      password: '',
      roleId: roles[0]?.id ?? '',
      status: 'active',
      superAdmin: false,
      ...emptyPersonalFields(),
    })
  }, [open, isEdit, user, roles, reset])

  return (
    <FormOverlay
      open={open}
      onOpenChange={onOpenChange}
      title={title}
      description={description}
      onSubmit={handleSubmit((values) => {
        const normalized = {
          ...values,
          phone: values.phone?.trim() ?? '',
          general: values.general?.trim() ?? '',
          birthday: values.birthday?.trim() ? values.birthday.trim() : '',
          address: values.address?.trim() ?? '',
          socialLinks: values.socialLinks.map((link) => ({
            label: link.label?.trim() || undefined,
            url: link.url.trim(),
          })),
        }
        if (isEdit) {
          onUpdate(normalized as UpdateUserFormValues)
        } else {
          onCreate(normalized as CreateUserFormValues)
        }
      })}
      footer={
        <UserFormFooter
          isEdit={isEdit}
          isPending={isPending}
          onCancel={() => onOpenChange(false)}
        />
      }
    >
      <UserFormFields
        isEdit={isEdit}
        userId={user?.id}
        userName={name || user?.name || user?.email || ''}
        avatarUrl={avatarUrl}
        roles={roles}
        roleId={roleId}
        status={status}
        superAdmin={superAdmin}
        canGrantSuperAdmin={isSuperAdmin}
        errors={errors}
        register={register}
        setValue={setValue}
        socialFields={socialFields}
        appendSocial={appendSocial}
        removeSocial={removeSocial}
        onAvatarUpdated={setAvatarUrl}
      />
    </FormOverlay>
  )
}
