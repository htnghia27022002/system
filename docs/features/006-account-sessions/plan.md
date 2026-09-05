# Implementation Plan: Account Sessions (Active Sign-ins on Account Settings)

**Feature**: `006-account-sessions` | **Date**: 2026-09-04 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Contracts** (authoritative): [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)

**Input**: Feature specification from `/docs/features/006-account-sessions/spec.md`

## Summary

Add owner-only **active session** management as the **fourth Card** on existing Account settings (`/admin/profile`). A session is an unrevoked, unexpired `refresh_tokens` row. Users see started / expires / last activity / IP / raw user-agent, identify **this device**, revoke other sessions one-by-one or all-at-once, and Sign out **actually revokes** the current refresh token.

**Backend** extends existing auth layers (`routes → handlers → services → repository`). No new domain package. Session queries live on **AuthRepository** (ISP: tokens stay here, not `UserRepository`). Reuse `be/pkg/repo`. Add `last_used_at` on `refresh_tokens`; populate existing `ip_address` / `user_agent` on issue and refresh. Return `sessionId` = `refresh_tokens.id` on login/register/OAuth/refresh.

**Frontend** extends `fe/src/features/user-profile/` with `account-sessions-card.tsx`, `use-account-sessions.ts`, `sessions-api.ts`. Store `sessionId` in `auth-token-service.ts`. Admin `user-menu-content.tsx` must use `useSignOut`.

P2 (friendly device names, geo, admin-of-another-user, sign-out-everywhere including this device, dedicated route) is deferred.

## Technical Context

**Language/Version**: Go 1.22 (BE); TypeScript strict on Next.js 15 App Router + React 19 (FE)

**Primary Dependencies**: Gin, pgx, squirrel, golang-migrate, PostgreSQL (BE); TanStack Query, shadcn/ui, react-i18next (FE)

**Storage**: PostgreSQL `refresh_tokens` (alter: `last_used_at`); existing `ip_address`, `user_agent`

**Testing**: `make test-be` / `make test-fe`; BE unit/integration for list/revoke/owner scope/current flag; FE Vitest for current-row UX / API mapping; QA manual Independent Tests

**Target Platform**: Docker stack (`make up-d`) and standalone BE/FE deploys; modern browsers; any signed-in user who can open `/admin/profile`

**Project Type**: Full-stack monorepo feature (`be/` + `fe/` independent packages)

**Performance Goals**: Sessions section interactive in under 5 seconds of navigation (SC-001); P1 lists all active sessions (no pagination)

**Constraints**: Package independence (HTTP only); camelCase JSON; English docs; JWT self-service (no new RBAC keys); never put refresh tokens in query strings; P1 only; no `device_id` usage

**Scale/Scope**: Typical session counts are small; one profile page; owner-only

## Constitution Check

*GATE: Must pass before design lock. Re-checked after design below.*

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | No new env; FE calls `/api/auth/*` only; no cross-imports |
| II. Role-owned artifacts | **Pass** | Architect owns `tasks.md` / `plan.md` / `contracts/*`; no app code in this phase |
| III. Spec before code | **Pass** | `spec.md` → `tasks.md` → `plan.md` (project override: tasks before plan) |
| IV. API contract alignment | **Pass** | camelCase under `/api`; **permissions N/A** — owner-only JWT like `/auth/profile`, not a new admin RBAC resource (explicit exception; do not invent `sessions:view`) |
| V. English documentation | **Pass** | Feature docs English; UI EN/VI via i18n |

**Post-design re-check**: Still pass. Session APIs sit next to existing `/api/auth/profile` (JWT middleware, no catalog keys, no sidebar). That matches 004-user-profile, not 005-webhooks-tool admin Tools RBAC.

## Complexity Tracking

> GitNexus CLI **was available** in this session. No constitution violations. HIGH blast radius is recorded so implementers do not treat AuthRepository edits as local.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

**GitNexus (upstream, 2026-09-04)**

| Symbol | Risk | Impacted | Notes |
|--------|------|----------|-------|
| `AuthRepository` (`be/internal/repository/interfaces/auth.go`) | **HIGH** | 36 | Additive methods only: `ListActiveByUserID`, `RevokeByID`, `RevokeAllExcept`. Do not move tokens to `UserRepository`. Update `MockAuthRepo`. |
| `issueTokenPair` (`be/internal/services/auth/auth_service.go`) | **HIGH** | 4 direct callers | Login, Register, OAuth `Callback`. Refresh has a sibling create path — update both. Additive: IP/UA/`last_used_at` + `sessionId` on response. Keep existing `accessToken` / `refreshToken` / `user` keys. |

