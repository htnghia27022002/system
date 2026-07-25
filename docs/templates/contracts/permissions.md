# Permissions contract: [FEATURE NAME]

**Feature**: `NNN-short-name`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Draft | Ready for implement

> Authoritative RBAC keys for this feature. BE seeds + FE `PermissionKeys` / menu must match.  
> Model: `{resource}:view` and `{resource}:modify` only. See `.cursor/rules/feature-permissions.mdc`.

## Scope

- **Resource name**: `…` (lowercase, matches group)
- **N/A**: … (only if no admin API, no admin menu, no PermissionGuard — explain why)

## Keys

| Key | Name | Group | Used by |
|-----|------|-------|---------|
| `{resource}:view` | View … | `{resource}` | `RequireView`, admin page, **sidebar menu** |
| `{resource}:modify` | Modify … | `{resource}` | `RequireModify`, create/edit/delete UI |

## Admin menu

| Menu label (i18n key) | Href | Required permission |
|-----------------------|------|---------------------|
| `nav.…` | `/admin/…` | `{resource}:view` |

If this feature has no admin nav entry, write **None**.

## Route ↔ permission map

| Surface | Permission |
|---------|------------|
| `GET /api/admin/…` | `{resource}:view` |
| `POST|PATCH|DELETE /api/admin/…` | `{resource}:modify` |
| FE `/admin/…` page | `{resource}:view` (`PermissionGuard`) |

## Seed notes

- Add to `be/internal/database/seeders/catalog.go` → `DefaultPermissions()` with stable UUIDs
- Admin role receives all catalog keys via `RolePermissionSeeder`
- FE: extend `PermissionKeys` / `PermissionResource` in `permission-keys.ts`
