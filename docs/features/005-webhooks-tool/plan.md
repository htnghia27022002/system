# Implementation Plan: Tools → Webhooks (Per-Account Public URL + Request Inbox)

**Feature**: `005-webhooks-tool` | **Date**: 2026-07-25 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Contracts** (authoritative): [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)

**Input**: Feature specification from `/docs/features/005-webhooks-tool/spec.md`

## Summary

Replace the Phase-1 Tools shell at `/tools/webhooks` with a real per-account webhook inbox: one auto-created public capture URL, unauthenticated ingest of inbound HTTP requests, and an authenticated owner UI (list + detail, copy/regenerate URL, soft-delete, active + lifetime totals). Retention caps stored requests at 200; body capture max is 1 MiB.

**Backend** owns persistence and both API surfaces: JWT owner APIs under `/api/webhooks/*` and public capture at `ANY /api/webhooks/capture/:uuid`. **Frontend** evolves existing `WebhooksToolPage` for the owner inbox only. **Edge** maps the product path `/tools/webhooks/{uuid}` to BE capture so the public URL matches Tools IA without rendering another user’s inbox in Next.js.

P2 (live push, custom responses, export, rich binary viewers, team inboxes) is deferred.

## Technical Context

**Language/Version**: Go 1.22 (BE); TypeScript strict on Next.js 15 App Router + React 19 (FE)

**Primary Dependencies**: Gin, GORM, golang-migrate, PostgreSQL (BE); TanStack Query (or existing fetch patterns), shadcn/ui, react-i18next (FE); nginx reverse proxy (Docker)

**Storage**: PostgreSQL for `webhook_inboxes` + `webhook_requests` (JSON for headers/query/form; body as constrained text/bytea)

**Testing**: `make test-be` / `make test-fe`; BE unit tests for size limit, retention, regenerate, soft-delete counters; FE Vitest for URL helpers; QA manual Independent Tests

**Target Platform**: Docker stack (`make up-d`) and standalone BE/FE deploys; modern browsers; signed-in account owners + arbitrary HTTP clients for capture

**Project Type**: Full-stack monorepo feature (`be/` + `fe/` independent packages) + small nginx/Next rewrite for product URL shape

**Performance Goals**: Capture acknowledgment fast on LAN; inbox list/detail interactive after refresh; no WebSocket in P1; max 200 rows/inbox; body ≤ 1 MiB

**Constraints**: Package independence (HTTP only); camelCase JSON; English docs; public capture without JWT; owner APIs JWT-only; no second owner tool page; P2 deferred

**Scale/Scope**: One inbox per user account; single owner UI route; public UUID path; desktop-first inbox (mobile stack = P2)

## Constitution Check

*GATE: Must pass before design lock. Re-checked after design below.*

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | Capture + storage in `be/`; FE calls HTTP only; optional `NEXT_PUBLIC_*` for rewrite base in `fe/.env`; no app secrets in root `.env` |
| II. Role-owned artifacts | **Pass** | Architect owns `tasks.md` / `plan.md`; no app code in this phase |
| III. Spec before code | **Pass** | `spec.md` → `tasks.md` → `plan.md` before `@be` / `@fe` |
| IV. API contract alignment | **Pass** | Contracts below use camelCase under `/api` |
| V. English documentation | **Pass** | Feature docs English; UI EN/VI via i18n |

**Post-design re-check**: Still pass. Public product URL is an edge rewrite onto BE; standalone BE remains usable via `/api/webhooks/capture/:uuid`; standalone FE uses Next rewrite or documents API-absolute capture URL when no edge is present.

## Complexity Tracking

> No constitution violations requiring justification. Nginx/Next rewrite is wiring, not a third app package.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

---

## Research Decisions (Phase 0)

### R1 — Public URL shape vs capture ownership (critical)

- **Decision**:
  1. **Canonical product URL** shown to owners: `{origin}/tools/webhooks/{uuid}` (matches spec IA).
  2. **Canonical BE ingest**: `ANY /api/webhooks/capture/:uuid` — **no JWT**, returns fixed acknowledgment only (never inbox HTML/JSON list).
  3. **Docker/nginx**: regex location *before* FE catch-all proxies `/tools/webhooks/<uuid>` → BE capture (preserve method, query, body, `X-Real-IP` / `X-Forwarded-For`).
  4. **Local/split FE**: Next.js rewrite in `fe/next.config.ts` for UUID-shaped paths → BE capture derived from `NEXT_PUBLIC_API_BASE_URL` (optional dedicated public base only if rewrite needs absolute upstream).
  5. **Do not** add `fe/src/app/(public)/tools/webhooks/[uuid]/page.tsx` as owner UI or as HTML capture — would violate FR-015 and couple capture to FE availability.
