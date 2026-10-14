'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import { FormOverlay } from '@/components/common/form-overlay'
import { Button } from '@/components/ui/button'
import { Field, FieldDescription, FieldError, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'

import {
  inviteUserSchema,
  type InviteUserFormValues,
} from '../schemas/access-control-schemas'
import type { Role } from '../types'

type UserInviteDialogProps = {
  open: boolean
  roles: Role[]
  isPending: boolean
  onOpenChange: (open: boolean) => void
  onInvite: (values: InviteUserFormValues) => void
}

/** Invite by email: the user sets their own password from the emailed link. */
export function UserInviteDialog({
  open,
  roles,
  isPending,
  onOpenChange,
  onInvite,
}: UserInviteDialogProps) {
  const { t } = useTranslation('admin')
  const {
    register,
    handleSubmit,
    control,
    reset,
    setValue,
    formState: { errors },
  } = useForm<InviteUserFormValues>({
    resolver: zodResolver(inviteUserSchema),
    defaultValues: { name: '', email: '', roleId: '' },
  })
  const roleId = useWatch({ control, name: 'roleId' })

  useEffect(() => {
    if (!open) return
    const member = roles.find((r) => r.slug === 'user') ?? roles[0]
    reset({ name: '', email: '', roleId: member?.id ?? '' })
  }, [open, roles, reset])

  return (
    <FormOverlay
      open={open}
      onOpenChange={onOpenChange}
      title={t('access.users.invite.title')}
      description={t('access.users.invite.description')}
      className="sm:max-w-lg"
      onSubmit={handleSubmit(onInvite)}
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
          <Button type="submit" disabled={isPending} className="w-full sm:w-auto">
            {isPending ? <Spinner className="size-4" /> : null}
            {t('access.users.invite.submit')}
          </Button>
        </>
      }
    >
      <FieldGroup>
        <Field data-invalid={Boolean(errors.email)}>
          <FieldLabel htmlFor="invite-email">{t('access.users.fields.email')}</FieldLabel>
          <Input
            id="invite-email"
            type="email"
            autoComplete="off"
            aria-invalid={Boolean(errors.email)}
            {...register('email')}
          />
          <FieldError errors={[errors.email]} />
        </Field>
        <Field data-invalid={Boolean(errors.name)}>
          <FieldLabel htmlFor="invite-name">{t('access.users.fields.name')}</FieldLabel>
          <Input id="invite-name" aria-invalid={Boolean(errors.name)} {...register('name')} />
          <FieldError errors={[errors.name]} />
        </Field>
        <Field data-invalid={Boolean(errors.roleId)}>
          <FieldLabel>{t('access.users.fields.role')}</FieldLabel>
          <Select
            value={roleId}
            onValueChange={(value) => setValue('roleId', value, { shouldValidate: true })}
          >
            <SelectTrigger className="w-full min-w-0" aria-label={t('access.users.fields.role')}>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {roles.map((role) => (
                <SelectItem key={role.id} value={role.id}>
                  {role.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <FieldDescription>{t('access.users.invite.hint')}</FieldDescription>
          <FieldError errors={[errors.roleId]} />
        </Field>
      </FieldGroup>
    </FormOverlay>
  )
}
