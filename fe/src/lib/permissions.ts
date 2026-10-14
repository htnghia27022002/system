/**
 * Builds the `{resource}:view` / `{resource}:modify` keys of the RBAC model.
 * Shared by core (access-control) and module descriptors so neither imports the other.
 */
export function definePermissions<const R extends string>(resource: R) {
  return {
    view: `${resource}:view`,
    modify: `${resource}:modify`,
  } as const satisfies { view: `${R}:view`; modify: `${R}:modify` }
}
