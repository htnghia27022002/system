# FE tasks verify: Account Sessions

**Feature:** `docs/features/006-account-sessions/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@fe`  
**Date:** 2026-09-04

## Summary

Account sessions are the fourth Card on existing Account settings (`/admin/profile`). The client stores `sessionId` in localStorage, sends `X-Session-Id` on list/revoke, highlights this device, and revokes other sessions with confirm dialogs. Admin Sign out now calls `useSignOut` so the current refresh token is revoked on the server. Vitest passed (46 tests). No new PermissionKeys or sidebar item.

## GitNexus blast radius

Re-ran from repo root before editing `UserMenuContent`:

| Symbol | Direction | Risk | Impacted | Action taken |
|--------|-----------|------|----------|----------------|
| `UserMenuContent` | upstream | **HIGH** | 3 (direct: `NavUser`; processes: Admin layout, AppSidebar) | Menu structure unchanged. Sign out swapped from local `auth-store.signOut()` to `useSignOut` from `@/features/auth` (JSON `{ refreshToken }` logout). |

**Warning:** HIGH chrome component. Only the Sign out action was changed.

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T002 | Fourth Card on existing profile page; no dedicated route/sidebar | Done | `profile-page.tsx` mounts `AccountSessionsCard` after avatar, personal info, password; `fe/src/app/admin/profile/page.tsx` unchanged |
| T016 | Persist `sessionId` in token service + store + refresh | Done | `auth-token-service.ts` (`session_id` localStorage); `AuthTokens.sessionId`; `auth-store.ts` `signIn`; `api-client.ts` refresh |
| T018 | Vitest camelCase list mapping | Done | `sessions-api.test.ts` |
| T019 | Types + `sessions-api` with `X-Session-Id` | Done | `types.ts` `AccountSession`; `services/sessions-api.ts` |
| T020 | TanStack Query hook | Done | `hooks/use-account-sessions.ts` (no `useEffect` for server data) |
| T021 | Sessions card: loading / empty / error+retry | Done | `account-sessions-card.tsx` |
| T022 | Mount fourth Card; barrel export | Done | `profile-page.tsx`; `user-profile/index.ts` |
| T023 | i18n EN + VI (`profile.sections.sessions`, `profile.sessions.*`) | Done | `fe/src/locales/en/admin.json`, `fe/src/locales/vi/admin.json` (UTF-8 diacritics, e.g. “Thiết bị này”) |
| T024 | Current row highlighted; no Revoke on current | Done | `account-sessions-card.test.tsx` |
| T025 | `sessionId` from login/register/OAuth/refresh; mock returns it | Done | `auth-store.ts`; `api-client.ts`; `createMockAuthTokens` → `auth.mock.ts` |
| T026 | Highlight via `current` and/or stored `sessionId`; never in URL | Done | `markCurrentSessions` in hook; header only, no query |
| T029 | Revoke non-current with shadcn AlertDialog; invalidate on success | Done | `account-sessions-card.tsx` + `useAccountSessions.revokeSession` |
| T030 | 404/stale: toast + refetch | Done | `onError` invalidates query; `profile.sessions.errors.revokeNotFound` |
| T032 | Revoke all other sessions with confirm; hidden when none | Done | Card header control; test hides button when `hasOtherSessions` is false |
| T033 | Admin Sign out uses `useSignOut` | Done | `user-menu-content.tsx` (swap-only) |
| T036 | Unknown IP/UA i18n; long UA wrap + title; no fingerprint | Done | `profile.sessions.unknown`; `line-clamp-3 break-all` + `title` |
| T039 | No PermissionKeys / sidebar / PermissionGuard | Done | `permission-keys.ts` and `app-sidebar.tsx` unchanged |

## Verification commands

```bash
make test-fe
pnpm --filter fe lint
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass | **13 files, 46 tests** (includes 6 new: 4 mapping/API + 2 card UX) |
| `pnpm --filter fe lint` | Fail (pre-existing) | Repo-wide ESLint already fails in landing hero + unused imports. **Changed session files: clean** (`eslint` on the files listed below exited 0) |

Changed files linted clean:

`account-sessions-card.tsx`, `account-sessions-card.test.tsx`, `use-account-sessions.ts`, `sessions-api.ts`, `sessions-api.test.ts`, `profile-page.tsx`, `user-menu-content.tsx`, `auth-token-service.ts`, `api-client.ts`, `auth-store.ts`, `types/auth.ts`, `mock-jwt.ts`, `user-profile/index.ts`, `user-profile/types.ts`.

Browser E2E against `http://system.local:8080/admin/profile` was **not** run in this agent session (no browser tools). `@qa` owns Independent Tests T040.

## Acceptance coverage (FE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| US1 fourth Sessions section on `/admin/profile` | T002/T022; no new App Router page | Pass (code) |
| US1 list fields camelCase | `sessions-api.test.ts` | Pass |
| US1 loading / empty / error+retry | `account-sessions-card.tsx` | Pass (code; not separately unit-tested) |
| US2 current highlight; no Revoke on current | `account-sessions-card.test.tsx` | Pass |
| US2 `sessionId` after login/refresh, not in URL | token service + interceptor + header tests | Pass |
| US3 revoke other + confirm; 404 refetch | AlertDialog + hook `onError` invalidate | Pass (code) |
| US4 revoke-all-others hidden when alone | card test | Pass |
| US5 Sign out calls server logout | `useSignOut` in user menu | Pass (code) |
| US6 unknown IP/UA; long UA title | card + mapping nulls test | Pass |
| FR-002 no new RBAC / sidebar | T039 | Pass |
| Guest redirect / two-browser / Sign out list truth | T040 Independent Tests | **N/A — `@qa`** |

## Gaps / follow-ups

- [ ] Ready for `@qa` Independent Tests (T003/T040) on `http://system.local:8080/admin/profile`: two/three browsers, revoke one, revoke-all, Sign out from admin user menu, guest redirect, cross-account isolation, refresh identity (`sessionId` updates in storage not URL)
- [ ] Repo-wide `pnpm lint` still fails on **unrelated** landing/sidebar files (pre-existing)
- [ ] Loading / empty / error UI paths are implemented but not covered by dedicated Vitest cases (happy-path + current-row tests only)

## Sign-off (FE)

- [x] All claimed `[FE]` tasks done or explicitly deferred above
- [x] Tests listed above passed; lint skipped at repo level with reason (pre-existing); changed files clean
- [x] Matches `plan.md` FE sections (fourth Card, `user-profile/` kebab-case files, JWT self-service, no PermissionKeys)