- **Rationale**: Spec wants Tools-path URLs; capture must work without owner UI online; existing nginx already proxies `/api/` to BE; UUID regex avoids stealing `/tools/webhooks` owner page.
- **Alternatives considered**:
  - **API-only public URL** (`/api/webhooks/capture/{uuid}` in clipboard) — simpler ops; rejected for P1 because clarifications locked Tools-path IA (still available as direct BE URL for standalone BE clients).
  - **FE Route Handler proxy** — works without nginx but keeps FE in the capture path; weaker “FE offline” story; acceptable fallback only if rewrite unavailable.
  - **Separate `/hooks/{uuid}`** — cleaner separation; deferred unless rewrite proves impractical.

### R2 — Owner API placement

- **Decision**: Authenticated group under `/api/webhooks/...` with existing JWT middleware (same pattern as `/api/auth` protected routes), not under `/api/admin` (no admin permission required in P1).
- **Rationale**: One inbox per signed-in account; FR-013 owner-only; Tools product ≠ RBAC admin module.
- **Endpoints**:
  - `GET /api/webhooks/inbox` — get-or-create
  - `POST /api/webhooks/inbox/regenerate`
  - `GET /api/webhooks/inbox/requests` — list active
  - `GET /api/webhooks/inbox/requests/:id`
  - `DELETE /api/webhooks/inbox/requests/:id` — soft-delete

### R3 — One inbox auto-create

- **Decision**: Unique `user_id` on `webhook_inboxes`. First successful `GET /api/webhooks/inbox` for the JWT user creates inbox + `public_uuid` (UUID v4). Idempotent thereafter.
- **Rationale**: FR-001; no manual “create inbox” in P1.

### R4 — Soft-delete + dual counters

- **Decision**: `webhook_requests.soft_deleted_at` nullable. Denormalized on inbox: `lifetime_received`, `active_count`.
  - Capture → both ++ (then retention may hard-purge)
  - Soft-delete → `active_count--` only (idempotent)
  - Hard-purge oldest when row count &gt; 200 → delete row; if purged row was active, `active_count--`; always `lifetime_received--` for purged rows (lifetime tracks currently stored history, not “ever including purged”)
- **Rationale**: Matches clarifications; avoids expensive COUNT(*) on every chrome render; purge semantics align with SC “lifetime no longer includes purged”.
- **Note**: Soft-deleted rows **count toward the 200 cap** (FR-014).

### R5 — Body size and oversize policy

- **Decision**: Application `LimitReader` / size check at **1 MiB**. Prefer **HTTP 413** (or 400 with clear message) to client; **do not** store unbounded body. If feasible, insert a short error-marker request (method/URL/meta, empty/ truncated body, status flag) so owner sees the attempt; otherwise no silent “success” body.
- **Rationale**: FR-006; nginx already allows 20m — app must enforce 1 MiB.

### R6 — CORS for public capture

- **Decision**: Owner APIs keep existing credentialed CORS allow-list. Public capture path allows cross-origin requests **without credentials** (permissive Allow-Origin for capture only, or dedicated middleware), because senders may be browsers; most webhook senders are server-side (CORS irrelevant).
- **Rationale**: Avoid breaking third-party browser tests; do not open credentialed owner APIs.

### R7 — Live updates / custom responses

- **Decision**: P1 = manual refresh and/or short polling on owner UI. Fixed capture acknowledgment. Custom response editor and WebSocket = P2.
- **Rationale**: Spec phasing.

### R8 — FE module placement

- **Decision**: Keep feature under `fe/src/features/tools/` (evolve `webhooks-tool-page.tsx`, add inbox components/services/types). Route stays `fe/src/app/(public)/tools/webhooks/page.tsx` with auth guard. Catalog `href` remains `/tools/webhooks`.
- **Rationale**: Extends 003-landing-page Tools IA; FR-016; avoids a conflicting second tool page.

### R9 — Client IP

- **Decision**: Prefer first public-looking entry in `X-Forwarded-For`, else `X-Real-IP`, else remote addr. Nginx already sets these for `/api/`; capture location must set the same headers.
- **Rationale**: Spec “best effort behind proxies”.

### R10 — GitNexus / blast radius

GitNexus MCP was unavailable in this session; implementers **must** run impact before editing shared symbols:

| Area | Likely symbols / files | Risk note |
|------|------------------------|-----------|
| API bootstrap | `public.Run`, `Register*Routes` | Additive route registration |
| DI | `app.NewContainer`, dependency resolvers | New webhook wiring |
| FE tools shell | `WebhooksToolPage`, tools `content` / catalog | Replace placeholder; keep href |
| nginx | `be.conf` / new include + `nginx.*.conf` include order | Capture regex before `/` |
| Next config | `next.config.ts` rewrites | UUID path only |

---

## Data Model (Phase 1 design)

> Authoritative schema: [contracts/database.md](./contracts/database.md). Summary below.

