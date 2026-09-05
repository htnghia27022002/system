# Tasks: Account Sessions (Active Sign-ins on Account Settings)

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Contracts**: [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks). Generated from **spec.md + locked technical design** (project override: tasks before plan). `.specify/scripts/bash/setup-tasks.sh` failed because `plan.md` was not present yet.

**Tests**: Spec Independent Tests plus locked design: BE unit/integration for list, revoke, owner scope, and current flag. FE tests for current-row UX and API mapping as practical. QA owns manual Independent Tests. TDD is not required.

**Organization**: Phases follow P1 user stories (US1–US6). US7 is P2 deferred — do not implement. Task labels use `[BE]` / `[FE]` / `[QA]` / `[BA]` per constitution; `[USn]` maps to spec user stories.

**GitNexus blast radius (design note — HIGH)**: Extending `AuthRepository` is **HIGH** (tokens, OAuth, user-service DI, mocks). Changing `issueTokenPair` is **HIGH** (Login, Register, OAuth callback). Prefer additive methods and extra `AuthResponse.sessionId` / IP-UA fields; do not move tokens onto `UserRepository`. Implementers must re-run `npx gitnexus impact AuthRepository` and `npx gitnexus impact issueTokenPair` before editing those symbols.

**P1 only**: No parsed device names, geo, admin-of-another-user, sign-out-everywhere-including-this-device, dedicated sessions route, or `device_id` usage.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[BE]` / `[FE]` / `[QA]` / `[BA]`**: Role ownership
- **`[P]`**: Parallelizable (different files, no incomplete dependency)
- **`[USn]`**: User story from spec (story phases only)
- Include exact file paths in descriptions

## Path Conventions

- Backend: `be/` (Go module) — layers `public/routes` → `public/handlers` → `internal/services` → `internal/repository`; reuse `be/pkg/repo`
- Frontend: `fe/src/` (App Router + `features/user-profile/`); kebab-case files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm migration number, evolve-in-place profile page, and P1 vs P2 boundary.

- [x] T001 [BE] Confirm next migration after `be/migrations/000007_webhook_request_is_read` is `be/migrations/000008_refresh_token_last_used.{up,down}.sql`; no new env vars (do not add sessions config to root `.env` or `be/.env`)
- [x] T002 [FE] [P] Confirm Sessions is a **fourth Card** on existing `fe/src/features/user-profile/components/profile-page.tsx` at route `fe/src/app/admin/profile/page.tsx`; **do not** add a dedicated sessions route or sidebar item
- [x] T003 [QA] [P] Note Independent Test matrix from spec US1–US6 (two/three browsers, revoke one, revoke-all, Sign out truthfulness, guest redirect, cross-account isolation) for later `qa-checklist.md`
- [ ] T004 [BA] [P] Record US7 (P2: parsed names, geo, admin viewing another user, sign-out-everywhere including this device, dedicated route) as **out of P1** — no implement tasks beyond this note

**Checkpoint**: Migration number, fourth-section placement, and P2 exclusion agreed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, AuthRepository session methods, token issuance metadata, `sessionId` on auth payloads, JWT session routes. **BLOCKS all user stories.**

**⚠️ CRITICAL**: No story UI until migration + list/revoke APIs + `sessionId` persistence contract exist.

- [x] T005 [BE] Add `be/migrations/000008_refresh_token_last_used.up.sql` (+ `.down.sql`): `ALTER TABLE refresh_tokens ADD COLUMN last_used_at TIMESTAMPTZ`; backfill existing rows from `created_at`; set NOT NULL after backfill if practical; **do not** add `device_id` usage; `ip_address` and `user_agent` already exist — populate them in application code
- [x] T006 [BE] Add `LastUsedAt` on `be/internal/models/auth/auth.go` (`RefreshToken`) with `db:"last_used_at"`; keep `IPAddress` / `UserAgent`; leave `DeviceID` unused in P1
- [x] T007 [BE] Extend `be/internal/repository/interfaces/auth.go` + `be/internal/repository/auth_repository.go` (reuse `be/pkg/repo` on `refresh_tokens`): `ListActiveByUserID` (not revoked, not expired; order later applied in service), `RevokeByID` (scoped by `user_id` + id), `RevokeAllExcept` (user_id + keep id). Tokens stay on **AuthRepository** — not `UserRepository`
- [x] T008 [BE] [P] Implement the new AuthRepository methods on `be/test/testutil/mocks.go` `MockAuthRepo` so existing unit tests compile
- [x] T009 [BE] [P] Add session DTOs in `be/internal/dto/auth/auth.go`: `sessionId` on `AuthResponse` and `TokenPairResponse`; list item (`id`, `createdAt`, `expiresAt`, `lastUsedAt`, `ipAddress`, `userAgent`, `current`); list wrapper `{ items }`; revoke-others body `{ sessionId }`; camelCase JSON; nullable `ipAddress` / `userAgent` when unknown
- [x] T010 [BE] Change `issueTokenPair` in `be/internal/services/auth/auth_service.go` to persist `ip_address`, `user_agent`, and `last_used_at` (now) on create and return `sessionId` = `refresh_tokens.id` on `AuthResponse` (**GitNexus HIGH** — additive fields only; keep existing token pair keys)
- [x] T011 [BE] Update `Refresh` in `be/internal/services/auth/auth_service.go` to record IP/UA/`last_used_at` on the renewed row, keep refresh-secret rotation, **copy `created_at` from the predecessor** so “started” stays the original sign-in, and return `sessionId` on `TokenPairResponse`
- [x] T012 [BE] Pass client IP via existing `be/internal/common/httpx/httpx.go` `ClientIP` (`X-Forwarded-For` / `X-Real-IP` / `RemoteAddr`) and `User-Agent` from `be/public/handlers/auth_handler.go` into Login, Register, Refresh, and OAuth callback (`be/internal/services/auth/oauth_service.go`); **do not** invent a second HTTP stack
- [x] T013 [BE] Add list/revoke/revoke-others on existing `be/internal/services/auth/auth_service.go` + `be/public/handlers/auth_handler.go` + JWT-protected routes in `be/public/routes/auth.go`: `GET /api/auth/sessions`, `DELETE /api/auth/sessions/:id`, `POST /api/auth/sessions/revoke-others`. JWT only — **no** `RequireView` / `RequireModify`. Owner scope by JWT `user_id`. Never put refresh tokens in query strings
- [x] T014 [BE] Mark `current` on list items when request identifies the session via `X-Session-Id` header (row uuid, **not** the refresh secret) in `be/public/handlers/auth_handler.go`; if header absent, `current` is false and the client compares stored `sessionId`
- [x] T015 [BE] Reject `DELETE /api/auth/sessions/:id` with **400** when `:id` equals the current session (`X-Session-Id`); **404** when the row is missing, not owned, already revoked, or expired (no silent success that still looks active)
- [x] T016 [FE] Persist `sessionId` in `fe/src/services/auth-token-service.ts` (localStorage, never URL/query); extend `AuthTokens` in `fe/src/types/auth.ts` (and thus `AuthResponse` in `fe/src/features/auth/types.ts`); update `setTokens` / `clearTokens`; store `sessionId` from login/register/OAuth in `fe/src/store/auth-store.ts` and from refresh in `fe/src/services/api-client.ts`

**Checkpoint**: Migration applies; login/register/OAuth/refresh return `sessionId` and write IP/UA/`last_used_at`; JWT session APIs callable; FE can store current session id without leaking refresh in URLs.

---

## Phase 3: User Story 1 — See active sessions on Account settings (Priority: P1) 🎯 MVP

**Goal**: Signed-in user opens `/admin/profile` and sees a fourth Sessions section listing only their own **active** sessions with started, expires, last activity, network address, and raw user-agent.

**Independent Test**: Sign in, open Account settings, confirm Sessions is the fourth section with the current sign-in listed; guests redirected; another account’s sessions never appear.

### Tests

- [x] T017 [P] [BE] [US1] Unit tests in `be/test/unit/auth/session_service_test.go` (or adjacent auth unit file): list returns only active rows for the JWT user; expired/revoked omitted; user B’s rows never appear
- [x] T018 [P] [FE] [US1] Vitest mapping tests in `fe/src/features/user-profile/services/sessions-api.test.ts` for camelCase list fields (`id`, `createdAt`, `expiresAt`, `lastUsedAt`, `ipAddress`, `userAgent`, `current`)

### Implementation

- [x] T019 [FE] [US1] Add types + API client `fe/src/features/user-profile/types.ts` (or session types in the same feature) and `fe/src/features/user-profile/services/sessions-api.ts` calling `GET /api/auth/sessions` with `X-Session-Id` from `authTokenService` (never refresh token in query)
- [x] T020 [FE] [US1] Add TanStack Query hook `fe/src/features/user-profile/hooks/use-account-sessions.ts` (query key factory; no `useEffect` + `useState` for server data)
- [x] T021 [FE] [US1] Implement `fe/src/features/user-profile/components/account-sessions-card.tsx` (kebab-case): loading, empty, error + retry; show started, expires, last activity, network address, raw user-agent; list order current first then last activity desc (sort client-side if API already marks `current`)
- [x] T022 [FE] [US1] Mount the card as the **fourth** Card in `fe/src/features/user-profile/components/profile-page.tsx` after avatar, personal information, and change password; export from `fe/src/features/user-profile/index.ts` if the public barrel should include it
- [x] T023 [P] [FE] [US1] Add i18n under `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json` (`profile.sections.sessions` and session strings); Vietnamese must use UTF-8 diacritics

**Checkpoint**: US1 Independent Test passes (FR-001–FR-006, FR-016, SC-001, SC-006, SC-007).

---

## Phase 4: User Story 2 — Identify this device (Priority: P1)

**Goal**: The current session is highlighted. Revoke is not offered on that row. Identity survives login and refresh without secrets in URLs.

**Independent Test**: Sign in on two browsers; each highlights only itself; after session renewal the highlight remains correct.

### Tests

- [x] T024 [P] [FE] [US2] Vitest in `fe/src/features/user-profile/components/account-sessions-card.test.tsx` (or hook test): current row highlighted; **no Revoke control** on the current row; other rows may show Revoke

### Implementation

- [x] T025 [FE] [US2] Ensure `sessionId` is saved after login, register, OAuth callback (`fe/src/features/auth/services/auth-api.ts` / `auth-store.ts`) and after refresh (`fe/src/services/api-client.ts`); mock `fe/src/services/mock/auth.mock.ts` returns a `sessionId`
- [x] T026 [FE] [US2] Highlight current row in `account-sessions-card.tsx` using server `current` and/or stored `sessionId` compared to list `id`s; never put refresh token or `sessionId` in the page address
- [x] T027 [P] [BE] [US2] Integration coverage in `be/test/integration/auth_sessions_integration_test.go`: after login, list marks current when `X-Session-Id` matches; after refresh, new `sessionId` identifies the same logical device when the client sends the updated header

**Checkpoint**: US2 Independent Test passes (FR-007, FR-008, FR-012, SC-002).

---

## Phase 5: User Story 3 — Revoke another session one at a time (Priority: P1)

**Goal**: User confirms Revoke on a **non-current** row; only that session ends; this device stays signed in.

**Independent Test**: Two devices; revoke B from A; B cannot continue; A remains signed in; B gone from the list.

### Implementation

- [x] T028 [BE] [US3] Complete `RevokeByID` owner-scope behavior in `be/internal/services/auth/auth_service.go`: success ends only that row; 404 for other user / already inactive; 400 if target is current; unit tests in `be/test/unit/auth/session_service_test.go`
- [x] T029 [FE] [US3] Revoke action on non-current rows in `account-sessions-card.tsx` using shadcn `AlertDialog` (`fe/src/components/ui/alert-dialog.tsx` — same pattern as webhooks); cancel leaves list unchanged; on success invalidate `use-account-sessions` query
- [x] T030 [FE] [US3] Show a clear error when revoke returns 404/stale and refetch the list in `use-account-sessions.ts` / card (FR-019)

**Checkpoint**: US3 Independent Test passes (FR-009, FR-011, FR-017, FR-019, SC-003).

---

## Phase 6: User Story 4 — Revoke all other sessions at once (Priority: P1)

**Goal**: One confirmed action ends every active session except the current one.

**Independent Test**: Three devices; from A revoke all others; A stays signed in; B and C cannot continue; list shows only A.

### Implementation

- [x] T031 [BE] [US4] Implement `RevokeAllExcept` in `be/internal/services/auth/auth_service.go` + `POST /api/auth/sessions/revoke-others` with body `{ sessionId }` (and/or `X-Session-Id`); success `200` `{ items }` of remaining active sessions; no-op 200 when only current is active (does **not** sign out this device); 400 if `sessionId` missing/not owned/not active; tests in `be/test/unit/auth/session_service_test.go` and `be/test/integration/auth_sessions_integration_test.go`
- [x] T032 [FE] [US4] “Revoke all other sessions” control in `account-sessions-card.tsx` with confirm `AlertDialog`; **unavailable or no-op** when no other sessions; this device stays signed in; i18n in `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json`

**Checkpoint**: US4 Independent Test passes (FR-010, FR-011, FR-017, SC-004).

---

## Phase 7: User Story 5 — Sign out of this device ends the current session (Priority: P1)

**Goal**: Admin chrome Sign out calls existing logout so the **current refresh token is revoked** on the server. No Revoke on the current row.

**Independent Test**: Two devices; Sign out on A via admin user menu; from B, A is gone from Sessions; A must sign in again.

### Implementation

- [x] T033 [FE] [US5] Change `fe/src/components/common/user-menu-content.tsx` to use `useSignOut` from `@/features/auth` (which calls `authApi.logout()` with the refresh token in **JSON body** via `fe/src/features/auth/hooks/use-sign-out.ts`); stop local-only `signOut()` that skips the server. **GitNexus HIGH** on this chrome component — keep menu structure; swap the action only
- [x] T034 [P] [BE] [US5] Confirm `Logout` in `be/internal/services/auth/auth_service.go` + `POST /api/auth/logout` still revokes the presented refresh token; add/extend a unit or integration assertion in `be/test/unit/auth/auth_service_test.go` or `be/test/integration/auth_integration_test.go` that the row is no longer active after logout

**Checkpoint**: US5 Independent Test passes (FR-008, FR-015, SC-005).

---

## Phase 8: User Story 6 — Record network address, user-agent, and last activity (Priority: P1)

**Goal**: Create/renew writes IP + UA; last activity updates on renew; missing values show unknown; no fingerprint IDs.

**Independent Test**: New sign-in row has IP and UA; after renewal last activity (and recorded IP/UA) update; list usable without `device_id`.

### Implementation

- [x] T035 [BE] [US6] Ensure create/renew paths never invent IP/UA (empty → JSON `null`); do not read or write `device_id` in P1; cover in `be/test/unit/auth/session_service_test.go` / auth service tests
- [x] T036 [FE] [US6] In `account-sessions-card.tsx`, show unknown/unavailable i18n for null IP or UA; long raw user-agent may wrap or truncate visually with full string available (title/expand); no fingerprint field

**Checkpoint**: US6 Independent Test passes (FR-013, FR-014, FR-018, SC-008).

---

## Phase 9: User Story 7 — Deferred session enhancements (Priority: P2)

**Goal**: Explicitly out of P1. Do not implement.

- [ ] T037 [BA] [US7] Leave P2 (parsed friendly names, geo, admin viewing another user’s sessions, sign-out-everywhere including this device, dedicated sessions route) **unimplemented**; `/admin/profile` remains the P1 place to manage own sessions

**Checkpoint**: P2 is documented as deferred, not filed as P1 defects.

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: RBAC N/A enforcement, i18n completeness, verify docs owned by implementers later.

- [x] T038 [BE] [P] Do **not** add `sessions:view` / `sessions:modify` to `be/internal/database/seeders/catalog.go`; do **not** wrap session routes in `RequireView` / `RequireModify` — JWT self-service like `/api/auth/profile` ([contracts/permissions.md](./contracts/permissions.md))
- [x] T039 [FE] [P] Do **not** add PermissionKeys, sidebar entries in `fe/src/components/common/app-sidebar.tsx`, or `PermissionGuard` for sessions; `/admin/profile` stays behind existing admin layout guards
- [x] T040 [QA] After `@be` / `@fe` verify docs exist, run Independent Tests US1–US6 on `http://system.local:8080/admin/profile` (Sign out from admin user menu, two-browser revoke, refresh identity)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Immediate
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS** all P1 stories
- **US1 (Phase 3)**: Depends on Phase 2 — MVP list
- **US2 (Phase 4)**: Depends on US1 UI + `sessionId` storage (T016/T025)
- **US3 (Phase 5)**: Depends on US1 list + US2 current-row rule
- **US4 (Phase 6)**: Depends on US2 current identity (revoke-others must keep current)
- **US5 (Phase 7)**: Can start after Phase 2 (logout already exists); should land with US1 so the list stays truthful
- **US6 (Phase 8)**: Mostly completed by T010–T012; FE unknown-state follows US1 card
- **US7 (Phase 9)**: No code
- **Polish (Phase 10)**: After desired P1 stories

