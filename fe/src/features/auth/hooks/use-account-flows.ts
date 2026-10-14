'use client'

import { useMutation, useQuery } from '@tanstack/react-query'

import { accountApi } from '../services/account-api'

export function useForgotPassword() {
  return useMutation({
    mutationFn: (email: string) => accountApi.forgotPassword(email),
  })
}

export function useResetPassword() {
  return useMutation({
    mutationFn: ({ token, password }: { token: string; password: string }) =>
      accountApi.resetPassword(token, password),
  })
}

export function useVerifyEmail() {
  return useMutation({
    mutationFn: (token: string) => accountApi.verifyEmail(token),
  })
}

export function useInvitation(token: string) {
  return useQuery({
    queryKey: ['auth', 'invitation', token],
    queryFn: () => accountApi.inspectInvite(token),
    enabled: token.length > 0,
    retry: false,
  })
}

export function useAcceptInvite() {
  return useMutation({
    mutationFn: (input: { token: string; password: string; name?: string }) =>
      accountApi.acceptInvite(input.token, input.password, input.name),
  })
}
