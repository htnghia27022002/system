'use client'

import { zodResolver } from '@hookform/resolvers/zod'
import { useRouter, useSearchParams } from 'next/navigation'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { AuthLayout } from '@/components/common/auth-layout'
import { InputError } from '@/components/common/input-error'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { apiErrorStatus } from '@/lib/api-error'

import { useAcceptInvite, useInvitation } from '../hooks/use-account-flows'
import {
  acceptInviteSchema,
  type AcceptInviteFormValues,
} from '../schemas/auth-schemas'
import { InvalidLinkNotice } from './invalid-link-notice'
import { PasswordFields } from './password-fields'

export function AcceptInvitePage() {
  const { t } = useTranslation('common')
  const router = useRouter()
  const token = useSearchParams().get('token') ?? ''
  const invitation = useInvitation(token)
  const accept = useAcceptInvite()
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AcceptInviteFormValues>({
    resolver: zodResolver(acceptInviteSchema),
    defaultValues: { name: '', password: '', confirmPassword: '' },
  })

  useEffect(() => {
    if (invitation.data) {
      reset({ name: invitation.data.name, password: '', confirmPassword: '' })
    }
  }, [invitation.data, reset])

  const onSubmit = handleSubmit((values) =>
    accept.mutate(
      { token, password: values.password, name: values.name },
      {
        onSuccess: () => {
          toast.success(t('auth.invite.success'))
          router.replace('/login')
        },
        onError: (error) => {
          if (apiErrorStatus(error) !== 400) toast.error(t('auth.invite.failed'))
        },
      },
    ),
  )

  const invalid = !token || invitation.isError || apiErrorStatus(accept.error) === 400

  return (
    <AuthLayout title={t('auth.invite.title')} description={t('auth.invite.description')}>
      {invalid ? (
        <InvalidLinkNotice />
      ) : invitation.isLoading ? (
        <div className="grid gap-4" aria-busy>
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-9 w-full" />
          <Skeleton className="h-9 w-full" />
        </div>
      ) : (
        <form className="flex flex-col gap-6" onSubmit={onSubmit} noValidate>
          <div className="grid gap-5">
            <div className="grid gap-2">
              <Label htmlFor="invite-email">{t('auth.fields.email')}</Label>
              <Input id="invite-email" value={invitation.data?.email ?? ''} readOnly disabled />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="name">{t('auth.fields.name')}</Label>
              <Input id="name" autoComplete="name" aria-invalid={Boolean(errors.name)} {...register('name')} />
              <InputError message={errors.name?.message} />
            </div>
            <PasswordFields register={register} errors={errors} autoFocus />
          </div>
          <Button type="submit" className="w-full" disabled={accept.isPending} aria-busy={accept.isPending}>
            {accept.isPending ? <Spinner /> : null}
            {t('auth.invite.submit')}
          </Button>
        </form>
      )}
    </AuthLayout>
  )
}
