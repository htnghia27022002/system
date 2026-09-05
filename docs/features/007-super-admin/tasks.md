# Tasks: Super Admin Flag and Full User Edit

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Contracts**: [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks). Generated from **spec.md + locked technical design** (project override: tasks before plan).

**Tests**: Spec Independent Tests plus: BE unit tests for Super admin bypass, last-Super-admin protection, and flag-only grant. FE tests for `usePermissions` Super admin allow-all. QA owns live Independent Tests after `be` restart.

**Organization**: Phases follow P1 user stories (US1–US6). Admin-of-another-user sessions is P2 / 006 — do not implement. Task labels use `[BE]` / `[FE]` / `[QA]` / `[BA]`.

**GitNexus blast radius (design note — HIGH)**: Editing `SignAccessToken`, `Auth` middleware, `RequirePermission`, `issueTokenPair` / `buildAuthUser`, `usePermissions`, and user Create/Update/Delete is **HIGH**. Prefer **additive** `superAdmin` / `IsSuperAdmin` fields. Do not remove existing JSON keys or delete RBAC. GitNexus MCP was **not available** in the implementing session — treat these symbols as HIGH and keep changes additive.

**P1 only**: No admin sessions-for-other-user, no Super admin role, no new catalog keys.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[BE]` / `[FE]` / `[QA]` / `[BA]`**: Role ownership
- **`[P]`**: Parallelizable
- **`[USn]`**: User story from spec

## Path Conventions

- Backend: `be/` — `public/routes` → `public/handlers` → `internal/services` → `internal/repository`
- Frontend: `fe/src/features/access-control/` + auth store/types; kebab-case

---

## Phase 1: Setup

- [x] T001 [BE] Confirm next migration after `be/migrations/000008_refresh_token_last_used` is `be/migrations/000009_user_is_super_admin.{up,down}.sql`; no new env vars
- [x] T002 [FE] [P] Confirm Users stays at `fe/src/app/admin/users/page.tsx` with dialog edit — **do not** add `/admin/users/:id` or a Super admin sidebar item
- [x] T003 [QA] [P] Note Independent Tests from spec US1–US6 for later `qa-checklist.md`

**Checkpoint**: Migration number and no-new-route lock agreed.

---

## Phase 2: Foundational (Blocking)

- [x] T004 [BE] Add `be/migrations/000009_user_is_super_admin.up.sql` (+ `.down.sql`): `users.is_super_admin BOOLEAN NOT NULL DEFAULT FALSE`; backfill seeded admin; partial index optional
- [x] T005 [BE] Add `IsSuperAdmin` on `be/internal/models/user/user.go` (`db:"is_super_admin"`)
- [x] T006 [BE] Add `SuperAdmin` on user DTOs (`be/internal/dto/user/user.go`) and `AuthUserResponse` (`be/internal/dto/auth/auth.go`); map in `ToResponse` / `buildAuthUser`
- [x] T007 [BE] Add `CountSuperAdmins` on `be/internal/repository/interfaces/user.go`, `be/internal/repository/user_repository.go`, and `be/test/testutil/mocks.go` `MemoryUserRepo`
- [x] T008 [BE] Seed first Super admin: migration UPDATE + `be/internal/database/seeders/user_seeder.go` insert column for new installs
- [x] T009 [BE] Add `superAdmin` claim on `be/internal/common/jwt/jwt.go` `Claims` / `SignAccessToken`; update `issueTokenPair` and refresh in `be/internal/services/auth/auth_service.go`; update `be/test/unit/jwt/jwt_test.go`
- [x] T010 [BE] Expose `UserRepo` on `be/internal/app/container.go`; pass it into `middleware.Auth` from `be/public/routes/admin.go`, `auth.go`, `webhook.go`
- [x] T011 [BE] `be/internal/middleware/auth.go`: set Super admin from JWT then refresh from `userRepo.GetByID`; `RequirePermission` allows Super admin; add `IsSuperAdmin` / `GetUserID` helpers as needed

**Checkpoint**: Flag persisted, JWT + middleware can bypass, seed admin is Super admin.

---

## Phase 3: User story 1–2 — Bypass + table flag

- [x] T012 [BE] [US1] Middleware Super admin bypass covered by unit test (`be/test/unit/middleware/` or adjacent)
- [x] T013 [BE] [US2] User list/get responses include `superAdmin` via existing handlers (no new routes)
- [x] T014 [FE] [US1] Add `superAdmin` on `fe/src/types/auth.ts` (`AuthUser`, `JwtPayload`); hydrate in `fe/src/store/auth-store.ts`; `isAdmin` true if role admin **or** Super admin
- [x] T015 [FE] [US1] `fe/src/features/access-control/hooks/use-permissions.ts`: Super admin → all `hasPermission` / `canView` / `canModify` / `hasAny` true; expose `isSuperAdmin`; test in `use-permissions.test.tsx`
- [x] T016 [FE] [US2] `ManagedUser` + create/update inputs in `fe/src/features/access-control/types.ts`; Users table + mobile card Super admin column in `users-table.tsx`
- [x] T017 [FE] [P] [US1] Mock path: `AuthResolvedUser.superAdmin`, seed admin flag, `createMockAuthTokens` claim, `auth.mock.ts` user object (`fe/src/services/mock/`, `fe/src/utils/mock-jwt.ts`)

**Checkpoint**: Super admin with empty permissions can use admin UI; table shows the flag.

---

## Phase 4: User story 3–4 — Full edit + last Super admin

- [x] T018 [BE] [US3] [US4] `usersvc.Actor` + Create/Update/Delete rules in `be/internal/services/user/user_service.go`; handlers pass `middleware.GetUserID` + `IsSuperAdmin`; only Super admin may set the flag; protect last Super admin; keep self delete/deactivate blocks
- [x] T019 [BE] [US4] Unit tests in `be/test/unit/user/user_service_test.go` for grant forbidden, last-flag, last-delete
- [x] T020 [FE] [US3] Schemas add `superAdmin` (`access-control-schemas.ts`)
- [x] T021 [FE] [US3] Tabbed user dialog (Profile / Password / Account) in `user-form-dialog.tsx` + `user-form-fields.tsx`; generate password via `generateSecurePassword`; Super admin checkbox only if `isSuperAdmin`
- [x] T022 [FE] [P] [US3] EN + VI copy in `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json` (UTF-8 Vietnamese)
- [x] T023 [FE] [P] [US3] Export `generateSecurePassword` from `fe/src/features/user-profile/index.ts` if imported across features
- [x] T024 [FE] [US4] Mock `createUser` / `updateUser` persist `superAdmin` and reject last-Super-admin removal when practical

**Checkpoint**: Users edit matches Account settings field class; last Super admin is protected.

---

## Phase 5: User story 5–6 — Non-bypass unchanged

- [x] T025 [BE] [US5] Existing RequireView/RequireModify still forbids non–Super admin without keys (covered by existing tests + middleware test)
- [x] T026 [FE] [US6] Do **not** treat Administrator role as Super admin in `usePermissions` or JWT mapping
- [x] T027 [QA] Independent Tests after **restart `be`** (dev is `go run .`)

---

## Phase 6: Polish

- [x] T028 [BE] `make test-be` + `docs/features/007-super-admin/be-tasks-verify.md`
- [x] T029 [FE] `make test-fe` + `docs/features/007-super-admin/fe-tasks-verify.md`
- [x] T030 [QA] `qa-checklist.md` + `make test` + Independent Tests
- [x] T031 [BA] Point `.specify/feature.json` and `.cursor/rules/specify-rules.mdc` current plan at `docs/features/007-super-admin`

---

## Dependencies

- T004–T011 block T012–T021
- T014–T015 block nav/page bypass (US1)
- T018 blocks T019
- T028/T029 before T030

## Parallel opportunities

- T002/T003 with T001
- T016/T017 after T014
- T022/T023 with T021
