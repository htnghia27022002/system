/**
 * Enabled optional modules — the only place that lists them.
 *
 * Add a module: create `features/<name>/module.ts` exporting an `AppModule`, add its
 * `app/` routes, and append it below. Remove one: delete its feature folder, its
 * `app/` routes and its line here. Keep this list in step with the BE registry
 * (`be/internal/modules/registry.go`); the sidebar is also permission-gated, so a
 * module disabled on the BE (no permissions granted) stays hidden.
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
