# Endpoints contract: Super Admin Flag

**Feature**: `007-super-admin`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Additive JSON only. No new routes. camelCase.

## Scope

- Existing `/api/auth/*` user object and `/api/admin/users*` payloads gain `superAdmin`
- Access token claims gain `superAdmin`
- No new HTTP paths

## Auth user object

Returned on login, register, refresh (nested `user`), OAuth callback, and `GET /api/auth/me`.

| Field | Type | Notes |
|-------|------|-------|
| *(existing fields unchanged)* | | |
| `superAdmin` | boolean | Always present. `true` = skip RBAC after sign-in |

## Access token claims

| Claim | Type | Notes |
|-------|------|-------|
| `superAdmin` | boolean | FE hydrate. Server Auth middleware re-reads DB flag when `UserRepository` is available |

`SignAccessToken` gains a `superAdmin` argument. Existing claims (`sub`, `email`, `name`, `role`, `roleId`, `permissions`) stay.

## Admin users

### `GET /api/admin/users` and `GET /api/admin/users/:id`

`UserResponse` adds:

| Field | Type | Notes |
|-------|------|-------|
| `superAdmin` | boolean | Always present |

### `POST /api/admin/users` and `PATCH /api/admin/users/:id`

Request may include optional `superAdmin` (boolean).

| Actor | Body includes `superAdmin` | Result |
|-------|----------------------------|--------|
| Super admin | omitted | Create: false. Update: unchanged |
| Super admin | true/false | Applied, unless it would remove the last Super admin |
| Not Super admin | omitted | Success (other fields as today) |
| Not Super admin | present | **403** |

### Last Super admin

Clearing `superAdmin`, setting `status` to `inactive`, or `DELETE /api/admin/users/:id` on the last Super admin → **403**.

Self delete / self deactivate remain **403** (existing).

## Permission middleware

All existing `RequireView` / `RequireModify` routes: if the authenticated user is Super admin → allow. Unauthenticated → 401 unchanged.

## Out of scope

- Admin list/revoke of another user’s `/api/auth/sessions*`
- New `/api/admin/super-admins` resource
