# Implementation Plan: User Profile, Access-Control Alignment, Admin Locale Flags

**Feature**: `004-user-profile` | **Date**: 2026-07-25 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Input**: Feature specification from `/docs/features/004-user-profile/spec.md`

## Summary

Deliver authenticated self-service profile at `/admin/profile` (personal fields, avatar upload, change password), wire **Account settings** in the admin sidebar user menu, align access-control user create/edit with the same personal field class (including optional avatar), and replace the admin menu language “Switch to …” toggle with clickable **EN/VI flags** consistent with public `LocaleSelect` behavior.

**Backend** extends the existing `users` table (`phone`, `avatar_url` already present) with `general`, `birthday`, `address`, and `social_links` (JSONB). Self-service APIs live under `/api/auth/*`; admin APIs extend `/api/admin/users`. Avatars use **local disk storage** owned by BE (`UPLOAD_DIR` in `be/.env`), served under `/api/media/...` via the existing nginx `/api/` proxy.

**Frontend** adds `fe/src/features/user-profile/`, route `fe/src/app/admin/profile/page.tsx`, extends access-control forms, and updates admin chrome (`UserMenuContent`, `UserInfo`).

## Technical Context

**Language/Version**: Go 1.22 (BE); TypeScript strict on Next.js 15 App Router + React 19 (FE)

**Primary Dependencies**: Gin, GORM, golang-migrate, PostgreSQL (BE); TanStack Query, Zustand, react-hook-form, Zod, shadcn/ui, react-i18next (FE)

**Storage**: PostgreSQL for user profile fields; local filesystem for avatar binaries (`UPLOAD_DIR`); `avatar_url` stores public API-relative or absolute URL string

**Testing**: `make test-be` / `make test-fe`; BE unit tests for validation/password/avatar rejection; FE Vitest for Zod schemas; QA manual Independent Tests

**Target Platform**: Docker stack (`make up-d`) and standalone BE/FE deploys; modern browsers; admin-capable authenticated users

**Project Type**: Full-stack monorepo feature (`be/` + `fe/` independent packages)

**Performance Goals**: Profile load/save feels interactive on stable LAN; avatar ≤2 MB; no mandatory CDN in P1

**Constraints**: Package independence (HTTP only); camelCase JSON; English docs; no self-service email change; P2 polish deferred; GitNexus HIGH/CRITICAL on `ToResponse`, `UserInfo`, `UserMenuContent` — additive changes only

**Scale/Scope**: One profile route; extend existing admin users CRUD; max 5 social links per user; EN/VI locales only

## Constitution Check

*GATE: Must pass before design lock. Re-checked after design below.*

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | Upload config in `be/.env`; FE uses `NEXT_PUBLIC_API_BASE_URL` only; no cross-imports |
| II. Role-owned artifacts | **Pass** | Architect owns `tasks.md` / `plan.md`; no app code in this phase |
| III. Spec before code | **Pass** | `spec.md` → `tasks.md` → `plan.md` before `@be` / `@fe` |
| IV. API contract alignment | **Pass** | Contracts below use camelCase under `/api` |
| V. English documentation | **Pass** | Feature docs English; UI EN/VI via i18n |

**Post-design re-check**: Still pass. Media serving stays inside BE `/api` so FE and nginx need no app-secret coupling. Shared personal Zod schemas may live in `user-profile` and be imported by access-control via the feature public barrel (or a thin shared schema module under `fe/src/`) without importing `be/`.

## Complexity Tracking

> No constitution violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

---

## Research Decisions (Phase 0)

### R1 — Avatar storage: local disk + URL field (not S3, not base64)

