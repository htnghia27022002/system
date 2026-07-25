# Tasks: Tools → Webhooks (Per-Account Public URL + Request Inbox)

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks)

**Tests**: Spec does not require TDD. Prefer BE unit tests for capture limits, retention purge, regenerate invalidation, and soft-delete counters. FE Vitest for URL builders / list filters where practical. QA owns manual Independent Tests.

**Organization**: Phases follow user stories (US1–US6 in P1; US7 = P2 deferred). Task labels use `[BE]` / `[FE]` / `[QA]` per constitution; `[USn]` maps to spec user stories. Evolve existing `WebhooksToolPage` — do not invent a second owner tool page.

**Architecture note (public URL)**: Canonical user-facing capture path is `/tools/webhooks/{uuid}`. Ingest is owned by BE (`ANY /api/webhooks/capture/:uuid`, no JWT). Edge wiring: nginx regex location + optional Next.js rewrite so capture never renders the owner inbox UI. See plan.md R1.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[BE]` / `[FE]` / `[QA]`**: Role ownership
- **`[P]`**: Parallelizable (different files, no incomplete dependency)
- **`[USn]`**: User story from spec (story phases only)
- Include exact file paths in descriptions

## Path Conventions

- Backend: `be/` (Go module)
- Frontend: `fe/src/` (App Router + `features/`)
- Edge: `docker/nginx/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm targets, env ownership, and that the existing Tools shell is the evolution point.

- [x] T001 [BE] Confirm next migration version after `be/migrations/000005_*` will be `000006_webhooks_inbox` and document any webhooks-related env (if needed) only in `be/.env.example` — never root `.env`
- [x] T002 [FE] [P] Confirm owner route remains `fe/src/app/(public)/tools/webhooks/page.tsx` rendering evolved `fe/src/features/tools/components/webhooks-tool-page.tsx` (or extracted components under `fe/src/features/tools/`); do **not** add a competing owner page
- [x] T003 [FE] [P] Confirm public capture must **not** be a Next owner UI page: no inbox at `fe/src/app/(public)/tools/webhooks/[uuid]/page.tsx`; capture is BE + rewrite only
- [x] T004 [QA] [P] Note Independent Test matrix from spec US1–US6 for later checklist (capture methods, soft-delete counters, regenerate, guest redirect) — see [qa-checklist.md](./qa-checklist.md)

**Checkpoint**: Migration number, FE evolve-in-place, and capture-vs-owner routing agreed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, domain layers, owner + public API surfaces, nginx/Next capture wiring. **BLOCKS all user stories.**

**⚠️ CRITICAL**: No story UI or capture work depends on unfinished migration + contracts.

