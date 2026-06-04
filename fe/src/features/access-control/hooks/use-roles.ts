'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'

import {
  accessControlApi,
  MockAccessControlError,
} from '../services/access-control-api'
import type { CreateRoleInput, UpdateRoleInput } from '../types'

const rolesKey = ['admin', 'access-control', 'roles'] as const

export function useRolesList() {
  return useQuery({
    queryKey: rolesKey,
    queryFn: () => accessControlApi.listRoles(),
  })
}

export function useRoleMutations() {
  const queryClient = useQueryClient()

  const invalidate = () => {
    void queryClient.invalidateQueries({ queryKey: rolesKey })
  }

  const createRole = useMutation({
    mutationFn: (input: CreateRoleInput) => accessControlApi.createRole(input),
    onSuccess: () => {
      invalidate()
      toast.success('Role created')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not create role',
      )
    },
  })

  const updateRole = useMutation({
    mutationFn: ({ id, input }: { id: string; input: UpdateRoleInput }) =>
      accessControlApi.updateRole(id, input),
    onSuccess: () => {
      invalidate()
      toast.success('Role updated')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not update role',
      )
    },
  })

  const deleteRole = useMutation({
    mutationFn: (id: string) => accessControlApi.deleteRole(id),
    onSuccess: () => {
      invalidate()
      toast.success('Role deleted')
    },
    onError: (error) => {
      toast.error(
        error instanceof MockAccessControlError
          ? error.message
          : 'Could not delete role',
      )
    },
  })

  return { createRole, updateRole, deleteRole }
}