- **Decision**: Store files under configurable `UPLOAD_DIR` (default e.g. `./data/uploads` in container/host). Persist public path in existing `users.avatar_url`. Serve via Gin static/media handler at **`GET /api/media/avatars/:filename`** (auth optional for GET if filenames are unguessable UUIDs; prefer require auth or use opaque names). Upload via multipart **`POST /api/auth/profile/avatar`** and **`POST /api/admin/users/:id/avatar`**.
- **Rationale**: No object-storage dependency in current stack; `avatar_url` already exists; nginx already proxies `/api/`; BE remains independently deployable with a volume/path env var.
- **Limits**: JPEG, PNG, WebP; max **2 MB**; reject with clear 400; on failure keep previous `avatar_url` and file.
- **Alternatives considered**:
  - **S3/MinIO** — better for multi-instance prod later; deferred (adds infra).
  - **Base64 in DB** — rejected (payload bloat, backups).
  - **FE-only / external URL paste** — does not satisfy upload FR.
  - **Nginx alias outside `/api`** — couples root compose/nginx to app media; weaker standalone BE story.

### R2 — New columns vs sidecar table

- **Decision**: Add columns on `users`: `general TEXT`, `birthday DATE NULL`, `address VARCHAR(500)`, `social_links JSONB NOT NULL DEFAULT '[]'`. Reuse `phone`, `avatar_url`.
- **Rationale**: Spec is 1:1 with user; JSONB list matches 0–5 social links without join table for P1.
- **Alternatives**: Separate `user_profiles` table — unnecessary until profile outgrows user row.

### R3 — Self-service API placement under `/api/auth`

- **Decision**:
  - Extend **`GET /api/auth/me`** (and login/register user payload) with profile display fields + `hasPassword`.
  - **`PATCH /api/auth/profile`** — update self personal fields (not email/role/status).
  - **`POST /api/auth/profile/avatar`** — multipart upload for self.
  - **`POST /api/auth/change-password`** — `{ currentPassword, newPassword }` (confirm validated on FE; BE may also accept `confirmPassword` or rely on single newPassword after FE match).
- **Rationale**: Profile is identity/self-service, not admin RBAC; reuses JWT middleware on auth protected group.
- **Alternatives**: `/api/admin/profile` — rejected (implies admin permission beyond “signed-in admin shell user”).

### R4 — Auth DTO `name` vs DB `fullName`

- **Decision**: Keep existing API field **`name`** (maps to `full_name`) for auth and admin user DTOs for compatibility; do not rename to `fullName` in P1 JSON.
- **Rationale**: FE `AuthUser.name` and access-control already use `name`.
- **Note**: Spec “full name” is the product label; JSON stays `name`.

### R5 — `hasPassword` for OAuth-only accounts

- **Decision**: `hasPassword: password_hash != ""`. FE hides/disables change-password when false.
- **Rationale**: Matches auth login already rejecting empty hash; FR-008.

### R6 — Admin personal fields on existing user CRUD

- **Decision**: Extend `CreateUserRequest` / `UpdateUserRequest` / `UserResponse` with optional personal fields; optional avatar via separate multipart endpoint (cleaner than mixing multipart into JSON create).
- **Rationale**: FR-011/FR-012; blank password on update unchanged in `user_service.Update`.

### R7 — Admin locale flags vs embedding LocaleSelect

- **Decision**: Do **not** drop a full shadcn `Select` into the dropdown (awkward nested interactive). Extract **shared flag icons + locale option list** from `locale-select.tsx`; render two clickable flag items (or compact flag buttons) in `UserMenuContent` that call `i18n.changeLanguage('en'|'vi')` with clear active state.
- **Rationale**: Spec wants clickable flags consistent with public behavior (immediate apply, active state), not necessarily the Select chrome inside a menu.
- **Alternatives**: Reuse `LocaleSelect` inside dropdown — poor UX; keep text toggle — rejected by FR-014.

### R8 — Feature module placement (FE)

- **Decision**: New feature `fe/src/features/user-profile/` for page, forms, schemas, services. Access-control keeps admin CRUD but reuses shared personal Zod where practical. Route only in `fe/src/app/admin/profile/page.tsx`.
- **Rationale**: FE AGENTS feature-first rule; admin layout guards already protect `/admin/*`.

### R9 — GitNexus blast radius (known HIGH/CRITICAL)

| Symbol | Risk | Mitigation |
|--------|------|------------|
| `ToResponse` (user service) | HIGH | Additive JSON fields only; keep existing keys |
| `UserInfo` | CRITICAL | Add `avatarUrl` display; preserve initials fallback |
| `UserMenuContent` | HIGH | Enable settings link + flag locale; keep theme/logout |

