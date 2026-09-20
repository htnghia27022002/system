# Permissions contract: Maps

**Feature**: `008-maps`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative RBAC keys for this feature. BE seeds + FE `PermissionKeys` / menu must match.  
> Model: `{resource}:view` and `{resource}:modify` only. See `.cursor/rules/feature-permissions.mdc`.

## Scope

- **Resource name**: `maps` (lowercase, matches group)
- **N/A**: not applicable — admin API, admin menu, and `PermissionGuard` are required

Super admin bypass from `007-super-admin` applies to every Maps gate **after sign-in**. It is not a catalog key. Guests remain `401`.

## Keys

| Key | Name | Group | Used by |
|-----|------|-------|---------|
| `maps:view` | View maps | `maps` | `RequireView("maps")`, `/admin/maps` page, **sidebar menu**, pin list/detail, Sources list, ingest status |
| `maps:modify` | Modify maps | `maps` | `RequireModify("maps")`, Sources create/edit/delete, Load data, Place status |

Do **not** add extra verbs (ingest, sources, pins). Sources and ingest share `maps`.

## Admin menu

| Menu label (i18n key) | Href | Required permission |
|-----------------------|------|---------------------|
| `nav.maps` | `/admin/maps` | `maps:view` |

Sidebar: top-level **Maps** item in the same chrome family as Users / Roles / Webhooks (same `useAdminNavItems` tree, filtered by `hasPermission`). Suggested placement: after Dashboard, before Tools.

- Set `permission: PermissionKeys.maps.view`.
- Filter with `hasPermission` so users without the key never see the link.
- Do **not** add a second menu item for Sources.

EN copy: `Maps`. VI copy (`fe/src/locales/vi/admin.json`): `Bản đồ` (UTF-8 with diacritics).

## Route ↔ permission map

| Surface | Permission |
|---------|------------|
| `GET /api/admin/maps/places` | JWT + `maps:view` |
| `GET /api/admin/maps/places/:id` | JWT + `maps:view` |
| `PATCH /api/admin/maps/places/:id` | JWT + `maps:modify` |
| `GET /api/admin/maps/categories` | JWT + `maps:view` |
| `GET /api/address/countries` | JWT |
| `GET /api/address/divisions` | JWT |
| `GET /api/admin/maps/sources` | JWT + `maps:view` |
| `GET /api/admin/maps/sources/:id` | JWT + `maps:view` |
| `POST /api/admin/maps/sources` | JWT + `maps:modify` |
| `PATCH /api/admin/maps/sources/:id` | JWT + `maps:modify` |
| `DELETE /api/admin/maps/sources/:id` | JWT + `maps:modify` |
| `POST /api/admin/maps/ingest` | JWT + `maps:modify` |
| `GET /api/admin/maps/ingest` | JWT + `maps:view` |
| `GET /api/admin/maps/ingest/:id` | JWT + `maps:view` |
| FE `/admin/maps` page | `maps:view` (`PermissionGuard`) |
| FE Sources create/edit/delete, Load data, Place status | `maps:modify` (`PermissionGate`) |
| FE sidebar Maps | `maps:view` (`hasPermission`) |

| Actor | Menu | View APIs | Sources write / Load data / status |
|-------|------|-----------|-------------------------------------|
| Guest | no | 401 | 401 |
| Signed in, no Maps keys, not Super admin | hidden | 403 | 403 |
| `maps:view` only | yes | 200 | 403 |
| `maps:view` + `maps:modify` | yes | 200 | allowed (validation applies) |
| Super admin (empty role permissions) | yes | 200 | allowed |
| Super admin, not signed in | — | 401 | 401 |

Clearing Super admin falls back to catalog keys on the next authorized request (same as `007-super-admin`).

## Seed notes

- Add to `be/internal/database/seeders/catalog.go` → `DefaultPermissions()` with stable UUIDs:
  - `maps:view` → `10000001-0000-4000-8000-000000000019`
  - `maps:modify` → `10000001-0000-4000-8000-000000000020`
- Use `rbac.Key("maps", rbac.ActionView)` / `rbac.Key("maps", rbac.ActionModify)`
- Suggested names/descriptions:
  - View maps — View the admin street map, Places, Sources, and ingest status
  - Modify maps — Manage Sources, run Load data, and change Place status
- Admin role receives all catalog keys via `RolePermissionSeeder` (iterates catalog). Demo member does **not** receive Maps keys unless assigned.
- FE: extend `PermissionKeys` / `PermissionResource` in `fe/src/features/access-control/permission-keys.ts` with `maps: { view, modify }`
- Update `permission-keys.test.ts` and any mock permission catalog
- GitNexus: `DefaultPermissions`, `PermissionKeys`, `useAdminNavItems` are HIGH — **additive only**
