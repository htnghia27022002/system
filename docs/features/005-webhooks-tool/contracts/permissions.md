# Permissions contract: Tools → Webhooks

**Feature**: `005-webhooks-tool`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative RBAC keys for this feature. BE seeds + FE `PermissionKeys` / menu must match.  
> Model: `{resource}:view` and `{resource}:modify` only. See `.cursor/rules/feature-permissions.mdc`.

## Scope

- **Resource name**: `webhooks`
- **N/A**: not applicable — owner UI is under admin; keys are required for roles and sidebar

## Keys

| Key | Name | Group | Used by |
|-----|------|-------|---------|
| `webhooks:view` | View webhooks | `webhooks` | `RequireView`, admin page `/admin/tools/webhooks`, **sidebar menu** |
| `webhooks:modify` | Modify webhooks | `webhooks` | `RequireModify`, regenerate URL, soft-delete request |

## Admin menu

| Menu label (i18n key) | Href | Required permission |
|-----------------------|------|---------------------|
| `nav.tools` / `nav.webhooks` | `/admin/tools/webhooks` | `webhooks:view` |

Sidebar: **Tools** group with child **Webhooks** (filtered by `hasPermission(webhooks:view)`).

## Route ↔ permission map

| Surface | Permission |
|---------|------------|
| `GET /api/webhooks/inbox` | JWT + `webhooks:view` |
| `GET /api/webhooks/inbox/requests` | JWT + `webhooks:view` |
| `GET /api/webhooks/inbox/requests/:id` | JWT + `webhooks:view` |
| `POST /api/webhooks/inbox/regenerate` | JWT + `webhooks:modify` |
| `DELETE /api/webhooks/inbox/requests/:id` | JWT + `webhooks:modify` |
| `ANY /api/webhooks/capture/:uuid` | Public (no JWT / no RBAC) |
| FE `/admin/tools/webhooks` | `webhooks:view` (`PermissionGuard`) |
| FE regenerate / soft-delete actions | `webhooks:modify` (`PermissionGate`) |
| Product capture URL `/tools/webhooks/{uuid}` | Public rewrite → capture (unchanged) |

## Seed notes

- Add to `be/internal/database/seeders/catalog.go` → `DefaultPermissions()` with stable UUIDs:
  - `webhooks:view` → `10000001-0000-4000-8000-000000000017`
  - `webhooks:modify` → `10000001-0000-4000-8000-000000000018`
- Admin role receives all catalog keys via `RolePermissionSeeder`
- FE: extend `PermissionKeys` / `PermissionResource` with `webhooks`
- Owner UI path is **`/admin/tools/webhooks`** (not `/tools/webhooks` for the inbox shell)