- [x] T005 [BE] Add migration `be/migrations/000006_webhooks_inbox.up.sql` (+ `.down.sql`): `webhook_inboxes` (id, user_id UNIQUE, public_uuid UNIQUE, lifetime_received, active_count, timestamps) and `webhook_requests` (id, inbox_id, method, url, client_ip, headers JSONB, query JSONB, form JSONB, body BYTEA/TEXT, body_truncated/status, content_type, soft_deleted_at, created_at) with indexes for newest-first list and retention
- [x] T006 [BE] Add models under `be/internal/models/webhook/` (`inbox.go`, `request.go`) with camelCase JSON tags matching contracts
- [x] T007 [BE] [P] Add DTOs under `be/internal/dto/webhook/` for inbox summary, request list item, request detail, list query (method, q, page/cursor), regenerate response
- [x] T008 [BE] Add repository interfaces in `be/internal/repository/interfaces/` + GORM impls in `be/internal/repository/` for inbox (get-or-create by user, get by public UUID, regenerate UUID, update counters) and requests (insert, list active with filter/pagination, get by id+inbox, soft-delete, count/purge oldest beyond 200)
- [x] T009 [BE] Implement `be/internal/services/webhook/` service: get-or-create inbox, build public path `/tools/webhooks/{uuid}`, capture with 1 MiB body limit + retention purge, list/detail/soft-delete/regenerate; only owner user_id may manage
- [x] T010 [BE] Add handlers `be/public/handlers/webhook_handler.go` + routes `be/public/routes/webhook.go`: register in `be/public/api.go`; wire DI in `be/internal/app/dependency/` + `be/internal/app/container.go`
- [x] T011 [BE] Owner JWT routes (auth middleware): `GET /api/webhooks/inbox` (get-or-create), `POST /api/webhooks/inbox/regenerate`, `GET /api/webhooks/inbox/requests`, `GET /api/webhooks/inbox/requests/:id`, `DELETE /api/webhooks/inbox/requests/:id` (soft-delete)
- [x] T012 [BE] Public no-auth capture: `Any /api/webhooks/capture/:uuid` accepting GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS (and other methods when practical); return fixed acknowledgment; 404 for unknown/invalidated UUID; enforce body ≤ 1 MiB with explicit error (prefer 413) and optional error-marker row
- [x] T013 [BE] [P] Client IP best-effort from `X-Forwarded-For` / `X-Real-IP` / `RemoteAddr` in capture handler (document behind nginx)
- [x] T014 [BE] [P] Public capture CORS: allow unauthenticated cross-origin capture without credentials on `/api/webhooks/capture/*` (or equivalent) without weakening credentialed owner CORS defaults in `be/public/api.go`
- [x] T015 [BE] Ensure nginx `client_max_body_size` already ≥ 1 MiB; document capture body limit is application-enforced at 1 MiB in `be/.env.example` comments if useful (no new root `.env` keys)
- [x] T016 [BE] Add nginx location **before** FE catch-all in `docker/nginx/config/` (dev + prod includes): regex `^/tools/webhooks/([0-9a-fA-F-]{36})$` → proxy to `http://be_upstream/api/webhooks/capture/$1` preserving method, query, body, and forwarded IP headers
- [x] T017 [FE] [P] Add Next.js rewrite in `fe/next.config.ts` for local/split deploy: `/tools/webhooks/:uuid` (UUID-shaped) → BE capture URL derived from `NEXT_PUBLIC_API_BASE_URL` (or dedicated public base env in `fe/.env` / `fe/src/config/env.ts` if needed)
- [x] T018 [FE] [P] Add feature types + API client under `fe/src/features/tools/` (e.g. `types/webhooks.ts`, `services/webhooks-api.ts`) matching camelCase contracts; helper to build display public URL `{origin}/tools/webhooks/{uuid}`

**Checkpoint**: Migration applies; owner get-or-create + public capture callable; `/tools/webhooks/{uuid}` reaches BE via nginx (Docker) and Next rewrite (local FE).

---

## Phase 3: User Story 1 — Open Webhooks tool and receive a public URL (Priority: P1) 🎯 MVP

**Goal**: Signed-in user opening `/tools/webhooks` auto-creates one inbox and can copy the public URL; guests redirected to login.

**Independent Test**: Sign in → open `/tools/webhooks` → UUID URL shown and copyable; reload same URL; guest redirected.

### Implementation

- [x] T019 [FE] [US1] Protect owner tool page with auth redirect (reuse `ProtectedGuard` or equivalent tools-scoped guard) on `fe/src/app/(public)/tools/webhooks/page.tsx` / layout so guests go to login with return path
- [x] T020 [FE] [US1] Replace placeholder body in `fe/src/features/tools/components/webhooks-tool-page.tsx` with inbox chrome that calls `GET /api/webhooks/inbox` on load (get-or-create) and displays the public URL
- [x] T021 [FE] [P] [US1] Implement Copy URL control with clipboard + success feedback in tools webhooks components
- [x] T022 [FE] [P] [US1] Update `fe/src/features/tools/content.ts` / catalog copy so Webhooks is a real tool (not “placeholder / later phase”) for signed-in users; keep Tools hub link `/tools/webhooks` (FR-016)
- [x] T023 [FE] [P] [US1] Add i18n keys for inbox URL chrome under `fe/src/locales/en/` and `fe/src/locales/vi/`

**Checkpoint**: US1 Independent Test passes (FR-001–FR-003, FR-011, FR-016).

---

## Phase 4: User Story 2 — Capture inbound HTTP requests (Priority: P1)

**Goal**: Unauthenticated clients can hit the public URL; requests land in the owner inbox; invalid UUID / oversize body handled explicitly.

**Independent Test**: Send GET/POST (headers, query, form, JSON) to public URL; owner sees entries; bad UUID → not found; body > 1 MiB → error, no silent unbounded store.

### Implementation

