import type { SearchEntityType, SearchHit } from '../types'

export type SearchNavigationAccess = {
  canViewUsers: boolean
  canViewRoles: boolean
  canViewPermissions: boolean
}

export function canAccessSearch(access: SearchNavigationAccess): boolean {
  return access.canViewUsers || access.canViewRoles || access.canViewPermissions
}

export function canNavigateToSearchHit(
  hit: SearchHit,
  access: SearchNavigationAccess,
): boolean {
  switch (hit.entityType) {
    case 'user':
      return access.canViewUsers
    case 'role':
      return access.canViewRoles
    case 'permission':
      // Permission hits deep-link into the roles list filter.
      return access.canViewRoles
    default:
      return false
  }
}

export function searchHitFilterParams(hit: SearchHit): URLSearchParams {
  const params = new URLSearchParams()

  switch (hit.entityType) {
    case 'user':
      if (hit.metadata?.email) {
        params.set('search', hit.metadata.email)
      } else {
        params.set('search', hit.title)
      }
      break
    case 'role':
      if (hit.metadata?.slug) {
        params.set('search', hit.metadata.slug)
      } else {
        params.set('search', hit.title)
      }
      break
    case 'permission':
      if (hit.metadata?.key) {
        params.set('permissionKey', hit.metadata.key)
      } else {
        params.set('permissionKey', hit.title)
      }
      break
  }

  return params
}

export function searchHitHref(
  hit: SearchHit,
  access?: SearchNavigationAccess,
): string | null {
  if (access && !canNavigateToSearchHit(hit, access)) {
    return null
  }

  const params = searchHitFilterParams(hit)

  switch (hit.entityType) {
    case 'user':
      return `/admin/users?${params.toString()}`
    case 'role':
    case 'permission':
      return `/admin/roles?${params.toString()}`
    default:
      return null
  }
}

export function entityTypeLabel(type: SearchEntityType): string {
  switch (type) {
    case 'user':
      return 'User'
    case 'role':
      return 'Role'
    case 'permission':
      return 'Permission'
    default:
      return type
  }
}
