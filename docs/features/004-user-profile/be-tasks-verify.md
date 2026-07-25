# BE tasks verify: User Profile (004-user-profile)

**Feature:** `docs/features/004-user-profile/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@be`  
**Date:** 2026-07-25

## Summary

Backend profile support is implemented: migration `000005` adds personal columns; self-service APIs under `/api/auth/*` cover extended `/me`, profile PATCH, avatar upload, and change-password; admin user CRUD returns/accepts the same personal fields with optional avatar upload; local avatar storage uses `UPLOAD_DIR` and is served at `GET /api/media/avatars/:filename`. Unit tests pass via `make test-be`.

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T001 | Confirm migration version + document upload env in `be/.env.example` | Done | `000005_*`; `UPLOAD_DIR` in `be/.env.example` |
| T004 | Migration `general`, `birthday`, `address`, `social_links` | Done | `be/migrations/000005_user_profile_fields.*.sql` |
| T005 | Extend user model | Done | `be/internal/models/user/user.go` |
| T006 | Personal-field DTOs + validation | Done | `be/internal/dto/user/personal.go` |
| T007 | Local avatar storage service | Done | `be/internal/services/media/avatar.go` |
| T008 | Media serving `/api/media/...` | Done | `be/public/routes/media.go`, `media_handler.go` |
| T009 | Extend `AuthUserResponse` + `buildAuthUser` (+ `hasPassword`) | Done | `be/internal/dto/auth/auth.go`, `auth_service.go` |
| T010 | Self-service routes + DI | Done | `auth.go` routes; handlers; `dependency/` |
| T011 | Admin user DTOs / `ToResponse` personal fields | Done | `dto/user/user.go`, `user_service.go` (additive; GitNexus HIGH on `ToResponse`) |
| T012 | Admin `POST /admin/users/:id/avatar` | Done | `admin.go`, `user_handler.go` |
| T019 | `PATCH /auth/profile` service logic | Done | `auth_service.UpdateProfile` |
| T023 | Finalize self avatar upload | Done | `auth_service.UploadAvatar` + handler |
| T026 | `POST /auth/change-password` | Done | `auth_service.ChangePassword` |
| T029 | Admin create/update persist personal fields | Done | `user_service.Create` / `Update` |
| T039 | Unit tests for validation / password / avatar | Done | `test/unit/user`, `auth`, `media` |

## Verification commands

```bash
make test-be
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | Pass | All `./test/unit/...` packages OK |

## Acceptance coverage (BE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| Extended `/me` profile + `hasPassword` | Unit + API surface | Pass |
| Profile PATCH personal fields / validation | Unit (`ValidatePersonalFields`, `UpdateProfile`) | Pass |
| Avatar type/size reject; prior URL kept | Unit media + auth upload | Pass |
| Change password (success / wrong current / OAuth-only) | Unit | Pass |
| Admin personal fields + blank password keep | Unit `ToResponse`, `UpdateKeepsPasswordWhenBlank` | Pass |
| Media GET for opaque avatar filenames | Route + resolve helpers | Pass (manual smoke optional) |

## Env / contract notes for `@fe`

- `UPLOAD_DIR` in `be/.env` (default `data/uploads`); never root `.env`
- Avatar URLs stored as **`/api/media/avatars/<uuid>.<ext>`** (relative to site origin behind nginx `/api` proxy)
- Auth JSON keeps field **`name`** (not `fullName`); new fields: `phone`, `avatarUrl`, `general`, `birthday`, `address`, `socialLinks`, `hasPassword`
- Multipart field: `file` or `avatar`
- Change password body: `{ currentPassword, newPassword }`

## Gaps / follow-ups

- [x] None blocking for `@fe` against the BE contract
- [ ] Integration/e2e against live Postgres + multipart not run in this pass (`make test-be` unit only)
- [ ] Docker named volume for uploads optional; bind-mount `./be` already persists `data/uploads` on host

## Sign-off (BE)

- [x] All claimed `[BE]` tasks done
- [x] `make test-be` passed
- [x] Matches `plan.md` BE sections (additive `ToResponse` / auth user fields; HIGH blast radius accepted as additive-only)