### WebhookInbox

| Field | DB | JSON | Rules |
|-------|-----|------|--------|
| id | uuid PK | `id` | |
| user_id | uuid UNIQUE NOT NULL FK → users | `userId` | exactly one inbox per account |
| public_uuid | uuid UNIQUE NOT NULL | `publicUuid` | regenerated on demand |
| lifetime_received | int NOT NULL DEFAULT 0 | `lifetimeReceived` | stored history count |
| active_count | int NOT NULL DEFAULT 0 | `activeCount` | non-soft-deleted count |
| created_at / updated_at | timestamptz | `createdAt` / `updatedAt` | |

Derived (not necessarily stored): `publicPath` = `/tools/webhooks/{publicUuid}`; FE may also expose absolute `publicUrl` using browser origin.

### WebhookRequest

| Field | DB | JSON | Rules |
|-------|-----|------|--------|
| id | uuid PK | `id` | |
| inbox_id | uuid NOT NULL FK | `inboxId` | cascade delete with inbox/user per global rules |
| method | varchar | `method` | GET/POST/… |
| url | text | `url` | as received (path + query reconstruction ok) |
| client_ip | varchar | `clientIp` | best effort; may be empty |
| headers | jsonb | `headers` | map/array of pairs; size-reasonable |
| query | jsonb | `query` | |
| form | jsonb | `form` | when present |
| body | bytea/text | `body` / `bodyText` | ≤ 1 MiB; binary → safe preview flag |
| content_type | varchar | `contentType` | optional |
| body_truncated | boolean / status enum | `bodyTruncated` / `captureStatus` | error/truncated marker |
| soft_deleted_at | timestamptz NULL | `softDeletedAt` | null = active |
| created_at | timestamptz | `createdAt` | list sort desc |

### State transitions

```text
(no inbox) --GET /api/webhooks/inbox--> Inbox(publicUuid)
Capture(valid uuid) --> Request(active) ; counters++
Soft-delete --> Request(hidden) ; active_count--
Regenerate --> new publicUuid ; old uuid 404 on capture
Retention >200 --> hard-delete oldest rows ; adjust counters
```

### Entity relationships

```text
User (1) ---- (1) WebhookInbox
WebhookInbox (1) ---- (*) WebhookRequest
Public client --ANY /tools/webhooks/{uuid}--> nginx/Next --> BE capture
Owner UI --JWT /api/webhooks/*--> inbox management
```

---

## Contracts

> **Authoritative sources:** [contracts/database.md](./contracts/database.md), [contracts/endpoints.md](./contracts/endpoints.md), [contracts/permissions.md](./contracts/permissions.md).  
> Summary below is a convenience mirror; prefer the contract files when they differ.

Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. All JSON **camelCase**.  
**Permissions:** N/A for P1 (JWT owner tool, not admin RBAC) — see [permissions.md](./contracts/permissions.md).

### Owner (JWT required)

| Method | Path | Notes | Response |
|--------|------|-------|----------|
| GET | `/webhooks/inbox` | Get-or-create | `WebhookInboxResponse` |
| POST | `/webhooks/inbox/regenerate` | Confirm on FE | updated inbox |
| GET | `/webhooks/inbox/requests` | Query: `method`, `q`, `page`/`cursor`, `limit` | `{ items, activeCount, lifetimeReceived, page info }` |
| GET | `/webhooks/inbox/requests/:id` | 404 if missing/other user/soft-deleted (P1) | detail |
| DELETE | `/webhooks/inbox/requests/:id` | Soft-delete; idempotent | `{ ok, activeCount, lifetimeReceived }` |

**WebhookInboxResponse (example):**

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

**Request list item (example):**

```json
{
  "id": "uuid",
  "method": "POST",
  "url": "/tools/webhooks/…?x=1",
  "clientIp": "203.0.113.10",
  "createdAt": "2026-07-25T00:00:00Z",
  "snippet": "{\"ok\":true}"
}
```

### Public capture (no JWT)

| Method | Path | Notes | Response |
|--------|------|-------|----------|
| ANY | `/webhooks/capture/:uuid` | Also reached via `/tools/webhooks/:uuid` rewrite | Fixed ack on success; 404 unknown; 413/400 oversize |

Success acknowledgment (illustrative; keep stable in implement verify):

```json
{ "ok": true, "message": "Request received" }
```

HEAD may return same status without body. OPTIONS should succeed for CORS preflight where enabled.

### UI routes (FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/tools/webhooks` | Required | Owner inbox UI (evolved shell) |
| `/tools/webhooks/{uuid}` | None | **Not** a Next page — proxied to BE capture |
| `/tools` | Public catalog | Link to `/tools/webhooks` unchanged |

### Validation / errors

