export const PermissionAction = {
  view: 'view',
  modify: 'modify',
} as const

export type PermissionResource = 'dashboard' | 'users' | 'roles' | 'permissions'

export const PermissionKeys = {
  dashboard: {
    view: 'dashboard:view',
  },
  users: {
    view: 'users:view',
    modify: 'users:modify',
  },
  roles: {
    view: 'roles:view',
    modify: 'roles:modify',
  },
  permissions: {
    view: 'permissions:view',
  },
} as const

const LEGACY_VIEW: Record<string, string[]> = {
  view: ['read'],
}

const LEGACY_MODIFY: Record<string, string[]> = {
  modify: ['create', 'update', 'delete'],
}

export function permissionKey(
  resource: PermissionResource,
  action: keyof typeof PermissionAction | string,
): string {
  return `${resource}:${action}`
}

export function isPermissionAllowed(
  granted: Iterable<string>,
  required: string,
): boolean {
  const set = granted instanceof Set ? granted : new Set(granted)
  if (set.has(required)) {
    return true
  }

  const [resource, action] = required.split(':')
  if (!resource || !action) {
    return false
  }

  if (action === PermissionAction.view) {
    for (const legacy of LEGACY_VIEW.view) {
      if (set.has(`${resource}:${legacy}`)) {
        return true
      }
    }
  }

  if (action === PermissionAction.modify) {
    for (const legacy of LEGACY_MODIFY.modify) {
      if (set.has(`${resource}:${legacy}`)) {
        return true
      }
    }
  }

  return false
}

export function canViewResource(
  granted: Iterable<string>,
  resource: PermissionResource,
): boolean {
  return isPermissionAllowed(granted, permissionKey(resource, PermissionAction.view))
}

export function canModifyResource(
  granted: Iterable<string>,
  resource: PermissionResource,
): boolean {
  return isPermissionAllowed(granted, permissionKey(resource, PermissionAction.modify))
}
