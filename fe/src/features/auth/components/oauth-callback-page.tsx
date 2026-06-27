'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { AuthLayout } from '@/components/common/auth-layout'
import { TextLink } from '@/components/common/text-link'
import { Spinner } from '@/components/ui/spinner'
import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi } from '../services/auth-api'
import { getOAuthCallbackUri } from '../services/oauth'

export function OAuthCallbackPage() {
  const { t } = useTranslation('common')
  const searchParams = useSearchParams()
  const router = useRouter()
  const signIn = useAuthStore((state) => state.signIn)
  const [error, setError] = useState<string | null>(null)
  const started = useRef(false)

  useEffect(() => {
    if (started.current) {
      return
    }
    started.current = true

    const oauthError = searchParams.get('error')
    if (oauthError) {
      setError(t('auth.oauth.cancelled'))
      return
    }

    const code = searchParams.get('code')
    if (!code) {
      setError(t('auth.oauth.missingCode'))
      return
    }

    void (async () => {
      try {
        const data = await authApi.oauthCallback(
          'google',
          code,
          getOAuthCallbackUri(),
        )
        signIn(data)
        toast.success(t('auth.oauth.success'))
        router.replace(getPostLoginPath(data.user.role))
      } catch {
        setError(t('auth.oauth.failed'))
      }
    })()
  }, [router, searchParams, signIn, t])

  if (error) {
    return (
      <AuthLayout
        title={t('auth.oauth.callbackTitle')}
        description={t('auth.oauth.callbackErrorDescription')}
      >
        <div className="space-y-4 text-center">
          <p className="text-sm text-destructive" role="alert">
            {error}
          </p>
          <TextLink href="/login">{t('auth.actions.signIn')}</TextLink>
        </div>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout
      title={t('auth.oauth.callbackTitle')}
      description={t('auth.oauth.callbackDescription')}
    >
      <div className="flex justify-center py-6">
        <Spinner className="size-8" />
      </div>
    </AuthLayout>
  )
}