Implementers **must** re-run `npx gitnexus impact` before editing these symbols.

---

## Data Model (Phase 1 design)

### User (extended)

| Field | DB | JSON | Rules |
|-------|-----|------|--------|
| id | uuid PK | `id` | existing |
| email | varchar | `email` | read-only on profile; editable in admin |
| password_hash | varchar | — / drives `hasPassword` | never expose hash |
| full_name | varchar | `name` | required; min length aligned with auth (≥2 BE binding / product) |
| phone | varchar(50) | `phone` | optional; trim; max 50 |
| avatar_url | text | `avatarUrl` | optional URL/path |
| general | text | `general` | optional; max 1000 |
| birthday | date | `birthday` | optional ISO date `YYYY-MM-DD`; not after today |
| address | varchar(500) | `address` | optional; max 500 |
| social_links | jsonb | `socialLinks` | array 0–5 of `{ label?: string, url: string }` |
| role_id / status / timestamps | existing | as today | not edited on self-service profile |

### SocialLink

```json
{ "label": "GitHub", "url": "https://github.com/example" }
```

- `url` required; must be `http:` or `https:`
- `label` optional; max 50
- Duplicate URLs allowed in P1

### Avatar

- Binary on disk under `UPLOAD_DIR/avatars/`
- Filename: opaque UUID + extension
- DB stores path usable by FE as `${API_BASE}/media/avatars/...` or relative `/api/media/avatars/...` (document one convention in implement verify)

### Entity relationships

```text
Auth session (JWT) --> User (self)
PATCH /auth/profile --> User (self fields)
POST /auth/profile/avatar --> file + User.avatarUrl
POST /auth/change-password --> User.passwordHash
Admin User CRUD --> User (incl. personal fields)
POST /admin/users/:id/avatar --> file + User.avatarUrl
FE AuthUser / UserInfo <-- GET /auth/me (incl. avatarUrl, hasPassword)
```

---

## Contracts

Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. All JSON **camelCase**.

### Auth / self-service

| Method | Path | Auth | Body / notes | Response |
|--------|------|------|--------------|----------|
| GET | `/auth/me` | JWT | — | `AuthUser` extended |
| PATCH | `/auth/profile` | JWT | personal fields (no email) | updated profile / user |
| POST | `/auth/profile/avatar` | JWT | `multipart/form-data` file field `file` (or `avatar`) | updated user + `avatarUrl` |
| POST | `/auth/change-password` | JWT | `{ currentPassword, newPassword }` | `{ message }` or empty success; 400/401 on bad current |

**Extended AuthUser (me / login user object):**

```json
{
  "id": "uuid",
  "email": "a@b.c",
  "name": "Display Name",
  "role": "admin",
  "roleId": "uuid",
  "permissions": ["..."],
  "phone": "",
  "avatarUrl": "",
  "general": "",
  "birthday": null,
  "address": "",
  "socialLinks": [],
  "hasPassword": true
}
```

**PATCH `/auth/profile` request:**

```json
{
  "name": "Display Name",
  "phone": "+84...",
  "general": "About text",
  "birthday": "1990-01-15",
  "address": "Full address string",
  "socialLinks": [{ "label": "Site", "url": "https://example.com" }]
}
```

Clearing optionals: send `""` / `null` / `[]` as agreed in implement (prefer empty string or null for scalars; `[]` for social links). Whitespace-only phone → store empty.

### Admin users (extended)

| Method | Path | Notes |
|--------|------|--------|
| GET/POST | `/admin/users` | Response/create include personal fields |
| PATCH/PUT | `/admin/users/:id` | Update personal fields; `password` omit/blank = keep |
| POST | `/admin/users/:id/avatar` | Optional multipart avatar for managed user |

`UserResponse` gains `phone`, `avatarUrl`, `general`, `birthday`, `address`, `socialLinks` (additive — HIGH blast radius on `ToResponse`).

### Validation errors

- Field-level messages via existing BE error envelope; FE maps to form errors.
- Avatar: 400 with clear message for type/size; previous avatar retained.

### UI routes (FE)

