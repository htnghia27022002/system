# Permissions contract: Super Admin Flag

**Feature**: `007-super-admin`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Super admin is **not** an RBAC resource. Do **not** add catalog keys.

## Scope

- **Resource name**: N/A (flag on user)
- Existing Users keys stay: `users:view`, `users:modify`
- Super admin **bypasses** every existing `{resource}:view` / `{resource}:modify` check after sign-in
- Non–Super admin behavior is unchanged

## Keys

| Key | Name | Group | Used by |
|-----|------|-------|---------|
| N/A | — | — | Do **not** add catalog rows (`super-admin:view`, etc.) |

Do **not** extend `be/internal/database/seeders/catalog.go` → `DefaultPermissions()`.  
Do **not** add keys to `fe/src/features/access-control/permission-keys.ts`.

## Admin menu

| Menu label | Href | Required permission |
|------------|------|---------------------|
| None new | — | Super admin sees **all existing** items because FE `hasPermission` returns true |

No new item in the sidebar. Users stays under Access control → `/admin/users`.

## Route ↔ permission map

| Surface | Non–Super admin | Super admin |
|---------|-----------------|-------------|
| `GET/PATCH/POST/DELETE /api/admin/users*` | existing `users:view` / `users:modify` | allow if signed in |
| Other `/api/admin/*` and webhook inbox | existing view/modify | allow if signed in |
| `/api/auth/*` owner routes | JWT only (unchanged) | JWT only (unchanged) |
| FE `/admin/users` `PermissionGuard` | `users:view` | allow |
| FE Users create/edit/delete gates | `users:modify` | allow |
| Change `superAdmin` field | forbidden | allowed (last Super admin protected) |

## Seed notes

- Seeded admin user (`aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa` / `admin@example.com`) is Super admin
- Administrator **role** is not Super admin by itself
- Demo member is not Super admin
