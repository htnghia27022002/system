# Implementation Plan: Super Admin Flag and Full User Edit

**Feature**: `007-super-admin` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Contracts** (authoritative): [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)

**Input**: Feature specification from `/docs/features/007-super-admin/spec.md`

## Summary

Add a **user-level Super admin flag** (`users.is_super_admin`) that, after a valid JWT, **skips all role/permission gates** on BE middleware and FE `usePermissions` / nav / `PermissionGuard`. Seed the existing administrator account as the first Super admin. Only a Super admin can grant or revoke the flag. The last Super admin cannot be cleared, deactivated, or deleted.

Extend **Users** create/edit to a **tabbed dialog** matching Account settings field class: Profile (avatar + personal), Password (set / generate), Account (role, status, Super admin). Show the flag on the Users table. Do **not** add catalog keys, a Super admin role, or admin-of-another-user sessions (006 P2).

## Technical Context

**Language/Version**: Go 1.22 (BE); TypeScript strict on Next.js 15 App Router + React 19 (FE)

**Primary Dependencies**: Gin, pgx, squirrel, golang-migrate, PostgreSQL (BE); TanStack Query, shadcn/ui, react-i18next (FE)

**Storage**: PostgreSQL `users` (alter: `is_super_admin`)

**Testing**: `make test-be` / `make test-fe`; BE unit for middleware bypass + last Super admin; FE Vitest for Super admin allow-all

**Target Platform**: Docker stack (`make up-d`) and standalone BE/FE; modern browsers

**Project Type**: Full-stack monorepo (`be/` + `fe/` independent)

**Performance Goals**: Users page and edit dialog interactive in under 5 seconds of navigation

**Constraints**: Package independence; camelCase JSON `superAdmin`; English docs; no new env vars; HIGH-risk auth symbols stay additive

**Scale/Scope**: One boolean on users; existing Users and auth surfaces

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | No new env; FE calls existing `/api/auth/*` and `/api/admin/users*` |
| II. Role-owned artifacts | **Pass** | Architect owns tasks/plan/contracts |
| III. Spec before code | **Pass** | spec → tasks → plan (project override) |
| IV. API contract alignment | **Pass** | Additive `superAdmin`; existing Users RBAC keys unchanged; Super admin is an authenticated bypass, not a new resource |
| V. English documentation | **Pass** | Feature docs English; UI EN/VI via i18n |

**Post-design re-check**: Still pass. Bypass lives in existing `Auth` / `RequirePermission` and `usePermissions`, not a new permission key.

## Complexity Tracking

GitNexus MCP **was not available**. Treat as HIGH: `SignAccessToken`, `Auth`, `RequirePermission`, `issueTokenPair`, `buildAuthUser`, `usePermissions`.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | Overloading Administrator role as bypass would mix catalog RBAC with a god-mode flag and break “role is not Super admin” (US6) |

## Project Structure

### Documentation

```text
docs/features/007-super-admin/
├── spec.md
├── tasks.md
├── plan.md
└── contracts/{database,endpoints,permissions}.md
```

### Backend (extend-in-place)

```text
be/migrations/000009_user_is_super_admin.{up,down}.sql
be/internal/models/user/user.go
be/internal/dto/user/user.go
be/internal/dto/auth/auth.go
be/internal/common/jwt/jwt.go
be/internal/middleware/auth.go
be/internal/services/user/user_service.go
be/internal/services/auth/auth_service.go
be/internal/repository/user_repository.go
be/public/handlers/user_handler.go
be/public/routes/{admin,auth,webhook}.go
be/internal/app/container.go
be/internal/database/seeders/user_seeder.go
```

### Frontend (extend-in-place)

```text
fe/src/types/auth.ts
fe/src/store/auth-store.ts
fe/src/features/access-control/hooks/use-permissions.ts
fe/src/features/access-control/types.ts
fe/src/features/access-control/components/{users-table,user-form-dialog,user-form-fields}.tsx
fe/src/locales/{en,vi}/admin.json
```

No new feature package. No FSD. No Kratos split.

## Design

### Super admin vs role

- Column on `users`, not a role slug and not a permission key.
- JWT claim `superAdmin` for FE hydrate; Auth middleware **re-reads** `users.is_super_admin` so grant/revoke applies before token expiry.
- `RequirePermission`: if context Super admin → `c.Next()`.
- FE `usePermissions`: if `user.superAdmin` → all checks true. Nav and `PermissionGuard` already use the hook.

### Grant / last Super admin

`usersvc.Actor{ID, SuperAdmin}` on Create / Update / Delete.

- `req.SuperAdmin != nil && !actor.SuperAdmin` → forbidden
- Clearing / deactivating / deleting when `CountSuperAdmins() <= 1` and target is Super admin → forbidden
- Keep existing cannot-delete / cannot-deactivate **self**

### Users UI

Keep list + dialog. Tabs inside the dialog (not a new route). Password generate reuses `generateSecurePassword`. Super admin checkbox only when actor `isSuperAdmin`.

### Auth wiring

`Container.UserRepo` passed into `middleware.Auth(jwt, roleRepo, userRepo)` so Super admin is fresh without a second JWT secret or env.

## Risks

- HIGH: middleware bypass must not skip **authentication**.
- HIGH: last Super admin protection must cover flag clear, deactivate, and delete.
- Dev BE is `go run .` — **restart `be` after BE edits** before Independent Tests.
