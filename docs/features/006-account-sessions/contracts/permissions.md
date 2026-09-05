# Permissions contract: Account Sessions

**Feature**: `006-account-sessions`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative RBAC keys for this feature. BE seeds + FE `PermissionKeys` / menu must match.  
> Model: `{resource}:view` and `{resource}:modify` only. See `.cursor/rules/feature-permissions.mdc`.

## Scope

- **Resource name**: N/A
- **N/A**: Owner-only JWT self-service, same class as `/api/auth/profile` (`004-user-profile`). This is **not** a new admin RBAC resource.
  - APIs live under `/api/auth/sessions*` with existing JWT middleware only
  - No admin sidebar/menu entry
  - No `PermissionGuard` (Account settings stays on `/admin/profile` behind the existing admin layout)
  - Any signed-in user who can open Account settings can list and revoke **their own** sessions
  - **Do not** invent `sessions:view` or `sessions:modify`

This N/A is an explicit product lock (spec FR-002), not an omitted file.

## Keys

| Key | Name | Group | Used by |
|-----|------|-------|---------|
| N/A | — | — | Do **not** add catalog rows |

## Admin menu

| Menu label (i18n key) | Href | Required permission |
|-----------------------|------|---------------------|
| None | — | — |

No new item in `fe/src/components/common/app-sidebar.tsx`. Sessions are a section on existing Account settings (`userMenu.settings` → `/admin/profile`).

## Route ↔ permission map

| Surface | Permission |
|---------|------------|
| `GET /api/auth/sessions` | JWT only (no RBAC) |
| `DELETE /api/auth/sessions/:id` | JWT only (no RBAC); owner scope in service |
| `POST /api/auth/sessions/revoke-others` | JWT only (no RBAC); owner scope in service |
| `POST /api/auth/logout` | Existing; JWT optional + refresh token in JSON body |
| FE `/admin/profile` Sessions Card | Existing admin layout guards only — **no** `PermissionGuard` for sessions |
| FE sidebar | **None** — do not add a sessions nav item |

Owner isolation is **user_id from JWT**, not a permission key. Administrators see **their own** sessions only (P2 would be required for viewing another user).

## Seed notes

- Do **not** add rows to `be/internal/database/seeders/catalog.go` → `DefaultPermissions()`
- Do **not** extend FE `PermissionKeys` / `PermissionResource` in `fe/src/features/access-control/permission-keys.ts`
- Do **not** wrap session routes with `RequireView` / `RequireModify`
- Admin role seeder needs no change
