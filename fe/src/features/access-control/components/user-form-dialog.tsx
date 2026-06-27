'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
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

  const {
    register,
    handleSubmit,
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
    },
  })

  const roleId = watch('roleId')
  const status = watch('status')
  const title = isEdit
    ? t('access.users.editTitle')
    : t('access.users.createTitle')

  useEffect(() => {
    if (!open) return
    if (isEdit && user) {
      reset({
        name: user.name,
        email: user.email,
        password: '',
        roleId: user.roleId,
        status: user.status,
      })
      return
    }
    reset({
      name: '',
      email: '',
      password: '',
      roleId: roles[0]?.id ?? '',
      status: 'active',
    })
  }, [open, isEdit, user, roles, reset])

  const onSubmit = handleSubmit((values) => {
    if (isEdit) {
      onUpdate(values as UpdateUserFormValues)
    } else {
      onCreate(values as CreateUserFormValues)
    }
  })

  const inputs = (
    <UserFormFields
      isEdit={isEdit}
      roles={roles}
      roleId={roleId}
      status={status}
      errors={errors}
      register={register}
      setValue={setValue}
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
          className="flex max-h-[min(92dvh,100%)] flex-col gap-0 rounded-t-xl p-0"
        >
          <SheetHeader className="shrink-0 border-b px-4 py-4 text-left">
            <SheetTitle>{title}</SheetTitle>
          </SheetHeader>
          <form
            className="flex min-h-0 flex-1 flex-col"
            onSubmit={onSubmit}
            noValidate
          >
            <div className="min-h-0 flex-1 overflow-y-auto px-4 py-4">
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
      <DialogContent className="grid max-h-[90vh] gap-0 overflow-hidden p-0 sm:max-w-xl">
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