Implementers **must** re-run `npx gitnexus impact AuthRepository --direction upstream` and `npx gitnexus impact issueTokenPair --direction upstream` before editing those symbols.

---

## Research Decisions (Phase 0)

### R1 — Session entity = `refresh_tokens` row

- **Decision**: An **active session** is a `refresh_tokens` row with `revoked_at IS NULL` and `expires_at > now()`. `sessionId` is `refresh_tokens.id` (UUID), never the raw refresh secret.
- **Rationale**: Spec maps 1:1 to stay-signed-in tokens from `001-auth`; no new sessions table.
- **Alternatives considered**: Separate `sessions` table — extra join for P1; rejected.

### R2 — `last_used_at` column; reuse IP/UA; no `device_id` in P1

- **Decision**: Migration `000008_refresh_token_last_used` adds `last_used_at TIMESTAMPTZ`, backfilled from `created_at`. Populate existing `ip_address` and `user_agent` on create and renew. Do **not** read or write `device_id`.
- **Rationale**: Locked P1; fingerprint IDs are out of scope.
- **Alternatives**: Geo/ASN columns — P2.

### R3 — Layers and ISP

- **Decision**: Keep `routes → handlers → services → repository`. No new domain package. Session methods on **AuthRepository** only. Reuse `be/pkg/repo`.
- **Rationale**: Tokens already live there; UserRepository must not grow token concerns.
- **Repo methods (locked)**: `ListActiveByUserID`, `RevokeByID` (scoped `user_id`), `RevokeAllExcept`. Existing `CreateRefreshToken` / `FindRefreshTokenByHash` / `RevokeRefreshToken` remain for issue/logout/refresh.

### R4 — Owner APIs under `/api/auth` (JWT, no RBAC)

- **Decision**:
  - `GET /api/auth/sessions`
  - `DELETE /api/auth/sessions/:id`
  - `POST /api/auth/sessions/revoke-others`
- Same JWT middleware group as `/me` and `/profile`. **No** `RequireView` / `RequireModify`. **No** `sessions:view` keys.
- **Rationale**: Spec FR-002; same pattern as 004 profile self-service.
- **Alternatives**: `/api/admin/sessions` + catalog keys — rejected (would imply admin-of-another-user and a sidebar item).

### R5 — Current session identity without secrets in URLs

- **Decision**:
  1. Every AuthResponse / TokenPairResponse includes `sessionId` (`refresh_tokens.id`).
  2. FE stores it in `auth-token-service.ts` (localStorage), never in the page address or query string.
  3. List: FE sends `X-Session-Id` (row uuid, not the refresh secret). Server sets `current: true` on the matching **owned active** row. If the header is absent, all `current` are false and the client compares stored `sessionId` to list `id`s.
  4. **Never** put the refresh token in query strings, path, or `X-Session-Id`.
- **Rationale**: FR-012; header is not a credential (knowing a uuid without the refresh secret does not continue a session).
- **Alternatives**: Encode session id in JWT access claims — larger blast radius on JWT manager; deferred.

### R6 — Refresh rotation vs stable “started” time

- **Decision**: Keep **refresh-secret rotation** (new raw token + hash). On Refresh, **copy `created_at` from the predecessor** onto the new row so “started” stays the original sign-in; set `last_used_at` / IP / UA on the new row; return the new `sessionId`. FE **must** persist `sessionId` on every refresh (`api-client.ts`).
- **Rationale**: Minimizes change to existing revoke-then-insert Refresh while satisfying “last activity updates on renew” and “started” ≠ last refresh.
- **Alternatives**: In-place update of the same row id — nicer identity stability, slightly more repo surface; allowed if implementers prefer, as long as `sessionId` remains the row id and last activity updates. Not required.

### R7 — Client IP helper

- **Decision**: Reuse `httpx.ClientIP` in `be/internal/common/httpx/httpx.go` (`X-Forwarded-For` first non-empty hop, else `X-Real-IP`, else `RemoteAddr`). User-Agent from the request header. Empty/missing → store empty and serialize JSON `null`.
- **Rationale**: Locked “reuse existing helpers; do not invent a second HTTP stack”.
- **Note**: Prefer passing IP/UA into `issueTokenPair` / `Refresh` from the handler, not importing Gin into the service.

### R8 — FE module placement

