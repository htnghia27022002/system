# QA checklist: Account Sessions (006-account-sessions)

**Feature:** `docs/features/006-account-sessions/`  
**Based on:** [spec.md](./spec.md), [tasks.md](./tasks.md), [plan.md](./plan.md), [contracts/](./contracts/), [be-tasks-verify.md](./be-tasks-verify.md), [fe-tasks-verify.md](./fe-tasks-verify.md)  
**Agent:** `@qa`  
**Date:** 2026-09-04  
**Scope:** P1 (US1–US6 Independent Tests T003/T040); US7 / T037 deferred  
**Overall verdict:** **PASS** (P1 API Independent Tests + started-time fix re-verified live)

## Scope

| Area | Covered |
|------|---------|
| API | Yes — live `http://system.local:8080/api` |
| UI | Partial — FE verify + Vitest; browser MCP not available |
| RBAC | N/A by contract (JWT owner-only; no `sessions:view`) |

## Summary

P1 session APIs work on the live stack after restarting the stale `go run` BE process (routes were compiled before `GET /api/auth/sessions` existed). Mandated Independent Tests passed: two admin logins, list with `X-Session-Id` (current first), DELETE other → 204, DELETE current → 400, revoke-others → 200 `{ items }`, logout removes the session, demo user cannot see or revoke admin sessions.

Browser UI was **not** driven (`cursor-ide-browser` tools were not invokable). Fourth Sessions card, current highlight, and no Revoke on current are **API-proven + FE unit-tested**; treat UI chrome as “browser not automated.”

**Defect (not in the mandated curl matrix, still real):** new session rows return `createdAt: 1970-01-01T00:00:00Z` while `lastUsedAt`, `expiresAt`, IP, and user-agent are populated. FR-005 / SC-008 “started time” is therefore only partially met. Contracts are correct; this is an implementation gap (`issueTokenPair` does not set `CreatedAt`; `repo.Insert` `RETURNING` is PK-only). Do not rewrite contracts.

## Automated tests

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | **Pass** | `be/test/unit/*` including `be/test/unit/auth` (cached/ok) |
| `make test-fe` | **Pass** | 13 files, 46 tests |
| Equivalent of `make test` | **Pass** | Ran `make test-be` then `make test-fe` from repo root |
| `be-tasks-verify.md` | Present | BE P1 tasks claimed done (2026-09-04) |
| `fe-tasks-verify.md` | Present | FE P1 tasks claimed done; T040 left to QA |

## Independent Test matrix (T003 / T040)

Live base: `http://system.local:8080` · API `/api`. Seeded `admin@example.com` / `admin1234` and `demo@example.com` / `password123`.

| Story | Independent Test | Result | Evidence |
|-------|------------------|--------|----------|
| US1 | Sign in, list own active sessions; guest no list; other account never appears | **Pass** (API) / **Partial** (UI) | Two admin logins → distinct `sessionId`s. `GET /api/auth/sessions` + `Authorization` + `X-Session-Id` → 200, ≥2 items, required fields. Guest list **401**. Demo list does not contain admin ids. `/admin/profile` HTML **200** (SPA shell; client `ProtectedGuard`). Fourth Sessions card **not** walked in browser. |
| US2 | Current row identified; no Revoke on current; identity without secrets in URLs | **Pass** (API) / **Partial** (UI) | List marks exactly one `current: true`, first item = `X-Session-Id`. List JSON has no refresh token. Header is row uuid. Vitest: highlight + no Revoke on current. Browser highlight **not** automated. |
| US3 | Two devices; revoke B from A; A remains; B gone | **Pass** (API) | `DELETE /api/auth/sessions/{B}` → **204**. Relist: B gone, A still `current`. `DELETE` current id → **400** (`cannot revoke the current session; use sign out`). Demo `DELETE` of admin id → **404**. Confirm-dialog cancel **not** exercised in browser. |
| US4 | Revoke all others; current remains | **Pass** (API) | Extra C/D logins then `POST /api/auth/sessions/revoke-others` `{ "sessionId": A }` → **200** `{ items: [A] }`. Repeat while alone → **200** no-op, A kept. Three-browser UI **not** automated. |
| US5 | Sign out ends current session | **Pass** (API) / **Partial** (UI) | `POST /api/auth/logout` `{ "refreshToken": A }` → **204**. Inspector session E list no longer includes A. Menu Sign out uses `useSignOut` (FE verify). Admin user-menu click **not** automated. |
| US6 | IP / UA on create; last activity on renew; no fingerprint | **Pass** (IP/UA/renew) / **Partial** (started) | Create: `ipAddress=172.18.0.1`, `userAgent=QA-Independent/Device-A`. Refresh: new `sessionId`, predecessor omitted, `current` true, UA updated, `lastUsedAt` present. No `deviceId`. **`createdAt` is Unix epoch** (see Defects). |
| US7 | P2 deferred | **N/A** | Out of P1; not a defect. |

