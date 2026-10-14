import type { ComponentType } from 'react'

type Icon = ComponentType<{ className?: string }>

/** Admin sidebar entry contributed by a module (labels are `admin` i18n keys). */
export type ModuleNavItem = {
  titleKey: string
  href: string
  icon: Icon
  /** Permission required to see the entry (and its page). */
  permission: string
  /** Entries sharing a group key render under one collapsible parent. */
  group?: { key: string; titleKey: string; icon: Icon }
}

/** Breadcrumb for an admin route owned by a module (labels are `admin` i18n keys). */
export type ModuleBreadcrumb = {
  path: string
  labelKey: string
  parent?: { labelKey: string; href?: string }
}

/** Permission catalog row (mirrors the BE module's Permissions()); used by the mock API. */
export type ModulePermission = {
  key: string
  name: string
  group: string
  description: string
}

/** Public site link (labels are `common` i18n keys). */
export type ModulePublicLink = {
  href: string
  labelKey: string
}

/**
 * Descriptor an optional feature module exports from `features/<name>/module.ts`.
 * Core chrome (sidebar, breadcrumbs, public header, sitemap, mock RBAC) is composed
 * from the list in `src/config/modules.ts` — core never imports a module directly.
 */
export type AppModule = {
  id: string
  adminNav?: ModuleNavItem[]
  breadcrumbs?: ModuleBreadcrumb[]
  permissions?: ModulePermission[]
  publicLinks?: ModulePublicLink[]
  /** Public paths for sitemap.xml */
  sitemap?: string[]
}
