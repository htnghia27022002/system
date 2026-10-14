'use client'

import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { apiErrorCode } from '@/lib/api-error'
import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { RegisterFormValues } from '../schemas/auth-schemas'

export function useRegister() {
  const router = useRouter()
  const signIn = useAuthStore((state) => state.signIn)
  const { t } = useTranslation('common')

  return useMutation({
    mutationFn: (values: RegisterFormValues) => {
      const { name, email, password } = values
      return authApi.register({ name, email, password })
    },
    onSuccess: (data) => {
      signIn(data)
      toast.success('Account created successfully')
      router.replace(getPostLoginPath())
    },
    onError: (error) => {
      // AUTH_REQUIRE_EMAIL_VERIFICATION: the account exists, sign-in waits for the link.
      if (apiErrorCode(error) === 'registration_disabled') {
        toast.error(t('auth.register.closed'))
        return
      }
      if (apiErrorCode(error) === 'email_not_verified') {
        toast.success(t('auth.register.checkEmail'))
        router.replace('/login')
        return
      }
      toast.error(
        error instanceof MockAuthError
          ? error.message
          : 'Could not create account. Try again.',
      )
    },
  })
}
