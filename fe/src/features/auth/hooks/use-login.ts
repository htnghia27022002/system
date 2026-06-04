'use client'

import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { LoginFormValues } from '../schemas/auth-schemas'

export function useLogin() {
  const router = useRouter()
  const signIn = useAuthStore((state) => state.signIn)

  return useMutation({
    mutationFn: (values: LoginFormValues) => authApi.login(values),
    onSuccess: (data) => {
      signIn(data)
      toast.success('Signed in successfully')
      router.replace(getPostLoginPath(data.user.role))
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAuthError
          ? error.message
          : 'Invalid email or password',
      )
    },
  })
}
