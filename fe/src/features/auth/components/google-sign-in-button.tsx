'use client'

import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'

import { GoogleIcon } from '@/components/icons/google-icon'
import { getGoogleOAuthStartUrl } from '../services/oauth'

export function GoogleSignInButton() {
  const { t } = useTranslation('common')

  return (
    <Button
      type="button"
      variant="outline"
      className="w-full"
      onClick={() => {
        window.location.href = getGoogleOAuthStartUrl()
      }}
    >
      <GoogleIcon className="size-4" />
      {t('auth.oauth.continueWithGoogle')}
    </Button>
  )
}
