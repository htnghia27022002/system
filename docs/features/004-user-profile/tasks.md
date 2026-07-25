# Tasks: User Profile, Access-Control Alignment, Admin Locale Flags (004-user-profile)

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks; expanded personal fields included)

**Tests**: Spec does not require TDD. Prefer unit tests for BE validation/password change and FE schema validation where practical. QA owns manual acceptance against Independent Tests.

**Organization**: Phases follow user stories (US1–US6 in P1; US7 = P2 polish). Task labels use `[BE]` / `[FE]` / `[QA]` per constitution; `[USn]` maps to spec user stories.

**GitNexus blast radius (design note)**: Extending `ToResponse` is **HIGH** (admin user List/Create/Update). Changing `UserInfo` is **CRITICAL** (admin chrome). Changing `UserMenuContent` is **HIGH**. Prefer additive fields and keep chrome updates backward-compatible.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[BE]` / `[FE]` / `[QA]`**: Role ownership
- **`[P]`**: Parallelizable (different files, no incomplete dependency)
- **`[USn]`**: User story from spec (story phases only)
- Include exact file paths in descriptions

## Path Conventions

- Backend: `be/` (Go module)
- Frontend: `fe/src/` (App Router + `features/`)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm stack targets, feature folder ownership, and env ownership for uploads.

- [x] T001 [BE] Confirm next migration version after `be/migrations/000004_*` and document upload env ownership in `be/.env.example` (`UPLOAD_DIR`, `PUBLIC_BASE_URL` or equivalent) — never root `.env`
- [x] T002 [FE] [P] Confirm App Router target `fe/src/app/admin/profile/page.tsx` under existing `fe/src/app/admin/layout.tsx` (ProtectedGuard + AdminGuard) and feature home `fe/src/features/user-profile/`
- [x] T003 [FE] [P] Confirm reuse of `fe/src/components/common/locale-select.tsx` flag assets / extract shared flag control for admin menu (US6)

**Checkpoint**: Env ownership and route/feature targets agreed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, DTOs, shared validation, media storage, and profile API surface. **BLOCKS all user stories.**

**⚠️ CRITICAL**: No story UI work depends on unfinished migration + profile contract.

- [x] T004 [BE] Add migration `be/migrations/000005_user_profile_fields.up.sql` (+ `.down.sql`) adding nullable `general` (TEXT, max enforced in app), `birthday` (DATE), `address` (VARCHAR(500)), `social_links` (JSONB NOT NULL DEFAULT `[]`); keep existing `phone` / `avatar_url`
- [x] T005 [BE] Extend `be/internal/models/user/user.go` with `General`, `Birthday`, `Address`, `SocialLinks` (JSON slice) matching camelCase JSON tags
- [x] T006 [BE] [P] Add personal-field DTOs + validation helpers in `be/internal/dto/user/` (and/or `be/internal/dto/auth/`) for social links (0–5, http(s) URL, optional label ≤50), phone ≤50, general ≤1000, address ≤500, birthday not future
- [x] T007 [BE] Implement local avatar storage service under `be/internal/services/media/` (or `be/pkg/` if storage-agnostic) writing under `UPLOAD_DIR`; enforce JPEG/PNG/WebP and **2 MB** max; return public path stored in `avatar_url`
- [x] T008 [BE] Wire static/media serving for uploaded avatars (Gin route under `/api/media/...` or equivalent) so nginx `/api/` proxy serves files without a separate FE host; document Docker volume for `UPLOAD_DIR` in `be/.env.example` / compose notes if needed
- [x] T009 [BE] Extend `AuthUserResponse` in `be/internal/dto/auth/auth.go` + `buildAuthUser` in `be/internal/services/auth/auth_service.go` with profile display fields (`phone`, `avatarUrl`, `general`, `birthday`, `address`, `socialLinks`) and `hasPassword` (true when `password_hash` non-empty)
- [x] T010 [BE] Add self-service routes in `be/public/routes/auth.go` + handlers: `GET` profile (or rely on extended `/me`), `PATCH /api/auth/profile`, `POST /api/auth/profile/avatar` (multipart), `POST /api/auth/change-password`; wire DI in `be/internal/app/dependency/`
- [x] T011 [BE] Extend admin user DTOs/responses in `be/internal/dto/user/user.go` + `Create`/`Update`/`ToResponse` in `be/internal/services/user/user_service.go` with the same personal fields (phone, general, birthday, address, socialLinks, avatarUrl); blank password on update still keeps hash
- [x] T012 [BE] [P] Add optional admin avatar upload `POST /api/admin/users/:id/avatar` in `be/public/routes/` + `be/public/handlers/user_handler.go` reusing media service
- [x] T013 [FE] [P] Extend `fe/src/types/auth.ts` `AuthUser` (+ any API mappers) with profile fields and `hasPassword` / `avatarUrl` for chrome + profile
- [x] T014 [FE] [P] Add shared personal-field Zod schemas (reusable by profile + access-control) under `fe/src/features/user-profile/schemas/` (or shared module exported for access-control)

**Checkpoint**: Migration applies; profile + change-password + avatar APIs callable; `/me` returns new fields; admin user JSON includes personal fields.

---

## Phase 3: User Story 1 — Open profile from Account settings (Priority: P1) 🎯 MVP entry

**Goal**: Account settings navigates to `/admin/profile`; guests redirected; page shows current personal info.

**Independent Test**: Sign in → sidebar user menu → Account settings → `/admin/profile` shows name, email, avatar/initials, personal fields.

### Implementation

- [x] T015 [FE] [US1] Enable Account settings in `fe/src/components/common/user-menu-content.tsx` to navigate to `/admin/profile` (Link / router); remove `disabled`
- [x] T016 [FE] [US1] Add thin route `fe/src/app/admin/profile/page.tsx` (`robots: noindex`) rendering profile page from `@/features/user-profile`
- [x] T017 [FE] [US1] Implement profile page shell `fe/src/features/user-profile/components/profile-page.tsx` loading current user via existing auth/`/me` (or profile GET) and showing name, read-only email, avatar/initials, phone, General, birthday, address, social links (empty-safe)
- [x] T018 [FE] [P] [US1] Export public API from `fe/src/features/user-profile/index.ts`; wire `UserInfo` in `fe/src/components/common/user-info.tsx` to show `avatarUrl` when present (initials fallback) without breaking admin chrome

**Checkpoint**: US1 Independent Test passes (FR-001, FR-009, FR-010).

---

## Phase 4: User Story 2 — Edit personal information (Priority: P1)

**Goal**: Self-service edit of name + optional personal fields; email read-only; validation per FR-019–FR-024.

**Independent Test**: Edit valid fields → save → reload persists; invalid future birthday / bad URL / over-long General blocked.

### Implementation

- [x] T019 [BE] [US2] Implement `PATCH /api/auth/profile` service logic (self only; never accept email change) in `be/internal/services/auth/` or dedicated profile service; trim whitespace-only optionals to empty; enforce validation
- [x] T020 [FE] [US2] Build profile edit form `fe/src/features/user-profile/components/profile-form.tsx` + API client methods under `fe/src/features/user-profile/services/` calling `PATCH /auth/profile`
- [x] T021 [FE] [US2] Add social-links editor UI (add/remove up to 5; label optional; URL required) and field-level error display; success toast/message on save
- [x] T022 [FE] [P] [US2] Add i18n keys for profile fields/errors under `fe/src/locales/en/` and `fe/src/locales/vi/` (label **General** in English)

**Checkpoint**: US2 Independent Test passes (SC-002, SC-008).

---

## Phase 5: User Story 3 — Upload or change avatar (Priority: P1)

**Goal**: Upload/replace avatar within type/size limits; initials when absent; failed upload keeps prior avatar.

**Independent Test**: Upload valid image; reject oversize/wrong type with message; reload shows new avatar.

### Implementation

- [x] T023 [BE] [US3] Finalize `POST /api/auth/profile/avatar` handler: multipart field, validate type/size, replace file, update `avatar_url`, return updated user/profile
- [x] T024 [FE] [US3] Add avatar upload control on profile page (`fe/src/features/user-profile/components/avatar-upload.tsx`) with client-side pre-checks + server error surfacing; refresh auth store / displayed avatar after success
- [x] T025 [FE] [US3] Ensure initials fallback via existing `useInitials` + `Avatar` when `avatarUrl` empty across profile and `UserInfo`

**Checkpoint**: US3 Independent Test passes (SC-003).

---

## Phase 6: User Story 4 — Change own password (Priority: P1)

**Goal**: Self-service change password with current + new + confirm; hide/disable when `hasPassword` is false.

**Independent Test**: Correct current password → can login with new; wrong current / mismatch / short password rejected; OAuth-only sees disabled/hidden path with explanation.

### Implementation

- [x] T026 [BE] [US4] Implement `POST /api/auth/change-password` in auth service: verify current hash, enforce min 8 chars, reject empty password accounts with clear error
- [x] T027 [FE] [US4] Add change-password section `fe/src/features/user-profile/components/change-password-form.tsx` (current, new, confirm); wire API; show success feedback without forced re-login in P1
- [x] T028 [FE] [US4] When `hasPassword === false`, hide or disable change-password and show brief explanation (FR-008)

**Checkpoint**: US4 Independent Test passes (SC-004).

---

## Phase 7: User Story 5 — Admin edits user personal info in access control (Priority: P1)

**Goal**: Access-control create/edit supports same personal field class + optional avatar; password create required / edit blank = keep.

**Independent Test**: Edit another user’s personal fields + optional password/avatar; validations match profile; blank password keeps old.

### Implementation

- [x] T029 [BE] [US5] Ensure admin create/update accept and persist personal fields in `be/internal/services/user/user_service.go` + handlers; list/get responses include them via `ToResponse`
- [x] T030 [FE] [US5] Extend `fe/src/features/access-control/types.ts`, `schemas/access-control-schemas.ts`, and `services/access-control-api.ts` with personal fields (+ optional avatar URL)
- [x] T031 [FE] [US5] Extend `fe/src/features/access-control/components/user-form-fields.tsx` (+ dialog) with personal section: phone, General, birthday, address, social links; keep email/role/status/password semantics
- [x] T032 [FE] [US5] Add optional avatar upload/change on access-control user form using admin avatar endpoint; initials/image display consistent with profile
- [x] T033 [FE] [P] [US5] Update access-control i18n under `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json` for new fields

**Checkpoint**: US5 Independent Test passes (SC-005).

---

## Phase 8: User Story 6 — Locale via flags in admin user menu (Priority: P1)

**Goal**: Replace “Switch to …” text toggle with clickable EN/VI flags; immediate apply; active state clear; consistent with public LocaleSelect behavior.

**Independent Test**: Open admin user menu; click non-active flag; UI language switches; active indicator clear; switch back.

### Implementation

- [x] T034 [FE] [US6] Extract shared flag icons / locale options from `fe/src/components/common/locale-select.tsx` into a small shared module (e.g. `fe/src/components/common/locale-flags.tsx`) if needed for reuse
- [x] T035 [FE] [US6] Replace language `DropdownMenuItem` toggle in `fe/src/components/common/user-menu-content.tsx` with clickable EN and VI flag options (`i18n.changeLanguage`); show active locale visually; remove sole “Switch to …” pattern
- [x] T036 [FE] [P] [US6] Clean unused `userMenu.switchToVi` / `switchToEn` copy or replace with accessible labels in `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json`

**Checkpoint**: US6 Independent Test passes (SC-006). US6 can ship in parallel with BE-heavy stories once chrome-only.

---

## Phase 9: User Story 7 — Profile UX polish (Priority: P2) *(deferred)*

**Goal**: Improve messaging, upload feedback, a11y without changing routes or field set.

**Independent Test**: Walk profile + access-control personal sections; errors/labels accessible.

### Implementation

- [ ] T037 [FE] [US7] Polish success/error/empty states and optional upload progress on profile + access-control personal forms
- [ ] T038 [FE] [US7] Accessibility pass: labels, focus order, live error text for profile, avatar, change-password, and locale flags

**Checkpoint**: P2 polish complete; P1 field set unchanged.

---

## Phase 10: Polish & Cross-Cutting / QA

**Purpose**: Verify docs alignment and manual acceptance.

- [x] T039 [BE] [P] Add/adjust unit tests under `be/test/unit/` for profile validation, change-password, and avatar rejection cases
- [x] T040 [FE] [P] Add/adjust Vitest coverage for profile/access-control Zod schemas where practical
- [ ] T041 [QA] Manual: US1–US6 Independent Tests + SC-001–SC-008 on `http://system.local:8080` (or configured origin)
- [ ] T042 [QA] Manual: OAuth-only / no-password account sees safe change-password UX; guest cannot open `/admin/profile`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → story phases
- **US1** needs foundational `/me`/profile read + FE route/menu
- **US2–US4** need foundational APIs; FE can stub UI after contract exists
- **US5** needs BE personal fields on admin user APIs + FE form extension
- **US6** is FE-only on admin chrome; can parallel after Setup (T003) once LocaleSelect flags exist
- **US7** after P1 stories desired complete
- **QA** after BE/FE verify for implemented stories

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|--------|
| US1 | Phase 2 (read profile fields) | MVP discoverability |
| US2 | Phase 2 + US1 shell | Edit persistence |
| US3 | Phase 2 media + US1 shell | Avatar |
| US4 | Phase 2 `hasPassword` + change-password API | Can parallel US2/US3 |
| US5 | Phase 2 admin DTO extensions | Can parallel profile FE after BE |
| US6 | Setup (flags) | Parallel FE-only |
| US7 | P1 stories | Deferred polish |