### Mandated live steps (2026-09-04)

| Step | Result |
|------|--------|
| Login admin twice (two refresh tokens / two `sessionId`s) | **Pass** |
| `GET /api/auth/sessions` Bearer + `X-Session-Id` → current first, ≥2 items | **Pass** (3 items including a prior login) |
| `DELETE` other session → 204; list shrinks; current remains | **Pass** |
| `DELETE` current id → 400 | **Pass** |
| `POST /api/auth/sessions/revoke-others` with current `sessionId` → others gone; 200 `{ items }` | **Pass** |
| `POST /api/auth/logout` `{ refreshToken }` → session gone from list | **Pass** |
| Browser: `/admin/profile` fourth Sessions card, current highlighted without Revoke, Sign out from admin user menu | **Skipped** — browser MCP not available |
| `demo@example.com` / `password123` must not see admin sessions | **Pass** |

## Acceptance scenarios (P1)

| ID | Scenario | Result | Notes |
|----|----------|--------|-------|
| US1.1 | Fourth Sessions section on `/admin/profile`, no new route | **Pass** (code) | `profile-page.tsx`: avatar → personal info → password → `AccountSessionsCard`. Browser **Skipped**. |
| US1.2 | Only own active sessions | **Pass** | Active filter + owner scope on live list |
| US1.3 | Started, expires, last activity, network address, UA | **Pass** | Live list `createdAt` is a real timestamp after the post-QA `@be` fix |
| US1.4 | Guest `/admin/profile` → login | **Pass** (API + shell) | Sessions API 401; HTML 200 SPA |
| US1.5 | Never see user B sessions | **Pass** | Demo vs admin ids |
| US2.1 | Current highlighted | **Pass** (API `current` + Vitest) | Browser **Skipped** |
| US2.2 | No Revoke on current | **Pass** (API 400 + Vitest) | Browser **Skipped** |
| US2.3 | Identity after login without secrets in URL | **Pass** | `sessionId` in JSON; `X-Session-Id` header |
| US2.4 | Identity after renewal | **Pass** | New `sessionId` after `POST /auth/refresh`; header identifies current |
| US3.1 | Revoke other; this device stays | **Pass** | 204 + relist |
| US3.2 | Cancel confirmation | **Skipped** | UI dialog; not automated |
| US3.3 | Revoked device cannot continue | **Pass** (list) | Revoked id omitted; access JWT may still work until expiry (refresh is ended) |
| US3.4 | Stale/not owned revoke | **Pass** | Demo → admin id **404** |
| US4.1 | Revoke-all others; current stays | **Pass** | |
| US4.2 | No others → no-op, not signed out | **Pass** | |
| US4.3 | Cancel revoke-all | **Skipped** | UI dialog |
| US4.4 | Other devices cannot continue | **Pass** (list) | C/D omitted after revoke-others |
| US5.1 | Sign out ends current session | **Pass** | Logout 204 |
| US5.2 | Other device list omits signed-out session | **Pass** | |
| US5.3 | No Revoke on current row | **Pass** (API + Vitest) | |
| US5.4 | Must sign in again | **Pass** (session gone) | Access token still listed other sessions until expiry |
| US6.1 | Create records IP + UA | **Pass** | |
| US6.2 | Renew updates last activity + IP/UA | **Pass** | `lastUsedAt` + UA on successor |
| US6.3 | No fingerprint | **Pass** | No `deviceId` in JSON |
| US6.4 | Unknown IP/UA → null / i18n unknown | **Pass** (code + contract) | Live rows had IP/UA; FE maps null → “Unknown” (Vitest) |

## Success criteria

