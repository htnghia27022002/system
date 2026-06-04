'use client'

import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { toast } from 'sonner'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { RegisterFormValues } from '../schemas/auth-schemas'

export function useRegister() {
  const router = useRouter()
  const signIn = useAuthStore((state) => state.signIn)

  return useMutation({
    mutationFn: (values: RegisterFormValues) => {
      const { name, email, password } = values
      return authApi.register({ name, email, password })
    },
    onSuccess: (data) => {
      signIn(data)
      toast.success('Account created successfully')
      router.replace(getPostLoginPath(data.user.role))
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAuthError
          ? error.message
          : 'Could not create account. Try again.',
      )
    },
  })
}
