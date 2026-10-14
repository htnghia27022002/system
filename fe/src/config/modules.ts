/**
 * Enabled optional modules — the only place that lists them.
 *
 * Add a module: create `features/<name>/module.ts` exporting an `AppModule`, add its
 * `app/` routes, and append it below. Remove one: delete its feature folder, its
 * `app/` routes and its line here. A module that sets `backendModule` only shows its
 * admin entries when the BE lists that module in GET /api/app/info (BE MODULES env),
 * on top of the usual permission gating.
 */
// Module descriptors are the one allowed deep import into a feature: they are tiny,
// free of components, and safe to load from server files such as sitemap.ts.
import { toolsModule } from '@/features/tools/module'
import type {
  AppModule,
  ModuleBreadcrumb,
  ModuleNavItem,
  ModulePermission,
  ModulePublicLink,
} from '@/types/modules'

export const appModules: readonly AppModule[] = [toolsModule]

export function isModuleEnabled(id: string): boolean {
  return appModules.some((m) => m.id === id)
}

export const moduleNavItems: readonly ModuleNavItem[] = appModules.flatMap((m) => m.adminNav ?? [])

/** BE module names the FE knows about (used by the mock /app/info). */
export const backendModuleNames: readonly string[] = appModules.flatMap((m) =>
  m.backendModule ? [m.backendModule] : [],
)

/**
 * Admin nav entries of modules whose BE module is running. While the BE list is
 * unknown (undefined), modules that need a BE module are left out.
 */
export function navItemsForBackend(backend: ReadonlySet<string> | undefined): ModuleNavItem[] {
  return appModules
    .filter((m) => !m.backendModule || backend?.has(m.backendModule))
    .flatMap((m) => m.adminNav ?? [])
}
export const moduleBreadcrumbs: readonly ModuleBreadcrumb[] = appModules.flatMap(
  (m) => m.breadcrumbs ?? [],
)
export const modulePermissions: readonly ModulePermission[] = appModules.flatMap(
  (m) => m.permissions ?? [],
)
export const modulePublicLinks: readonly ModulePublicLink[] = appModules.flatMap(
  (m) => m.publicLinks ?? [],
)
export const moduleSitemapPaths: readonly string[] = appModules.flatMap((m) => m.sitemap ?? [])
