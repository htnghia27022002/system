import { useQuery } from '@tanstack/react-query'

import { accessControlApi } from '../services/access-control-api'

const permissionsKey = ['admin', 'access-control', 'permissions'] as const

export function usePermissionsCatalog() {
  return useQuery({
    queryKey: permissionsKey,
    queryFn: () => accessControlApi.listPermissions(),
  })
}