- Unknown UUID → 404, no cross-inbox leakage
- Oversize body → explicit error; no unbounded store
- Non-owner JWT → 401/403 on owner routes
- Soft-deleted detail → 404 for normal GET (no trash browser in P1)

---

## Project Structure

### Documentation (this feature)

```text
docs/features/005-webhooks-tool/
├── spec.md
├── tasks.md
├── plan.md                 ← this file
├── contracts/
│   ├── database.md
│   ├── endpoints.md
│   └── permissions.md      ← N/A (JWT owner tool, not admin RBAC)
├── be-tasks-verify.md
├── fe-tasks-verify.md
└── qa-checklist.md
```

### Source Code (concrete paths)

```text
be/
├── migrations/000006_webhooks_inbox.up.sql
├── migrations/000006_webhooks_inbox.down.sql
├── internal/models/webhook/
├── internal/dto/webhook/
├── internal/repository/interfaces/…webhook…
├── internal/repository/…webhook…
├── internal/services/webhook/
├── public/handlers/webhook_handler.go
├── public/routes/webhook.go
├── public/api.go                          # RegisterWebhookRoutes
├── internal/app/container.go
├── internal/app/dependency/
├── be/.env.example                        # comments only unless new vars needed
└── test/unit/…webhook…

docker/nginx/
├── config/be.conf                         # and/or new webhooks-capture.conf included before FE
├── nginx.dev.conf
└── nginx.prod.conf

fe/
├── next.config.ts                         # UUID rewrite → BE capture
├── src/config/env.ts                      # if rewrite base needs explicit public var
├── src/app/(public)/tools/webhooks/page.tsx
└── src/features/tools/
    ├── components/webhooks-tool-page.tsx  # evolve (remove placeholder)
    ├── components/webhooks-inbox.tsx      # NEW
    ├── components/webhooks-request-list.tsx
    ├── components/webhooks-request-detail.tsx
    ├── services/webhooks-api.ts
    ├── types/webhooks.ts
    ├── content.ts / catalog.ts            # real-tool copy
    └── index.ts
```

**Structure Decision**: Full-stack feature; BE owns capture + storage; FE owns owner inbox UX; nginx/Next provide product-shaped public URL without a FE capture page.

---

## Quickstart validation (for implementers / QA)

### Prerequisites

- `make env` / configured `be/.env`, `fe/.env`, root `.env`
- `make up-d` (preferred for nginx capture path)

### Smoke scenarios

1. Guest opens `/tools/webhooks` → redirected to login.
2. Sign in → `/tools/webhooks` shows public URL `/tools/webhooks/{uuid}`; Copy works; reload same UUID.
3. `curl -X POST "http://system.local:8080/tools/webhooks/{uuid}" -H 'Content-Type: application/json' -d '{"a":1}'` → ack; refresh inbox → request visible with headers/body.
4. Soft-delete → hidden; active decreases; lifetime unchanged (until purge).
5. Regenerate → old UUID 404; new UUID captures; history remains.
6. Body &gt; 1 MiB → client error; no silent full body stored.
7. Direct BE: `curl http://localhost:<be-port>/api/webhooks/capture/{uuid}` still works (standalone BE).

### Commands

```bash
make test-be
make test-fe
```

---

## Env & deploy notes

| Var | Owner file | Purpose |
|-----|------------|---------|
| `NEXT_PUBLIC_API_BASE_URL` | `fe/.env` | Owner API calls; Next rewrite upstream base |
| Optional `NEXT_PUBLIC_WEBHOOKS_PUBLIC_ORIGIN` | `fe/.env` | Only if display/rewrite origin must differ from window origin (avoid unless needed) |
| `CORS_ORIGINS` | `be/.env` | Owner APIs; capture may use separate permissive policy in code |
| `NGINX_HTTP_PORT` | root `.env` | Existing; public URL host:port for Docker demos |

Do **not** put JWT/OAuth/webhook secrets in root `.env`. Capture needs no new shared secret in P1 (UUID obscurity only; regenerate for rotation).

**Standalone BE**: publish `/api/webhooks/capture/:uuid` as the capture URL (or put any reverse proxy in front with the Tools-path rewrite).

**Standalone FE**: owner UI only; configure rewrite/API base to reachable BE; without rewrite, document copying API-absolute capture URL as fallback.

---

## Implementation order

1. **`@be`**: migration → models/DTOs/repos/service → owner + capture routes → unit tests → `be-tasks-verify.md`
2. **Edge early**: nginx UUID location (Docker) + FE `next.config.ts` rewrite
3. **`@fe`**: types/API client → auth-guarded evolve of `WebhooksToolPage` → list/detail → soft-delete/totals/regenerate/refresh → `fe-tasks-verify.md`
4. **`@qa`**: checklist + Independent Tests on `system.local` origin

P1 done when US1–US6 acceptance criteria met; US7 (P2) may remain open.