| Route | Guard | Page |
|-------|-------|------|
| `/admin/profile` | Protected + Admin (existing admin layout) | `user-profile` feature |
| Account settings menu | — | navigates to `/admin/profile` |

---

## Project Structure

### Documentation (this feature)

```text
docs/features/004-user-profile/
├── spec.md
├── tasks.md
├── plan.md                 ← this file
├── be-tasks-verify.md      ← @be later
├── fe-tasks-verify.md      ← @fe later
└── qa-checklist.md         ← @qa later
```

### Source Code (concrete paths)

```text
be/
├── migrations/000005_user_profile_fields.up.sql
├── migrations/000005_user_profile_fields.down.sql
├── internal/models/user/user.go
├── internal/dto/auth/auth.go
├── internal/dto/user/user.go
├── internal/services/auth/auth_service.go          # me, change-password, profile
├── internal/services/user/user_service.go          # admin CRUD personal fields
├── internal/services/media/                        # NEW: store + validate uploads
├── public/routes/auth.go
├── public/routes/user.go                           # or admin users routes
├── public/handlers/auth_handler.go
├── public/handlers/user_handler.go
├── internal/app/dependency/
├── internal/config/config.go                       # UPLOAD_DIR, public media base if needed
├── be/.env.example
└── test/unit/...

fe/src/
├── app/admin/profile/page.tsx                      # NEW
├── features/user-profile/                          # NEW
│   ├── components/profile-page.tsx
│   ├── components/profile-form.tsx
│   ├── components/avatar-upload.tsx
│   ├── components/change-password-form.tsx
│   ├── schemas/
│   ├── services/
│   ├── types.ts
│   └── index.ts
├── features/access-control/
│   ├── components/user-form-fields.tsx
│   ├── schemas/access-control-schemas.ts
│   ├── types.ts
│   └── services/access-control-api.ts
├── components/common/
│   ├── user-menu-content.tsx                       # Account settings + flags
│   ├── user-info.tsx                               # avatarUrl
│   ├── locale-select.tsx
│   └── locale-flags.tsx                            # NEW optional extract
├── types/auth.ts
└── locales/{en,vi}/...
```

**Structure Decision**: Full-stack feature; BE owns persistence + media; FE owns admin profile UX and chrome; access-control extended in place.

---

## Quickstart validation (for implementers / QA)

### Prerequisites

- `make env` / configured `be/.env`, `fe/.env`, root `.env`
- `UPLOAD_DIR` set and writable for BE process
- `make up-d` (or local BE + FE)

### Smoke scenarios

1. Sign in → user menu → **Account settings** → `/admin/profile` shows fields.
2. Edit name + optional fields → save → reload persists.
3. Upload PNG &lt; 2 MB → avatar shows; upload 3 MB → rejected, prior avatar kept.
4. Change password with correct current → login with new password works.
5. Access control → edit user personal fields + blank password → password unchanged.
6. Admin menu → click VI flag → UI Vietnamese; EN flag → English; no “Switch to …” sole control.
7. Guest opens `/admin/profile` → redirected to login.

### Commands

```bash
make test-be
make test-fe
# optional integration as available
```

---

## Env & deploy notes

| Var | Owner file | Purpose |
|-----|------------|---------|
| `UPLOAD_DIR` | `be/.env` | Absolute/relative path for avatar files |
| Optional `MEDIA_PUBLIC_BASE_URL` | `be/.env` | If absolute URLs required behind reverse proxy |
| `NEXT_PUBLIC_API_BASE_URL` | `fe/.env` | Already used; avatar URLs must resolve via this origin `/api/...` |

Docker: mount a volume onto `UPLOAD_DIR` inside the `be` service so avatars survive container recreate. Do **not** add upload paths to root `.env`.

---

## Implementation order

1. **`@be`**: migration → model/DTOs → media service → auth profile/password/avatar → admin user DTO + avatar → unit tests → `be-tasks-verify.md`
2. **`@fe`** (after contract stable; US6 can start earlier): types → profile feature + route + menu → forms/avatar/password → access-control alignment → locale flags → `fe-tasks-verify.md`
3. **`@qa`**: checklist + Independent Tests

P1 done when US1–US6 acceptance criteria met; US7 (P2) may remain open.
