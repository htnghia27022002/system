'use client'

import { useMutation } from '@tanstack/react-query'
import { MailWarning } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'

import { profileApi } from '../services/profile-api'

type EmailVerificationNoticeProps = {
  email: string
}

/** Shown on the profile while the account email is unverified. */
export function EmailVerificationNotice({ email }: EmailVerificationNoticeProps) {
  const { t } = useTranslation('admin')
  const resend = useMutation({
    mutationFn: () => profileApi.resendVerification(),
    onSuccess: () => toast.success(t('profile.verification.sent', { email })),
    onError: () => toast.error(t('profile.verification.failed')),
  })

  return (
    <Alert variant="warning" className="max-w-3xl">
      <MailWarning aria-hidden />
      <AlertTitle>{t('profile.verification.title')}</AlertTitle>
      <AlertDescription>
        <p>{t('profile.verification.description', { email })}</p>
        <Button
          size="sm"
          variant="outline"
          className="mt-2"
          disabled={resend.isPending || resend.isSuccess}
          onClick={() => resend.mutate()}
        >
          {resend.isPending ? <Spinner /> : null}
          {resend.isSuccess ? t('profile.verification.resent') : t('profile.verification.resend')}
        </Button>
      </AlertDescription>
    </Alert>
  )
}
