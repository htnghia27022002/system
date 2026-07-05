'use client'

import { useMemo } from 'react'

import { PermissionKeys, usePermissions } from '@/features/access-control'

import {
  canAccessSearch,
  canNavigateToSearchHit,
  searchHitHref,
  type SearchNavigationAccess,
} from '../lib/search-routes'
import type { SearchHit } from '../types'

export function useSearchNavigation() {
  const { canView, hasPermission, sessionSynced } = usePermissions()

  const access = useMemo<SearchNavigationAccess>(
    () => ({
      canViewUsers: canView('users'),
      canViewRoles: canView('roles'),
      canViewPermissions: hasPermission(PermissionKeys.permissions.view),
    }),
    [canView, hasPermission],
  )

  const canUseSearch = canAccessSearch(access)

  function resolveHitHref(hit: SearchHit): string | null {
    return searchHitHref(hit, access)
  }

  function canNavigateToHit(hit: SearchHit): boolean {
    return canNavigateToSearchHit(hit, access)
  }

  return {
    access,
    canUseSearch,
    sessionSynced,
    resolveHitHref,
    canNavigateToHit,
  }
}
