'use client'

import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'

import { useAuthStore } from '@/store/auth-store'

import { authApi } from '../services/auth-api'

export function useSignOut() {
  const router = useRouter()
  const signOut = useAuthStore((state) => state.signOut)

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSettled: () => {
      signOut()
      router.replace('/login')
    },
  })
}
