import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'

import { getPostLoginPath, useAuthStore } from '@/store/auth-store'

import { authApi, MockAuthError } from '../services/auth-api'
import type { LoginFormValues } from '../schemas/auth-schemas'

export function useLogin() {
  const navigate = useNavigate()
  const signIn = useAuthStore((state) => state.signIn)

  return useMutation({
    mutationFn: (values: LoginFormValues) => authApi.login(values),
    onSuccess: (data) => {
      signIn(data)
      toast.success('Signed in successfully')
      void navigate(getPostLoginPath(data.user.role), { replace: true })
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