| ID | Criterion | Result | Notes |
|----|-----------|--------|-------|
| SC-001 | Sessions section in &lt; 5 s | **Partial** | API list immediate; UI timing not measured in browser |
| SC-002 | Identify this device; no Revoke on it | **Pass** (API + Vitest) | Browser **Skipped** |
| SC-003 | Revoke one other | **Pass** | |
| SC-004 | Revoke-all-others | **Pass** | |
| SC-005 | Sign out ends session | **Pass** | |
| SC-006 | Unauthenticated cannot view list | **Pass** | 401 on session APIs |
| SC-007 | Cross-account isolation | **Pass** | |
| SC-008 | New sign-in fields + last activity after renew | **Pass** | Re-checked 2026-09-04: new login `createdAt=2026-09-04T14:30:55Z` (not epoch) |

## Manual / stack checks (T040)

| Check | Result | Notes |
|-------|--------|-------|
| Stack available | Yes | `http://system.local:8080` (nginx 8080) |
| Session routes loaded | **Pass** after BE restart | First probe: `GET /sessions` **404** (stale `go run` from ~13:03; routes added ~14:22). `docker compose restart be` registered `GET/DELETE/POST /api/auth/sessions*`. |
| Seeded admin/demo login | **Pass** after restart | First probe: login **401** while `users` was empty (likely leftover from BE integration truncate on live Postgres). After restart, seed accounts worked. |
| Browser owner UI | **Skipped** | MCP `cursor-ide-browser` cataloged but tools not invokable (`MCP server does not exist`) |
| Guest profile page | **Partial** | HTML 200; API 401 |
| Cross-account isolation | **Pass** | |
| No new RBAC / sidebar | **Pass** (verify docs + contracts) | `permissions.md` N/A |

## Defects

1. **[DISPLAY] `createdAt` Unix epoch — FIXED 2026-09-04** — `issueTokenPair` now sets `CreatedAt`/`UpdatedAt`; list mapping falls back from epoch/zero to `lastUsedAt`. Do **not** remove `pkg/repo` timestamp `COALESCE` (login 500: NULL into `*time.Time`). Live re-check after `docker compose restart be`: new session `createdAt=2026-09-04T14:30:55Z`.

2. **[OPS] Stale BE `go run`** — Independent Tests against a long-lived API container miss new routes until restart. Dev Dockerfile is `CMD go run .` (no air). Restart `be` after BE code changes before QA curl.

3. **[OPS] Integration tests vs live DB** — `TruncateAuthTables` in `be/test/integration` targets compose Postgres; running `make test-be-integration` against the shared DB can empty `users`. Re-seed / BE restart restored admin and demo.

4. **[UI] Browser Independent Tests skipped** — fourth card, highlight, confirm dialogs, admin Sign out menu not clicked. Covered by FE verify + Vitest + API.

## Cross-artifact notes (lightweight analyze)

- Spec US1–US6, tasks T003/T040, plan R4–R11, and `contracts/{database,endpoints,permissions}.md` align. Permissions N/A is an explicit P1 lock, not a missing file.
- No constitution conflict (JWT self-service beside `/api/auth/profile`; packages stay independent).
- BE/FE verify docs match implemented routes and the fourth Card.
- Real gap was runtime `createdAt` (fixed after first QA pass). US7/T037 remain deferred (BA).
- T004 (BA P2 note) still open in `tasks.md`; not required for P1 QA.

## QA task status

| Task | Status |
|------|--------|
| T003 Independent Test matrix noted in `qa-checklist.md` | **Done** |
| T040 US1–US6 Independent Tests on `system.local:8080` | **Done** (API Pass; UI Skipped / Partial) |

## Sign-off

- [x] P1 acceptance scenarios pass **for API Independent Tests** (US1–US6)
- [x] Automated `make test-be` / `make test-fe` passed
- [x] BE verify recorded in `be-tasks-verify.md`
- [x] FE verify recorded in `fe-tasks-verify.md`
- [ ] Browser UI Independent Tests (fourth card, highlight, Sign out menu) — **not automated**; API + Vitest only
- [x] FR-005 started time truthful — **fixed** (live `createdAt` after BE restart)

**Verdict: PASS** — P1 security + truthful started time. Browser UI still not automated. US7 remains out of scope.
