import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import { Spinner } from '@/components/ui/spinner'
import { useIsMobile } from '@/hooks/use-is-mobile'

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
  const isMobile = useIsMobile()
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
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side={isMobile ? 'bottom' : 'right'}
        className={
          isMobile
            ? 'flex max-h-[min(92dvh,100%)] w-full flex-col gap-0 rounded-t-xl p-0'
            // The shadcn default applies data-[side=right]:sm:max-w-sm (384px) which has
            // higher specificity than plain sm:max-w-* utilities. Use ! to override both
            // the w-3/4 default and the sm:max-w-sm cap so the matrix can breathe.
            : 'flex h-full flex-col gap-0 p-0 !w-[min(860px,95vw)] !max-w-none'
        }
      >
        <SheetHeader className="shrink-0 border-b px-6 py-4">
          <SheetTitle>
            {isEdit
              ? t('access.roles.editTitle')
              : t('access.roles.createTitle')}
          </SheetTitle>
          <SheetDescription>
            {t('access.roles.permissionsSheetHint')}
          </SheetDescription>
        </SheetHeader>

        <form
          className="flex min-h-0 flex-1 flex-col"
          onSubmit={handleSubmit((values) => {
            if (isEdit) {
              onUpdate(values)
            } else {
              onCreate(values)
            }
          })}
          noValidate
        >
          <div className="min-h-0 flex-1 overflow-y-auto px-6 py-4">
            <div className="grid gap-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="grid gap-2">
                  <Label htmlFor="role-name">
                    {t('access.roles.fields.name')}
                  </Label>
                  <Input
                    id="role-name"
                    aria-invalid={Boolean(errors.name)}
                    {...register('name')}
                  />
                  <InputError message={errors.name?.message} />
                </div>

                <div className="grid gap-2">
                  <Label htmlFor="role-slug">
                    {t('access.roles.fields.slug')}
                  </Label>
                  <Input
                    id="role-slug"
                    aria-invalid={Boolean(errors.slug)}
                    {...register('slug')}
                  />
                  <InputError message={errors.slug?.message} />
                </div>
              </div>

              <PermissionPicker
                permissions={permissions}
                selectedKeys={selectedKeys}
                onChange={(keys) =>
                  setValue('permissionKeys', keys, { shouldValidate: true })
                }
                error={errors.permissionKeys?.message}
              />
            </div>
          </div>

          <SheetFooter className="shrink-0 flex-col-reverse gap-2 border-t bg-muted/30 px-4 py-4 sm:flex-row sm:justify-end md:px-6">
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
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  )
}
