# Endpoints contract: Tools → Webhooks

**Feature**: `005-webhooks-tool`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement (as-shipped)

> Authoritative HTTP API contract for this feature. Implementers follow this file; `plan.md` links here.  
> Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. JSON fields are **camelCase**.  
> Aligned with shipped routes in `be/public/routes/webhook.go`, DTOs in `be/internal/dto/webhook/`, and FE client `fe/src/features/tools/services/webhooks-api.ts`.

## Scope

- **New routes**: owner JWT group under `/api/webhooks/inbox*`; public capture `ANY /api/webhooks/capture/:uuid`
- **Changed routes**: none
- **N/A**: not applicable — feature adds HTTP APIs

## Conventions

- Prefix: `/api`
- Auth: JWT via `middleware.Auth` for owner routes, plus RBAC `RequireView("webhooks")` / `RequireModify("webhooks")` (see `contracts/permissions.md`)
- Owner UI: `/admin/tools/webhooks` (FE). Public product capture path remains `/tools/webhooks/{uuid}` → BE capture
- Public capture: no JWT; CORS via `middleware.WebhookCaptureCORS` (`Access-Control-Allow-Origin: *`, no credentials)
- Body limit on capture: **1 MiB** (`MaxBodyBytes`); oversize → **413** after optional marker row
- List defaults: `page=1`, `limit=20` (max 100)
- Errors: document status codes that matter for clients

## Endpoints

### `GET /api/webhooks/inbox`

| | |
|--|--|
| **Auth** | JWT + `webhooks:view` |
| **Request** | None |
| **Success** | `200` — inbox; **get-or-create** on first call |
| **Errors** | `401` unauthorized; `403` missing permission |

**Success example**

```json
{
  "id": "uuid",
  "publicUuid": "uuid",
  "publicPath": "/tools/webhooks/uuid",
  "activeCount": 12,
  "lifetimeReceived": 40,
  "createdAt": "2026-07-25T00:00:00Z",
  "updatedAt": "2026-07-25T00:00:00Z"
}
```

---

### `POST /api/webhooks/inbox/regenerate`

| | |
|--|--|
| **Auth** | JWT + `webhooks:modify` |
| **Request** | Empty body |
| **Success** | `200` — updated `InboxResponse` with new `publicUuid` / `publicPath`; prior request history remains |
| **Errors** | `401`; `403` |

Old UUID returns **404** on capture after regenerate.

---

### `GET /api/webhooks/inbox/requests`

| | |
|--|--|
| **Auth** | JWT + `webhooks:view` |
| **Request** | Query: `method` (optional), `q` (optional search), `read` (`all`\|`read`\|`unread`, optional), `page` (default 1), `limit` (default 20, max 100) |
| **Success** | `200` — paginated **active** (non-soft-deleted) requests, newest first, plus counters |
| **Errors** | `401`; `403` |

**Success example**

```json
{
  "items": [
    {
      "id": "uuid",
      "method": "POST",
      "url": "/tools/webhooks/…?x=1",
      "clientIp": "203.0.113.10",
      "createdAt": "2026-07-25T00:00:00Z",
      "snippet": "{\"ok\":true}",
      "isRead": false
    }
  ],
  "activeCount": 12,
  "lifetimeReceived": 40,
  "page": 1,
  "limit": 20,
  "total": 12,
  "hasMore": false
}
```

---

### `GET /api/webhooks/inbox/requests/:id`

| | |
|--|--|
| **Auth** | JWT + `webhooks:view` |
| **Request** | Path `:id` = request UUID |
| **Success** | `200` — full detail; **marks the request as read** if it was unread |
| **Errors** | `401`; `403`; `404` if missing, other user’s inbox, or soft-deleted |

**Success example**

```json
{
  "id": "uuid",
  "inboxId": "uuid",
  "method": "POST",
  "url": "/api/webhooks/capture/…",
  "clientIp": "203.0.113.10",
  "headers": { "Content-Type": "application/json" },
  "query": {},
  "form": {},
  "body": "{\"a\":1}",
  "bodyEncoding": "utf-8",
  "isBinary": false,
  "contentType": "application/json",
  "bodyTruncated": false,
  "captureStatus": "ok",
  "isRead": true,
  "createdAt": "2026-07-25T00:00:00Z"
}
```

Notes:

- `bodyEncoding`: `utf-8` (text) or `base64` (binary); when binary, `isBinary` is `true` and `body` is base64.
- `captureStatus`: `ok` | `oversized` | `error`.

---

### `PATCH /api/webhooks/inbox/requests/:id/read`

| | |
|--|--|
| **Auth** | JWT + `webhooks:view` |
| **Request** | `{ "isRead": true \| false }` |
| **Success** | `200` — updated detail |
| **Errors** | `401`; `403`; `404` |

---

### `DELETE /api/webhooks/inbox/requests`

| | |
|--|--|
| **Auth** | JWT + `webhooks:modify` |
| **Request** | None (clears **all** active requests for the owner inbox) |
| **Success** | `200` — soft-delete all active; `activeCount` → `0`; lifetime unchanged |
| **Errors** | `401`; `403` |

---

### `DELETE /api/webhooks/inbox/requests/:id`

| | |
|--|--|
| **Auth** | JWT + `webhooks:modify` |
| **Request** | Path `:id` |
| **Success** | `200` — soft-delete; **idempotent** (already deleted / missing still returns `ok` with current counters) |
| **Errors** | `401`; `403` |

**Success example**

```json
{
  "ok": true,
  "activeCount": 11,
  "lifetimeReceived": 40
}
```

---

### `ANY /api/webhooks/capture/:uuid`

| | |
|--|--|
| **Auth** | Public (no JWT) |
| **Request** | Any method/body/headers/query; body ≤ 1 MiB |
| **Success** | `200` — fixed acknowledgment JSON; **HEAD** returns `200` with empty body |
| **Errors** | `404` unknown / regenerated UUID; `413` body > 1 MiB (marker row may be stored); `400` failed body read |

**Success example**

```json
{
  "ok": true,
  "message": "Request received"
}
```

CORS:

- Capture path allows `Access-Control-Allow-Origin: *` without credentials.
- Browser preflight (`OPTIONS` + `Access-Control-Request-Method`) → `204` short-circuit.
- Bare `OPTIONS` (e.g. curl) continues into Capture so the method can be recorded.

Client IP: first non-empty `X-Forwarded-For` entry → `X-Real-IP` → remote addr.

---

## Product / edge paths

| Public path | Proxies to | Notes |
|-------------|------------|-------|
| `ANY /tools/webhooks/{uuid}` | `ANY /api/webhooks/capture/{uuid}` | Docker: `docker/nginx/config/webhooks-capture.conf` (UUID regex, before FE catch-all). Local FE: `fe/next.config.ts` rewrite → `{API}/webhooks/capture/:uuid` |
| Direct BE | `ANY /api/webhooks/capture/{uuid}` | Standalone BE without nginx Tools-path rewrite |

Owner UI path `/tools/webhooks` (no UUID) is **not** proxied — served by Next.js.

## UI routes (FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/tools/webhooks` | JWT required (`ProtectedGuard`) | Owner inbox UI |
| `/tools/webhooks/{uuid}` | None | **Not** a Next page — nginx/Next rewrite → BE capture |
| `/tools` | Public catalog | Link to `/tools/webhooks` unchanged |
