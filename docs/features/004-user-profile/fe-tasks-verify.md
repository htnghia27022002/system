# FE tasks verify: User Profile (004-user-profile)

**Feature:** `docs/features/004-user-profile/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@fe`  
**Date:** 2026-07-25

## Summary

P1 frontend for self-service profile, admin locale flags, and access-control personal-field alignment is implemented. Route `/admin/profile` is wired from Account settings; forms cover personal fields, avatar upload, and change password; access-control user create/edit shares the same personal Zod rules and optional admin avatar upload. `make test-fe` passed (34 tests).

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T002 | App Router profile target + feature folder | Done | `fe/src/app/admin/profile/page.tsx`, `fe/src/features/user-profile/` |
| T003 | Shared locale flags extract | Done | `fe/src/components/common/locale-flags.tsx`; `locale-select.tsx` reuses it |
| T013 | Extend `AuthUser` profile fields + `hasPassword` | Done | `fe/src/types/auth.ts` |
| T014 | Shared personal Zod schemas | Done | `fe/src/features/user-profile/schemas/personal-fields-schema.ts` |
| T015 | Account settings → `/admin/profile` | Done | `user-menu-content.tsx` Link (not disabled) |
| T016 | Thin profile route | Done | `app/admin/profile/page.tsx` (`robots: noindex`) |
| T017 | Profile page shell | Done | `profile-page.tsx` loads `/auth/me` |
| T018 | Feature barrel + `UserInfo` avatarUrl | Done | `user-profile/index.ts`, `user-info.tsx` + `resolveMediaUrl` |
| T020 | Profile form + API | Done | `profile-form.tsx`, `services/profile-api.ts` |
| T021 | Social links editor | Done | `social-links-editor.tsx` (0–5) |
| T022 | Profile i18n EN/VI | Done | `locales/en/admin.json`, `locales/vi/admin.json` (`profile.*`) |
| T024 | Avatar upload | Done | `avatar-upload.tsx` client checks + API |
| T025 | Initials fallback | Done | Profile + `UserInfo` via `useInitials` |
| T027 | Change password form | Done | `change-password-form.tsx` |
| T028 | Hide when `hasPassword === false` | Done | Explanation copy when unavailable |
| T030 | Access-control types/schemas/API | Done | Extended personal fields + `uploadUserAvatar` |
| T031 | User form personal section | Done | `user-form-fields.tsx` / dialog |
| T032 | Admin avatar upload on edit | Done | `POST /admin/users/:id/avatar` |
| T033 | Access-control i18n | Done | EN/VI admin.json personal fields |
| T034 | Extract locale flags module | Done | `locale-flags.tsx` |
| T035 | Flag locale in user menu | Done | EN/VI flag items + active check |
| T036 | Replace switchToVi/En copy | Done | `userMenu.language` / locale labels |
| T040 | Vitest schema coverage | Done | `personal-fields-schema.test.ts` + access-control schema tests |

**Deferred (P2 US7):** T037, T038 — polish / a11y pass left open per plan.

## Verification commands

```bash
make test-fe
pnpm lint
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass | 10 files / 34 tests |
| `pnpm lint` | Fail (pre-existing) | Errors in landing hero / unused imports in sidebar & admin-layout — not introduced by this feature. New profile/access-control files only show existing react-hook-form `watch` compiler warning. |

## Acceptance coverage (FE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| US1 Account settings → profile | Menu Link + route + page shell | Pass (code) |
| US2 Edit personal fields + validation | Profile form + Zod | Pass (unit + code) |
| US3 Avatar upload / reject / initials | AvatarUpload + validateAvatarFile | Pass (unit + code) |
| US4 Change password / OAuth-only | ChangePasswordForm + hasPassword | Pass (code) |
| US5 Access-control personal + avatar | User form + API | Pass (code) |
| US6 EN/VI flags in admin menu | UserMenuContent + locale-flags | Pass (code) |
| Guest blocked on `/admin/profile` | Existing admin layout guards | Pass (existing) |

## Gaps / follow-ups

- [x] P1 `[FE]` tasks T002–T036, T040 done
- [ ] T037 / T038 (US7 P2 polish) deferred
- [ ] Manual smoke on `http://system.local:8080` (avatar multipart, password change, live i18n) — for `@qa`
- [ ] `pnpm lint` clean-up of unrelated pre-existing errors outside this feature

## Sign-off (FE)

- [x] All claimed P1 `[FE]` tasks done; US7 explicitly deferred
- [x] `make test-fe` passed
- [x] Matches `plan.md` FE sections (feature module, route, chrome, access-control alignment, flag locale)
- [x] Ready for `@qa` with caveat: manual Independent Tests still required; lint noise is pre-existing