- [x] T024 [BE] [US2] Complete capture persistence: method, full URL, timestamp, IP, headers, query, form fields when present, body ≤ 1 MiB; increment lifetime/active counters; run retention hard-purge oldest when stored count &gt; 200 in `be/internal/services/webhook/`
- [x] T025 [BE] [US2] Fixed success acknowledgment response (status + short body) for successful captures; distinct 404 for unknown UUID; explicit oversize error; never capture into another inbox
- [x] T026 [BE] [P] [US2] Unit tests under `be/test/unit/` for capture size limit, unknown UUID, retention purge at 201st request, and method acceptance
- [x] T027 [FE] [P] [US2] Manual-refresh affordance on owner inbox (button and/or short poll interval) so newly captured requests appear without live WebSocket (P1)

**Checkpoint**: US2 Independent Test passes (FR-004–FR-006, FR-014–FR-015, SC-002, SC-005 capture half).

---

## Phase 5: User Story 3 — Browse inbox list and inspect detail (Priority: P1)

**Goal**: webhook.site-inspired list + detail; newest-first; search/filter; pagination; empty state still shows public URL.

**Independent Test**: Capture mixed methods; list newest-first; select detail sections; search; paginate/load-more.

### Implementation

- [x] T028 [BE] [US3] Finalize list API: newest-first, filter by method + free-text over URL/path/snippet, pagination or cursor; **active only** (exclude soft-deleted) in `be/internal/services/webhook/` + handler
- [x] T029 [BE] [US3] Finalize detail API returning URL, metadata, headers, query, form, body/preview (binary-safe note when non-text)
- [x] T030 [FE] [US3] Build inbox list + detail layout components under `fe/src/features/tools/components/` (e.g. `webhooks-inbox.tsx`, `webhooks-request-list.tsx`, `webhooks-request-detail.tsx`) wired into `webhooks-tool-page.tsx`
- [x] T031 [FE] [US3] List row chrome: method badge, short id/path snippet, IP, timestamp; empty state with public URL + how to send first request
- [x] T032 [FE] [US3] Detail panes/sections: URL, metadata, headers, query, form, body/content with scrollable overflow
- [x] T033 [FE] [US3] Search/filter controls + pagination or load-more calling list API

**Checkpoint**: US3 Independent Test passes (FR-007, FR-008, SC-006).

---

## Phase 6: User Story 4 — Soft-delete a request (Priority: P1)

**Goal**: Owner soft-deletes from list; item hidden; not returned by normal search.

**Independent Test**: Soft-delete one of several; disappears from list; search excludes it.

### Implementation

- [x] T034 [BE] [US4] Soft-delete sets `soft_deleted_at`, decrements `active_count` only (not lifetime); idempotent if already deleted
- [x] T035 [FE] [US4] Add remove/soft-delete action on list/detail; refresh list selection and counts after success

**Checkpoint**: US4 Independent Test passes (FR-009, SC-003 visibility half).

---

## Phase 7: User Story 5 — Active and lifetime totals (Priority: P1)

**Goal**: UI shows active (visible) and lifetime received counts with correct semantics under soft-delete and hard-purge.

**Independent Test**: Capture N, soft-delete K → active = N−K, lifetime = N; after retention purge, both exclude purged rows.

### Implementation

- [x] T036 [BE] [US5] Ensure inbox payload always returns `activeCount` and `lifetimeReceived` consistent with service counters (purge adjusts both as needed)
- [x] T037 [FE] [US5] Display dual totals in inbox chrome (e.g. “N active · M lifetime”) from inbox API; update after capture refresh, soft-delete, regenerate

**Checkpoint**: US5 Independent Test passes (FR-010, SC-003 counters).

---

## Phase 8: User Story 6 — Regenerate public URL (Priority: P1)

**Goal**: New UUID; old URL stops capturing; history remains.

**Independent Test**: Note old URL → regenerate → old fails → new captures → prior requests still listed.

### Implementation

- [x] T038 [BE] [US6] `POST /api/webhooks/inbox/regenerate` issues new `public_uuid`; old UUID lookup fails; requests remain tied to inbox id
- [x] T039 [FE] [US6] Regenerate control with confirm dialog; update displayed URL; keep inbox list/history visible
- [x] T040 [BE] [P] [US6] Unit test: old UUID 404 after regenerate; history still listable for owner

**Checkpoint**: US6 Independent Test passes (FR-012, SC-004).

---

