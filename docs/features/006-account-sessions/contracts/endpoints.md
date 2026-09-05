# Endpoints contract: Account Sessions

**Feature**: `006-account-sessions`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative HTTP API contract for this feature. Implementers follow this file; `plan.md` links here.  
> Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. JSON fields are **camelCase**.

## Scope

- **New routes**: JWT owner session APIs under `/api/auth/sessions*`
- **Changed routes**: login, register, OAuth callback, and refresh responses gain `sessionId`; issue/refresh persist IP, user-agent, and `last_used_at`
- **Unchanged but required for FR-015**: `POST /api/auth/logout` must continue to revoke the presented refresh token (JSON body, **not** query string)
- **N/A**: not applicable — feature adds HTTP APIs

## Conventions

- Prefix: `/api`
- Auth: JWT required on session routes (same `middleware.Auth` group as `/api/auth/me` and `/api/auth/profile`)
- RBAC: **none** — no `sessions:view` / `sessions:modify` (see [permissions.md](./permissions.md))
- Errors: document status codes that matter for clients
- **Never** put refresh tokens (or any session secret) in query strings, path params, or shareable URLs
- Session identity for “current”: header `X-Session-Id: <refresh_tokens.id UUID>` and/or JSON `sessionId` on POST body. This is the row id, **not** the raw refresh token
- List order: **current first**, then `lastUsedAt` descending
- No pagination in P1
- `ipAddress` / `userAgent`: JSON `null` when unknown — do not fabricate values

## Endpoints

### `GET /api/auth/sessions`

| | |
|--|--|
| **Auth** | JWT required (no RBAC key) |
| **Request** | No query. Optional header `X-Session-Id` (uuid of current refresh_tokens row) |
| **Success** | `200` — `{ items }` of **active** sessions for the JWT user only |
| **Errors** | `401` unauthorized |

**Success example**

```json
{
  "items": [
    {
      "id": "8c2a0c3e-1b2a-4d5e-9f10-111213141516",
      "createdAt": "2026-09-04T10:00:00Z",
      "expiresAt": "2026-09-18T10:00:00Z",
      "lastUsedAt": "2026-09-04T12:30:00Z",
      "ipAddress": "203.0.113.10",
      "userAgent": "Mozilla/5.0 …",
      "current": true
    },
    {
      "id": "9d3b1d4f-2c3b-5e6f-0a21-222324252627",
      "createdAt": "2026-09-03T08:00:00Z",
      "expiresAt": "2026-09-17T08:00:00Z",
      "lastUsedAt": "2026-09-03T09:00:00Z",
      "ipAddress": null,
      "userAgent": null,
      "current": false
    }
  ]
}
```

List item fields (locked): `id`, `createdAt`, `expiresAt`, `lastUsedAt`, `ipAddress`, `userAgent`, `current`.

`current` is `true` only when `X-Session-Id` matches that item’s `id` **and** the row belongs to the JWT user and is active. If the header is omitted, all items have `current: false` (client compares stored `sessionId`).

---

### `DELETE /api/auth/sessions/:id`

| | |
|--|--|
| **Auth** | JWT required (no RBAC key) |
| **Request** | Path `:id` = refresh_tokens UUID. Optional header `X-Session-Id` to identify current |
| **Success** | `204` No Content — row revoked; this device (current) stays signed in |
| **Errors** | `401`; `400` if `:id` is the current session; `404` if missing, not owned, already revoked, or expired |

Do not offer this for the current row in the UI. Ending this device uses Sign out (`POST /api/auth/logout`).

---

### `POST /api/auth/sessions/revoke-others`

| | |
|--|--|
| **Auth** | JWT required (no RBAC key) |
| **Request** | JSON body `{ "sessionId": "<uuid>" }` identifying the session to **keep**. May also accept `X-Session-Id` if body `sessionId` is empty. No refresh token in URL |
| **Success** | `200` — remaining active sessions (typically only current), e.g. `{ "items": [ … ] }` |
| **Errors** | `401`; `400` if `sessionId` missing, not owned, or not an active session of this user |

If the user has no other active sessions, return `200` with the current session still present (**no-op**, do not sign them out).

**Request example**

```json
{
  "sessionId": "8c2a0c3e-1b2a-4d5e-9f10-111213141516"
}
```

---

## Changed auth payloads

### `POST /api/auth/login`, `POST /api/auth/register`, `POST /api/auth/oauth/:provider/callback`

Add `sessionId` to the existing `AuthResponse`. Keep `accessToken`, `refreshToken`, `user`.

```json
{
  "accessToken": "…",
  "refreshToken": "…",
  "sessionId": "8c2a0c3e-1b2a-4d5e-9f10-111213141516",
  "user": { }
}
```

On success the new `refresh_tokens` row MUST store `ip_address`, `user_agent` (from `httpx.ClientIP` + `User-Agent`), and `last_used_at = now`.

### `POST /api/auth/refresh`

Add `sessionId` to `TokenPairResponse`:

```json
{
  "accessToken": "…",
  "refreshToken": "…",
  "sessionId": "8c2a0c3e-1b2a-4d5e-9f10-111213141516"
}
```

Refresh MUST rotate the refresh secret, record IP/UA/`last_used_at` on the successor row, and **copy `created_at` from the predecessor** so “started” remains the original sign-in. The client MUST persist the new `sessionId` (local storage, not URL).

### `POST /api/auth/logout` (existing — required behavior)

| | |
|--|--|
| **Auth** | Optional JWT; body carries the refresh token |
| **Request** | JSON `{ "refreshToken": "…" }` — **not** a query parameter |
| **Success** | `204` — current refresh token revoked when presented |
| **Errors** | Existing behavior; missing body/token may still clear client-side only — FE **must** send the token via `useSignOut` / `authApi.logout()` |

Admin chrome Sign out is not a new endpoint; it MUST call this API so the session disappears from other devices’ lists (FR-015).

---

## Product / edge paths (optional)

| Public path | Proxies to | Notes |
|-------------|------------|-------|
| (none) | — | No rewrite. Sessions are not a public product URL |

## UI routes (FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/admin/profile` | Required (existing `ProtectedGuard` + `AdminGuard`) | Fourth Card: Sessions. No new route in P1 |

## Client storage (not HTTP, but contract-adjacent)

- Store `sessionId` beside access/refresh tokens in `auth-token-service.ts`
- Compare to list `id`s to highlight current if the server did not mark `current`
- Clear `sessionId` on logout / token clear
- Refresh interceptor must save the new `sessionId` from `POST /auth/refresh`
