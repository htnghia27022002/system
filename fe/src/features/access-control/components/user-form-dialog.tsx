'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { emptyPersonalFields } from '@/features/user-profile'
import { useIsMobile } from '@/hooks/use-is-mobile'

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
  const isMobile = useIsMobile()
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
  const title = isEdit
    ? t('access.users.editTitle')
    : t('access.users.createTitle')

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
      ...emptyPersonalFields(),
    })
  }, [open, isEdit, user, roles, reset])

  const onSubmit = handleSubmit((values) => {
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
  })

  const inputs = (
    <UserFormFields
      isEdit={isEdit}
      userId={user?.id}
      userName={name || user?.name || user?.email || ''}
      avatarUrl={avatarUrl}
      roles={roles}
      roleId={roleId}
      status={status}
      errors={errors}
      register={register}
      setValue={setValue}
      control={control}
      socialFields={socialFields}
      appendSocial={appendSocial}
      removeSocial={removeSocial}
      onAvatarUpdated={setAvatarUrl}
    />
  )

  const footer = (
    <UserFormFooter
      isEdit={isEdit}
      isPending={isPending}
      onCancel={() => onOpenChange(false)}
    />
  )

  if (isMobile) {
    return (
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent
          side="bottom"
          className="flex max-h-[min(92dvh,100%)] w-full max-w-none flex-col gap-0 overflow-x-hidden rounded-t-xl p-0"
        >
          <SheetHeader className="shrink-0 border-b px-4 py-4 pr-12 text-left">
            <SheetTitle>{title}</SheetTitle>
          </SheetHeader>
          <form
            className="flex min-h-0 min-w-0 flex-1 flex-col"
            onSubmit={onSubmit}
            noValidate
          >
            <div className="min-h-0 min-w-0 flex-1 overflow-x-hidden overflow-y-auto px-4 py-4">
              {inputs}
            </div>
            <div className="shrink-0 border-t bg-muted/30 px-4 py-4">
              {footer}
            </div>
          </form>
        </SheetContent>
      </Sheet>
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="grid max-h-[90vh] gap-0 overflow-hidden p-0 sm:max-w-2xl">
        <DialogHeader className="border-b px-6 py-4">
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <form
          className="flex max-h-[calc(90vh-4rem)] flex-col"
          onSubmit={onSubmit}
          noValidate
        >
          <div className="min-h-0 flex-1 overflow-y-auto px-6 py-4">
            {inputs}
          </div>
          <div className="border-t bg-muted/30 px-6 py-4">{footer}</div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