- **Decision**: Fourth Card on `profile-page.tsx`. New files under `fe/src/features/user-profile/`: `components/account-sessions-card.tsx`, `hooks/use-account-sessions.ts`, `services/sessions-api.ts`. kebab-case. TanStack Query. i18n EN+VI. Loading/empty/error. Confirm dialogs via existing shadcn `AlertDialog`. No Revoke on the current row.
- **Rationale**: Spec FR-001; feature-module rule; evolve 004 profile rather than a new feature folder.
- **Sign out**: `user-menu-content.tsx` must call `useSignOut` so `POST /api/auth/logout` revokes the current refresh token (body JSON `{ refreshToken }`, not query).

### R9 — Permissions N/A (mandatory explicit)

- **Decision**: [contracts/permissions.md](./contracts/permissions.md) exists with **N/A**. No catalog seed, no PermissionKeys, no sidebar, no PermissionGuard.
- **Rationale**: Owner-only JWT self-service, same class as `/auth/profile`. User lock: do not invent `sessions:view`.
- **Contrast**: 005-webhooks-tool later used admin Tools RBAC; **this feature must not copy that**.

### R10 — List order and pagination

- **Decision**: No pagination in P1. Order: **current first**, then others by `last_used_at` DESC. Service (or FE) may sort after `ListActiveByUserID`.
- **Rationale**: Spec assumptions; small N.

### R11 — Destructive API errors

- **Decision**: Single revoke of current id → **400** (use Sign out). Target missing / not owned / already inactive → **404** with a clear message (FE refreshes list). Revoke-others with invalid/missing current `sessionId` → **400**. Revoke-others with no other rows → **200** no-op (current remains).
- **Rationale**: FR-008, FR-019; no silent no-op that still shows a stale active row.

---

## Data Model (Phase 1 design)

> Authoritative schema: [contracts/database.md](./contracts/database.md). Summary below.

### RefreshToken (session)

| Field | DB | JSON (list) | Rules |
|-------|-----|-------------|-------|
| id | uuid PK | `id` / `sessionId` on auth payloads | Current-session identity |
| user_id | uuid FK → users | (not listed; scoped by JWT) | Owner-only |
| token_hash | varchar unique | **never in JSON list** | Secret hash only |
| expires_at | timestamptz | `expiresAt` | Active if > now and not revoked |
| revoked_at | timestamptz null | (omit from list; filter out) | Set on logout / revoke |
| created_at | timestamptz | `createdAt` | Session started; copied across refresh rotation |
| last_used_at | timestamptz | `lastUsedAt` | Set on create; updated on renew |
| ip_address | varchar(64) null | `ipAddress` | Best-effort; JSON `null` if unknown |
| user_agent | text null | `userAgent` | Raw UA in P1; JSON `null` if unknown |
| device_id | varchar | **not used in P1** | Do not populate or display |
| current | (derived) | `current` | True when `X-Session-Id` matches this id |

### State transitions

```text
Login/Register/OAuth --issueTokenPair--> RefreshToken(active, last_used_at=now, ip, ua)
Refresh --rotate secret--> new hash (+ copied created_at); last_used_at=now; new sessionId to client
GET /sessions --> active rows for JWT user
DELETE /sessions/:id --> revoked (if owned, not current)
POST /sessions/revoke-others --> revoke all except current
POST /auth/logout --> revoke current by refresh-token hash
Expiry --> omitted from list (no job required in P1)
```

### Entity relationships

```text
User (1) ---- (*) RefreshToken (active subset = sessions)
Owner UI --JWT GET/DELETE/POST /api/auth/sessions*--> AuthRepository
Sign out --JWT POST /api/auth/logout { refreshToken }--> RevokeRefreshToken
```

---

## Contracts

> **Authoritative sources:** [contracts/database.md](./contracts/database.md), [contracts/endpoints.md](./contracts/endpoints.md), [contracts/permissions.md](./contracts/permissions.md).  
> Summary below is a convenience mirror; prefer the contract files when they differ.

Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. All JSON **camelCase**.  
**Permissions:** N/A — JWT owner self-service, not admin RBAC — see [permissions.md](./contracts/permissions.md).

### JWT session APIs (new)

| Method | Path | Notes | Response |
|--------|------|-------|----------|
| GET | `/auth/sessions` | Optional `X-Session-Id`; owner active list | `{ items: SessionListItem[] }` |
| DELETE | `/auth/sessions/:id` | Owner; not current | 204 |
| POST | `/auth/sessions/revoke-others` | Body `{ sessionId }` | `{ items }` remaining active (typically current only) |

### Changed auth payloads

`AuthResponse` and `TokenPairResponse` gain `sessionId` (uuid). Existing keys unchanged.