### User Story Dependencies

- **US1**: After Foundational
- **US2**: After US1 card + sessionId persistence
- **US3**: After US2 (must not revoke current)
- **US4**: After US2
- **US5**: Independent of list UI but required for list truthfulness
- **US6**: Foundational BE + US1 display
- **US7**: Deferred

### Parallel Opportunities

- T002/T003/T004 after T001
- T008/T009 while T007 is in progress (interface first)
- T017/T018 tests after DTOs exist
- T023 i18n parallel with card
- T024/T027 with US2 UI
- T034 parallel with T033
- T038/T039 are no-op guards (can be checked anytime after routes exist)

---

## Parallel Example: User Story 1

```text
Task: "Unit tests in be/test/unit/auth/session_service_test.go — owner-only active list"
Task: "Vitest mapping tests in fe/src/features/user-profile/services/sessions-api.test.ts"
Task: "i18n keys in fe/src/locales/en/admin.json and fe/src/locales/vi/admin.json"
```

---

## Implementation Strategy

### MVP First (User Story 1)

1. Phase 1 Setup
2. Phase 2 Foundational (migration, AuthRepository methods, `sessionId`, JWT routes)
3. Phase 3 US1 list on `/admin/profile`
4. **STOP and VALIDATE** Independent Test US1

### Incremental Delivery

1. Setup + Foundational
2. US1 list → demo
3. US2 current highlight + sessionId
4. US3 revoke one
5. US4 revoke others
6. US5 Sign out uses `useSignOut`
7. US6 unknown-state polish
8. Skip US7 (P2)

### Role handoff

1. **`@be` first**: T001, T005–T015, T017, T027–T028, T031, T034–T035, T038 → `be-tasks-verify.md`
2. **`@fe` second** (needs `sessionId` + list/revoke APIs): T002, T016, T018–T026, T029–T033, T036, T039 → `fe-tasks-verify.md`
3. **`@qa`**: T003, T040 after verify docs

---

## Notes

- `[P]` = different files, no incomplete dependency
- No new domain package; no tokens on `UserRepository`
- Never put refresh tokens in query strings or shareable URLs
- English docs; UI copy EN + VI
- Do not invent `sessions:view` keys
