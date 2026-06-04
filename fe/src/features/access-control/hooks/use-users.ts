'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import { useAuthStore } from '@/store/auth-store'

import {
  accessControlApi,
  MockAccessControlError,
} from '../services/access-control-api'
import type { CreateUserInput, ListUsersParams, UpdateUserInput } from '../types'

const usersKey = ['admin', 'access-control', 'users'] as const

export function useUsersList(params: ListUsersParams = {}) {
  return useQuery({
    queryKey: [...usersKey, params],
    queryFn: () => accessControlApi.listUsers(params),
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

  return { createUser, updateUser, deleteUser }
}