## Phase 9: User Story 7 — Responsive detail on smaller screens (Priority: P2) — Deferred

**Goal**: Stacked list/detail on narrow viewports without blocking P1 desktop layout.

**Independent Test**: Narrow viewport can reach detail and still use copy-URL + totals.

### Implementation (deferred — not required for P1 done)

- [ ] T041 [FE] [US7] (P2) Improve stacked list/detail navigation on small screens in webhooks inbox components while preserving copy URL + totals

**Checkpoint**: Optional; P1 may ship desktop-first.

---

## Phase 10: Polish & Cross-Cutting

**Purpose**: Verify reports, tests, docs hygiene, QA gate.

- [x] T042 [BE] [P] Additional unit tests under `be/test/unit/` for soft-delete counters and owner authorization (non-owner cannot list)
- [x] T043 [FE] [P] Vitest for public URL builder / list filter helpers under `fe/src/features/tools/` where practical
- [x] T044 [BE] Write `docs/features/005-webhooks-tool/be-tasks-verify.md` after `make test-be`
- [x] T045 [FE] Write `docs/features/005-webhooks-tool/fe-tasks-verify.md` after `make test-fe`
- [x] T046 [QA] Manual: US1–US6 Independent Tests + SC-001–SC-006 on `http://system.local:8080` (nginx capture path) — **PASS** 2026-07-25 re-check: product URL GET/POST → 200 ack after `\{36\}` + rewrite fallback fixes. See [qa-checklist.md](./qa-checklist.md)
- [x] T047 [QA] Manual: Guest `/tools/webhooks` → login; valid UUID capture without auth; owner cannot see another account’s inbox; HEAD/OPTIONS appear in list — **PASS** 2026-07-25 re-check: isolation Pass; bare OPTIONS captured+listed; CORS preflight 204. Browser UI optional/skipped. See [qa-checklist.md](./qa-checklist.md)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → story phases
- **US1** needs Phase 2 owner inbox API + FE guard
- **US2** needs Phase 2 public capture + nginx/rewrite; can validate before full inbox UI polish
- **US3** needs US1 shell + list/detail APIs (US2 for realistic data)
- **US4–US5** build on list UI + counters from inbox API
- **US6** after US1 URL chrome; can parallel US3–US5 once regenerate API exists
- **US7** deferred P2
- **QA** after BE/FE verify for P1 stories

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|--------|
| US1 | Phase 2 owner API + FE route | MVP: URL + copy |
| US2 | Phase 2 capture + edge rewrite | Can test with curl before rich UI |
| US3 | US1 + list/detail APIs | Needs captures for realism |
| US4 | US3 list | Soft-delete action |
| US5 | Inbox counters in API + chrome | Often lands with US1/US4 |
| US6 | US1 URL chrome + regenerate API | History must remain |
| US7 | P1 complete | Deferred |

### Parallel Opportunities

```text
After Phase 2:
  @be: finish capture/retention (US2) + list/detail (US3) + soft-delete/regenerate
  @fe: US1 shell + copy in parallel once GET /inbox exists
  @fe: US3 UI after list/detail contract stable
  Edge (T016) can parallel BE route registration
```

### Suggested MVP

1. Phase 1–2 foundational (including nginx/Next capture wiring)  
2. US1 (open tool → URL → copy)  
3. US2 (curl capture → appears after refresh)  
4. US3 list/detail → US4 soft-delete → US5 totals → US6 regenerate  
5. US7 optional P2  

---

## Implementation Strategy

### Recommended start order for agents

1. **`@be` first**: T001, T005–T015, T024–T026, T028–T029, T034, T036, T038, T040, T042 → `be-tasks-verify.md`
2. **Edge / FE config in parallel early**: T016 (nginx), T017–T018 (Next rewrite + client types)
3. **`@fe` after inbox + capture contracts stable**: T019–T023 (US1), T027, T030–T033, T035, T037, T039, T043 → `fe-tasks-verify.md`
4. **`@qa`**: T046–T047 after verify docs

### Notes

- `[P]` = different files, no incomplete dependency  
- Do not add a second owner page; evolve `WebhooksToolPage`  
- Do not implement P2 live push / custom responses / export in P1 tasks  
- JSON camelCase everywhere (`activeCount`, `lifetimeReceived`, `publicUuid`, `softDeletedAt`)  
- Public callers must never receive inbox list/detail JSON from capture URL
