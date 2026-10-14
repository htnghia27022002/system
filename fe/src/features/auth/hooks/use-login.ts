'use client'

import { useMutation } from '@tanstack/react-query'
import axios from 'axios'
import { useRouter } from 'next/navigation'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { LoginFormValues } from '../schemas/auth-schemas'

export function useLogin() {
  const router = useRouter()
  const signIn = useAuthStore((state) => state.signIn)
  const { t } = useTranslation('common')

  return useMutation({
    mutationFn: (values: LoginFormValues) => authApi.login(values),
    onSuccess: (data) => {
      signIn(data)
      toast.success('Signed in successfully')
      router.replace(getPostLoginPath())
    },
    onError: (error) => {
      if (error instanceof MockAuthError) {
        toast.error(error.message)
        return
      }
      // 429: per-IP rate limit or temporary account lockout on the BE.
      const tooMany =
        axios.isAxiosError(error) && error.response?.status === 429
      toast.error(
        t(tooMany ? 'auth.errors.tooManyAttempts' : 'auth.errors.invalidCredentials'),
      )
    },
  })
}