### Parallel Opportunities

```text
After Phase 2:
  @be: US2/US3/US4 API polish in parallel with @fe US1 shell + US6 flags
  @fe US2/US3/US4 after profile route exists
  @fe US5 after admin DTO contract ready (T011/T012)
```

### Suggested MVP

1. Phase 1–2 foundational  
2. US1 (Account settings → profile view)  
3. US6 (locale flags — quick chrome win)  
4. US2 → US3 → US4  
5. US5 access-control alignment  
6. US7 polish (optional for P1 done)

---

## Implementation Strategy

### Recommended start order for agents

1. **`@be` first**: T001, T004–T012 (schema → media → auth profile APIs → admin DTO)  
2. **`@fe` in parallel once `/me` + profile contract stable**: T013–T018 (US1), T034–T036 (US6)  
3. **`@fe`**: US2 → US3 → US4 → US5  
4. **`@qa`**: after `be-tasks-verify.md` / `fe-tasks-verify.md`

### Notes

- `[P]` = different files, no incomplete dependency  
- Do not put upload secrets/paths in root `.env`  
- JSON camelCase everywhere (`avatarUrl`, `socialLinks`, `hasPassword`, `fullName` mapped as existing `name` where product already uses `name` in auth DTO)
- Admin password set/reset remains blank-to-keep; self-service always requires current password
