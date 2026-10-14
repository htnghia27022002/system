---
description: "New features must define RBAC permissions; FE gates admin menu and routes by permission."
paths:
  - "docs/features/**/*"
  - "be/internal/database/seeders/**"
  - "be/public/routes/**"
  - "be/internal/middleware/**"
  - "fe/src/features/access-control/**"
  - "fe/src/components/common/app-sidebar.tsx"
  - ".claude/agents/**"
---

# Feature permissions (RBAC)

When adding a **new feature** that exposes protected admin/API surfaces, **generate permissions** and gate UI/API with them. Do not ship a new admin page or menu item without a matching permission key.

## Key format

Use the existing **view / modify** model only:

| Action | Key | Middleware / FE |
|--------|-----|-----------------|
| Read / list / open page | `{resource}:view` | BE `RequireView("resource")` · FE `PermissionKeys.{resource}.view` |
| Create / update / delete | `{resource}:modify` | BE `RequireModify("resource")` · FE `PermissionKeys.{resource}.modify` |

Helpers: BE `be/common/rbac` (`rbac.Key`, `rbac.ActionView`, `rbac.ActionModify`); FE `fe/src/features/access-control/permission-keys.ts`.

## Who does what

| Role | Must do |
|------|---------|
| `@agent-technical-architect` | Document keys in `contracts/permissions.md` (resource, `:view` / `:modify`, which routes/menu). Add `[BE]` / `[FE]` tasks for catalog + menu gate. |
| `@agent-be` | Core resource: add rows to `DefaultPermissions()` in `be/internal/database/seeders/catalog.go`; module resource: return them from the module's `Permissions()` (`be/internal/modules/<name>/module.go`). Stable UUIDs; duplicate keys/IDs abort seeding. Protect routes with `RequireView` / `RequireModify`. Admin role gets new keys via `RolePermissionSeeder` (iterates the merged catalog). |
| `@agent-fe` | Core resource: extend `PermissionKeys`. Module resource: `definePermissions` in the module's `permissions.ts` + catalog rows and `adminNav` in its `module.ts` (the mock catalog and sidebar read them from `src/config/modules.ts`). Gate admin routes with `PermissionGuard`. **Show sidebar/menu items only when `hasPermission(...)`**. Gate mutate actions with `PermissionGate`. |

## Admin menu (mandatory for FE)

Any new admin nav item — core items in `use-admin-nav-items.ts`, module items in the module's `adminNav` — **MUST**:

1. Set `permission: PermissionKeys.<resource>.view` (or the documented view key).
2. Be filtered with `hasPermission` so users without the key never see the link.
3. Match a page wrapped in `PermissionGuard` with the same key.

Do **not** add always-visible admin menu entries for permission-protected features.

## When N/A is allowed

Keep `contracts/permissions.md` and mark **N/A** only when the feature has:

- No JWT-protected admin or Tools API, **and**
- No admin sidebar/menu entry, **and**
- No `PermissionGuard` page

Examples: purely public marketing UI, public webhook **capture** URL only.

**Do not** mark N/A for owner Tools UIs — put them under `/admin/...`, seed `{resource}:view` / `{resource}:modify`, and gate menu + pages with those keys.

## Checklist (implement)

- [ ] `contracts/permissions.md` lists keys (or N/A)
- [ ] BE: `catalog.go` `DefaultPermissions()` (core) or module `Permissions()` updated
- [ ] BE: routes use `RequireView` / `RequireModify`
- [ ] FE: `PermissionKeys` (core) or module `permissions.ts` + `module.ts` updated
- [ ] FE: sidebar filtered by `hasPermission`
- [ ] FE: page `PermissionGuard` + modify `PermissionGate` where needed
