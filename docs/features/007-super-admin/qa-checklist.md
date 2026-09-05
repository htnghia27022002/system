# QA checklist: Super Admin Flag and Full User Edit

**Feature**: `007-super-admin`  
**Date**: 2026-09-05  
**Commands**: `make test-be`, `make test-fe`

## Automated

- [x] `make test-be` — pass (includes middleware Super admin bypass + last Super admin)
- [x] `make test-fe` — pass (16 files / 53 tests, including Super admin allow-all)

## Persistence

- [x] Migration `000009` applied: `users.is_super_admin` exists
- [x] Seeded `admin@example.com` is Super admin (`t`); `demo@example.com` is not (`f`)
- [x] `be` restarted after code change (dev is `go run .`)

## Independent Tests (manual)

Browser MCP tools were **not available** in this session. Confirm after a **hard refresh** (and re-login if the access token predates the new `superAdmin` claim — `/auth/me` also returns the flag).

1. Sign in as `admin@example.com`. Open `/admin/users`. Super admin column shows a badge on Admin User.
2. Edit Admin User (or another user). Confirm tabs **Profile**, **Password**, **Account**. Account tab has Super admin checkbox.
3. Password tab: Generate password fills the field.
4. Edit a member’s profile fields and save. Reopen and confirm values.
5. Attempt to clear Super admin on the last Super admin — API/UI must reject.
6. Sign in as a member without `users:view` — Users page and API stay forbidden.

## Out of P1

- Admin viewing another user’s sessions (006 P2)
