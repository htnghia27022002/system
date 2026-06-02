import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { RegisterFormValues } from '../schemas/auth-schemas'

export function useRegister() {
  const navigate = useNavigate()
  const signIn = useAuthStore((state) => state.signIn)

  return useMutation({
    mutationFn: (values: RegisterFormValues) => {
      const { name, email, password } = values
      return authApi.register({ name, email, password })
    },
    onSuccess: (data) => {
      signIn(data)
      toast.success('Account created successfully')
      void navigate(getPostLoginPath(data.user.role), { replace: true })
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