### UI routes (FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/admin/profile` | Required (existing admin layout) | Fourth Card: Sessions |
| (none new) | — | No dedicated sessions path in P1 |

---

## Project Structure

### Documentation (this feature)

```text
docs/features/006-account-sessions/
├── spec.md
├── tasks.md
├── plan.md                 ← this file
├── contracts/
│   ├── database.md
│   ├── endpoints.md
│   └── permissions.md      ← N/A (JWT self-service, not admin RBAC)
├── be-tasks-verify.md      ← @be
├── fe-tasks-verify.md      ← @fe
└── qa-checklist.md         ← @qa
```

### Source Code (concrete paths)

```text
be/
├── migrations/000008_refresh_token_last_used.up.sql
├── migrations/000008_refresh_token_last_used.down.sql
├── internal/models/auth/auth.go              # LastUsedAt
├── internal/dto/auth/auth.go                 # sessionId + session DTOs
├── internal/repository/interfaces/auth.go    # ListActiveByUserID, RevokeByID, RevokeAllExcept
├── internal/repository/auth_repository.go
├── internal/services/auth/auth_service.go    # issueTokenPair, Refresh, List/Revoke
├── internal/services/auth/oauth_service.go   # pass IP/UA into issueTokenPair
├── public/handlers/auth_handler.go
├── public/routes/auth.go
├── test/testutil/mocks.go
├── test/unit/auth/session_service_test.go
└── test/integration/auth_sessions_integration_test.go

fe/
├── src/services/auth-token-service.ts        # store sessionId
├── src/services/api-client.ts                # persist sessionId on refresh
├── src/types/auth.ts                         # AuthTokens.sessionId
├── src/store/auth-store.ts
├── src/services/mock/auth.mock.ts
├── src/components/common/user-menu-content.tsx  # useSignOut
├── src/features/auth/services/auth-api.ts
├── src/features/user-profile/components/profile-page.tsx
├── src/features/user-profile/components/account-sessions-card.tsx
├── src/features/user-profile/hooks/use-account-sessions.ts
├── src/features/user-profile/services/sessions-api.ts
├── src/features/user-profile/index.ts
├── src/locales/en/admin.json
└── src/locales/vi/admin.json
```

**Structure Decision**: Full-stack feature on existing auth + user-profile modules; no new BE domain package; no new FE feature folder; no new App Router page.

---

## Quickstart validation (for implementers / QA)

### Prerequisites

- `make env` / configured `be/.env`, `fe/.env`, root `.env`
- `make up-d`

### Smoke scenarios

1. Guest opens `/admin/profile` → redirected to login; no session JSON leaked.
2. Sign in → `/admin/profile` fourth Card lists one highlighted session with started, expiry, last activity, IP (or unknown), raw UA.
3. Second browser, same account → two rows; each browser highlights only itself; Revoke hidden on current.
4. Revoke the other row (confirm) → other browser cannot stay signed in; this device remains.
5. Cancel confirm → list unchanged.
6. Three browsers → revoke all others → only this device remains.
7. Sign out from admin user menu → other remaining device no longer sees that session; signed-out device must log in again.
8. Stay signed in until refresh → `sessionId` updates in storage (not in URL); highlight still this device; `lastUsedAt` changes.
9. User B never sees user A’s rows.

### Commands

```bash
make test-be
make test-fe
```

---

## Env & deploy notes

No new env vars. IP/UA come from the incoming request (nginx already forwards `X-Real-IP` / `X-Forwarded-For` for `/api/`).

| Var | Owner file | Purpose |
|-----|------------|---------|
| (none new) | — | — |
| `CORS_ORIGINS` | `be/.env` | Existing credentialed auth APIs |
| `NEXT_PUBLIC_API_BASE_URL` | `fe/.env` | Existing `/api` base |

**Standalone BE**: `/api/auth/sessions*` on the JWT group.  
**Standalone FE**: same `/admin/profile` Card; needs reachable BE.

---

## Implementation order

1. **`@be`**: migration → model → AuthRepository methods → issueTokenPair/Refresh metadata + `sessionId` → session routes → unit/integration tests → `be-tasks-verify.md`
2. **`@fe`**: `sessionId` storage + refresh interceptor → API client/hook/card → profile fourth Card → confirm dialogs → `useSignOut` on admin menu → i18n → `fe-tasks-verify.md`
3. **`@qa`**: checklist + Independent Tests

P1 done when US1–US6 acceptance criteria met; US7 (P2) may remain open.
