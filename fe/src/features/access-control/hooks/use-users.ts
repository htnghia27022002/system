'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { apiErrorStatus } from '@/lib/api-error'
import { useAuthStore } from '@/store/auth-store'

import {
  accessControlApi,
  MockAccessControlError,
} from '../services/access-control-api'
import type {
  CreateUserInput,
  InviteUserInput,
  ListUsersParams,
  UpdateUserInput,
} from '../types'

const usersKey = ['admin', 'access-control', 'users'] as const

export function useUsersList(params: ListUsersParams = {}) {
  return useQuery({
    queryKey: [...usersKey, params],
    queryFn: () => accessControlApi.listUsers(params),
    meta: { skipNavLoading: true },
  })
}

export function useUserMutations() {
  const queryClient = useQueryClient()
  const sessionUserId = useAuthStore((s) => s.user?.id)

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: usersKey })
  }

  const createUser = useMutation({
    mutationFn: (input: CreateUserInput) => accessControlApi.createUser(input),
    onSuccess: () => {
      invalidate()
      toast.success('User created')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not create user',
      )
    },
  })

  const updateUser = useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateUserInput }) =>
      accessControlApi.updateUser(id, input, sessionUserId),
    onSuccess: () => {
      invalidate()
      toast.success('User updated')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not update user',
      )
    },
  })

  const deleteUser = useMutation({
    mutationFn: (id: string) =>
      accessControlApi.deleteUser(id, sessionUserId),
    onSuccess: () => {
      invalidate()
      toast.success('User deleted')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not delete user',
      )
    },
  })

  const inviteUser = useMutation({
    mutationFn: (input: InviteUserInput) => accessControlApi.inviteUser(input),
    onSuccess: (user) => {
      invalidate()
      toast.success(`Invitation sent to ${user.email}`)
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : apiErrorStatus(error) === 409
            ? 'A user with this email already exists'
            : 'Could not send the invitation',
      )
    },
  })

  const resendInvite = useMutation({
    mutationFn: (id: string) => accessControlApi.resendInvite(id),
    onSuccess: () => toast.success('Invitation sent again'),
    onError: () => toast.error('Could not resend the invitation'),
  })

  return { createUser, updateUser, deleteUser, inviteUser, resendInvite }
}
