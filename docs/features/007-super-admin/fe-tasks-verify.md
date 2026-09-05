# FE tasks verify: Super Admin Flag

**Feature**: `007-super-admin`  
**Date**: 2026-09-05  
**Command**: `make test-fe`

## Result

**Pass.** 16 files / 53 tests.

## Coverage vs tasks

| Task | Status | Notes |
|------|--------|-------|
| T014 `superAdmin` on AuthUser / JWT hydrate | Done | `types/auth.ts`, `auth-store.ts` |
| T015 `usePermissions` Super admin allow-all | Done | + `use-permissions-super-admin.test.tsx` (53rd test) |
| T016 Users table Super admin column | Done | Desktop + mobile card |
| T017 Mock seed + JWT claim | Done | Seeded mock admin is Super admin |
| T020–T023 tabbed dialog, generate password, EN/VI | Done | Profile / Password / Account |
| T024 Mock last Super admin | Done | create/update/delete persist and protect |
| T026 Administrator role is not implicit Super admin | Done | Bypass only reads `user.superAdmin` |
| T029 `make test-fe` | Pass | |

## Notes

`superAdmin` is sent on create/update **only** when the signed-in actor is Super admin, so a member with `users:modify` does not trip the BE 403.
