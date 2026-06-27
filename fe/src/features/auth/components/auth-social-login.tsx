'use client'

import { useIsGoogleOAuthEnabled } from '../hooks/use-oauth-providers'

import { AuthDivider } from './auth-divider'
import { GoogleSignInButton } from './google-sign-in-button'

export function AuthSocialLogin() {
  const googleEnabled = useIsGoogleOAuthEnabled()

  if (!googleEnabled) {
    return null
  }

  return (
    <div className="grid gap-4">
      <AuthDivider />
      <GoogleSignInButton />
    </div>
  )
}
