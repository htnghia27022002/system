# BE tasks verify: Account Sessions

**Feature:** `docs/features/006-account-sessions/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@be`  
**Date:** 2026-09-04

## Summary

Owner-only JWT session APIs are implemented on existing auth layers. Login, register, OAuth callback, and refresh now persist IP / user-agent / `last_used_at`, return additive `sessionId`, and expose list / revoke / revoke-others under `/api/auth/sessions*`. Unit, integration, and e2e tests passed. No new RBAC catalog keys.

## GitNexus blast radius

Re-ran from repo root before editing shared symbols (index was current):

| Symbol | Direction | Risk | Impacted | Action taken |
|--------|-----------|------|----------|----------------|
| `AuthRepository` | upstream | **HIGH** | 36 | Additive methods only: `ListActiveByUserID`, `RevokeByID`, `RevokeAllExcept`. Tokens stay on AuthRepository (not UserRepository). `MockAuthRepo` updated. |
| `issueTokenPair` | upstream | **HIGH** | 4 (Login, Register, OAuth `Callback`; Refresh sibling path) | Additive: IP/UA/`last_used_at` + `sessionId`. Existing JSON keys `accessToken`, `refreshToken`, `user` unchanged. |

**Warning:** Both symbols are HIGH. Callers (auth service, OAuth callback, handlers, DI, mocks) were updated in-place. Do not treat AuthRepository edits as local.

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T001 | Confirm migration `000008`; no new env vars | Done | `be/migrations/000008_refresh_token_last_used.{up,down}.sql`; no `.env` changes |
| T005 | `last_used_at` column + backfill + NOT NULL | Done | Up SQL: add with `DEFAULT NOW()`, backfill from `created_at`, drop default, partial index |
| T006 | `LastUsedAt` on `RefreshToken`; DeviceID unused | Done | `be/internal/models/auth/auth.go` |
| T007 | AuthRepository session methods | Done | `interfaces/auth.go` + `auth_repository.go` |
| T008 | `MockAuthRepo` implements new methods | Done | `be/test/testutil/mocks.go` (`var _ interfaces.AuthRepository`) |
| T009 | Session DTOs + `sessionId` on auth payloads | Done | `be/internal/dto/auth/auth.go`; pointer IP/UA → JSON `null` |
| T010 | `issueTokenPair` metadata + `sessionId` | Done | `auth_service.go`; `pkg/repo.Insert` RETURNING fills `id` |
| T011 | Refresh rotation copies `created_at`, new `sessionId` | Done | `Refresh` in `auth_service.go` |
| T012 | Pass `httpx.ClientIP` + User-Agent from handlers | Done | `auth_handler.go` → Login/Register/Refresh/OAuth callback |
| T013 | JWT session routes, no RBAC | Done | `GET/DELETE/POST /api/auth/sessions*` in `routes/auth.go` |
| T014 | `current` via `X-Session-Id` | Done | `ListSessions` + handler header |
| T015 | DELETE current → 400; not owned/inactive → 404 | Done | `RevokeSession`; unit + integration |
| T017 | Unit list owner/active filter | Done | `be/test/unit/auth/session_service_test.go` |
| T027 | Integration current flag + refresh identity | Done | `be/test/integration/auth_sessions_integration_test.go` |
| T028 | RevokeByID owner scope | Done | Service + `TestRevokeSessionOwnerScopeAndCurrentGuard` |
| T031 | RevokeAllExcept + revoke-others HTTP | Done | Service + unit + integration |
| T034 | Logout still revokes refresh row | Done | `TestLogoutRevokesRefreshTokenSoItIsNoLongerActive` |
| T035 | Empty IP/UA → JSON null; no `device_id` | Done | `TestListSessionsNullIPAndUserAgentWhenEmpty` + login unit assertions |
| T038 | No `sessions:view` / `RequireView` | Done | `catalog.go` unchanged; session routes JWT only |

## Verification commands

```bash
make test-be
make test-be-integration
make test-be-e2e
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | Pass | `be/test/unit/auth` 1.796s |
| `make test-be-integration` | Pass | After recovering dirty schema_migrations v8 (see below) |
| `make test-be-e2e` | Pass | Existing auth e2e still green with additive `sessionId` |

First integration attempt failed: `000008` `SET NOT NULL` saw NULL `last_used_at` (live API inserts racing the add-column). golang-migrate marked version 8 dirty. Up SQL was changed to add the column with `DEFAULT NOW()`, then backfill from `created_at` and drop the default. Schema was forced back to version 7, then integration re-ran successfully.

## Acceptance coverage (BE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| US1 list only own active sessions | `TestListSessionsReturnsOnlyActiveOwnedRows` | Pass |
| US2 current via `X-Session-Id`; identity after refresh | `TestAuthSessionsListMarksCurrentAndSurvivesRefresh` | Pass |
| US3 revoke one; 400 current; 404 other/inactive | `TestRevokeSessionOwnerScopeAndCurrentGuard` + HTTP 400/204 | Pass |
| US4 revoke others; no-op 200; 400 missing id | `TestRevokeOtherSessionsKeepsCurrentAndNoOpsWhenAlone` + HTTP | Pass |
| US5 logout revokes current refresh token | `TestLogoutRevokesRefreshTokenSoItIsNoLongerActive` | Pass |
| US6 IP/UA recorded or JSON null; no device_id | login unit + `TestListSessionsNullIPAndUserAgentWhenEmpty` | Pass |
| FR-002 no new RBAC keys | T038; routes use `middleware.Auth` only | Pass |

## Gaps / follow-ups

- [x] `[BE]` tasks claimed above are done; ready for `@fe` then `@qa` Independent Tests
- [x] **Post-QA fix:** persist `CreatedAt`/`UpdatedAt` on `issueTokenPair`; list DTO falls back from epoch; do not drop `pkg/repo` timestamp COALESCE (NULL scan 500). Live list `createdAt` re-checked 2026-09-04 after `docker compose restart be`.
- [x] Existing clients that ignore unknown JSON fields keep working; FE stores `sessionId` (`auth-token-service`) and sends `X-Session-Id`

## Sign-off (BE)

- [x] All claimed `[BE]` tasks done or explicitly deferred above
- [x] Tests listed above passed
- [x] Matches `plan.md` BE sections (or deviations documented)

**Deviations:** None vs contracts. Migration up SQL uses a temporary `DEFAULT NOW()` during add/backfill then drops it so concurrent inserts cannot leave NULL — equivalent to the contract’s add → backfill → SET NOT NULL sequence.
