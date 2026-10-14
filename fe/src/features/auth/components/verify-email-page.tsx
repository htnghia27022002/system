'use client'

import { CheckCircle2 } from 'lucide-react'
import { useSearchParams } from 'next/navigation'
import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'

import { AuthLayout } from '@/components/common/auth-layout'
import { TextLink } from '@/components/common/text-link'
import { Spinner } from '@/components/ui/spinner'
import { useAuthStore } from '@/store/auth-store'

import { useVerifyEmail } from '../hooks/use-account-flows'
import { InvalidLinkNotice } from './invalid-link-notice'

/** Opened from the verification email; works signed in or out. */
export function VerifyEmailPage() {
  const { t } = useTranslation('common')
  const token = useSearchParams().get('token') ?? ''
  const verify = useVerifyEmail()
  const signedIn = useAuthStore((s) => Boolean(s.accessToken))
  const syncSession = useAuthStore((s) => s.syncSession)
  const started = useRef(false)

  useEffect(() => {
    if (started.current || !token) return
    started.current = true
    verify.mutate(token, {
      onSuccess: () => {
        if (signedIn) void syncSession()
      },
    })
  }, [token, verify, signedIn, syncSession])

  return (
    <AuthLayout title={t('auth.verify.title')} description={t('auth.verify.description')}>
      {!token || verify.isError ? (
        <InvalidLinkNotice
          retryHref={signedIn ? '/admin/profile' : undefined}
          retryLabel={t('auth.verify.resendFromProfile')}
        />
      ) : verify.isSuccess ? (
        <div className="flex flex-col items-center gap-4 text-center" role="status">
          <CheckCircle2 className="size-10 text-success" aria-hidden />
          <p className="text-sm text-muted-foreground">{t('auth.verify.success')}</p>
          <TextLink href={signedIn ? '/admin' : '/login'} className="text-sm">
            {signedIn ? t('auth.verify.continue') : t('auth.links.backToSignIn')}
          </TextLink>
        </div>
      ) : (
        <div className="flex justify-center py-6" role="status" aria-live="polite">
          <Spinner className="size-6" />
          <span className="sr-only">{t('auth.verify.pending')}</span>
        </div>
      )}
    </AuthLayout>
  )
}
