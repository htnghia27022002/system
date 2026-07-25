# Endpoints contract: [FEATURE NAME]

**Feature**: `NNN-short-name`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Draft | Ready for implement

> Authoritative HTTP API contract for this feature. Implementers follow this file; `plan.md` links here.  
> Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. JSON fields are **camelCase**.

## Scope

- **New routes**: …
- **Changed routes**: … (or none)
- **N/A**: … (only if this feature adds no HTTP API — keep this file and explain why)

## Conventions

- Prefix: `/api`
- Auth: JWT / public / RBAC permission key (`{resource}:view` or `{resource}:modify`) as noted per route
- Errors: document status codes that matter for clients
- Permissions catalog: also document keys in `contracts/permissions.md` (do not invent ad-hoc action names)

## Endpoints

### `METHOD /path`

| | |
|--|--|
| **Auth** | JWT required / public / permission `…` |
| **Request** | Query / body fields (camelCase) |
| **Success** | Status + response shape |
| **Errors** | e.g. 401, 403, 404, 413 |

**Success example**

```json
{
  "id": "uuid",
  "createdAt": "2026-01-01T00:00:00Z"
}
```

---

## Product / edge paths (optional)

| Public path | Proxies to | Notes |
|-------------|------------|-------|
| `/…` | `METHOD /api/…` | nginx / Next rewrite |

## UI routes (optional, FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/…` | Required / public | … |
