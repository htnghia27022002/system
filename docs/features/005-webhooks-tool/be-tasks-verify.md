# BE tasks verify: Tools → Webhooks

**Feature:** `docs/features/005-webhooks-tool/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md), [contracts/](contracts/)  
**Agent:** `@be`  
**Date:** 2026-07-25 (RBAC seed + route guards)

## Summary

Owner webhook APIs now require RBAC: seeded `webhooks:view` / `webhooks:modify` and applied `RequireView` / `RequireModify` on `/api/webhooks/inbox*` while public capture stays open. `make test-be` passed; seeder catalog unit test updated for the new keys.

**Ops note:** Restart BE (or run seed) so `PermissionSeeder` / `RolePermissionSeeder` insert the new keys and admin role bindings — otherwise Roles UI will not list them until seed runs.

## Contract alignment (this pass)

| Contract area | Result | Notes |
|---------------|--------|-------|
| `database.md` | Match | Unchanged |
| `endpoints.md` paths / auth / status codes | Match | Owner: JWT + `webhooks:view`/`modify`; capture public; 401/403 |
| `permissions.md` | Match | Catalog IDs `…017` / `…018`; routes under `/api/webhooks/*` (not `/api/admin`) |

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T001 | Migration version + `be/.env.example` comments | Done | `000006_webhooks_inbox_*`, `be/.env.example` |
| T005 | Schema migration | Done | `be/migrations/000006_webhooks_inbox.{up,down}.sql` |
| T006 | Models | Done | `be/internal/models/webhook/{inbox,request}.go` |
| T007 | DTOs | Done | `be/internal/dto/webhook/webhook.go` |
| T008 | Repositories | Done | `interfaces/webhook.go`, `webhook_*_repository.go` |
| T009 | Service | Done | `be/internal/services/webhook/service.go` (+ contract fixes) |
| T010 | Handlers + routes + DI | Done | `webhook_handler.go`, `routes/webhook.go`, `container.go`, `dependency/` |
| T011 | Owner JWT + RBAC routes | Done | `/api/webhooks/inbox*` Auth + `RequireView`/`RequireModify("webhooks")` |
| — | Seed `webhooks:view` / `webhooks:modify` | Done | `catalog.go` IDs `…017` / `…018`; admin via `RolePermissionSeeder` |
| T012 | Public capture | Done | `Any /api/webhooks/capture/:uuid`; 404/413; ack JSON |
| T013 | Client IP | Done | `ClientIP` in handler (XFF → X-Real-IP → RemoteAddr) |
| T014 | Capture CORS | Done | `middleware.WebhookCaptureCORS` before credentialed CORS |
| T015 | Body size docs / nginx ≥ 1 MiB | Done | nginx already `20m`; documented in `be/.env.example` |
| T016 | nginx UUID rewrite | Done | `docker/nginx/config/webhooks-capture.conf` |
| T024 | Capture persistence + retention | Done | service `Capture` + `enforceRetention` |
| T025 | Ack / 404 / oversize | Done | `{ok,message}`; marker row + 413 |
| T026 | Unit tests capture/retention | Done | `be/test/unit/webhook/service_test.go` |
| T028 | List API | Done | method + `q`, page/limit, active only |
| T029 | Detail API | Done | headers/query/form/body + `utf-8`/`base64` encoding |
| T034 | Soft-delete counters | Done | `activeCount--` only; idempotent (incl. no inbox) |
| T036 | Dual counters on inbox | Done | inbox + list payloads |
| T038 | Regenerate UUID | Done | `POST /inbox/regenerate` |
| T040 | Regenerate unit test | Done | `TestRegenerateInvalidatesOldUUID` |
| T042 | Soft-delete + auth tests | Done | + `TestSoftDeleteIdempotentWithoutInbox`, `TestDetailBodyEncoding` |
| T044 | This verify doc | Done | after `make test-be` |

Open `[BE]` tasks in `tasks.md`: **none** (all marked `[x]`).

## Verification commands

```bash
make test-be
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | Pass | `be/test/unit/database` (8 catalog keys); webhook package cached |

## Key APIs

| Method | Path | Auth |
|--------|------|------|
| GET | `/api/webhooks/inbox` | JWT + `webhooks:view` |
| POST | `/api/webhooks/inbox/regenerate` | JWT + `webhooks:modify` |
| GET | `/api/webhooks/inbox/requests` | JWT + `webhooks:view` |
| GET | `/api/webhooks/inbox/requests/:id` | JWT + `webhooks:view` |
| DELETE | `/api/webhooks/inbox/requests/:id` | JWT + `webhooks:modify` |
| ANY | `/api/webhooks/capture/:uuid` | None |
| ANY | `/tools/webhooks/{uuid}` | None (nginx → capture) |

Success ack: `{ "ok": true, "message": "Request received" }`

## Acceptance coverage (BE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| Get-or-create inbox + public path | unit + API | Pass |
| Capture methods / unknown UUID | unit | Pass |
| Body > 1 MiB → 413 + marker | unit | Pass |
| Retention hard-purge at 201 | unit | Pass |
| Soft-delete counters + missing/no-inbox idempotency | unit | Pass |
| Detail bodyEncoding utf-8 / base64 / empty | unit | Pass |
| Regenerate invalidates old UUID | unit | Pass |
| Non-owner cannot see other inbox | unit | Pass |
| nginx product URL rewrite | config | Pass (manual QA on stack) |

## Gaps / follow-ups

- [x] All claimed `[BE]` tasks done
- [x] Aligned with `contracts/{database,endpoints,permissions}.md` (RBAC applied this pass)
- [ ] Restart BE / run seed so new permission rows appear in Roles UI
- Owner UI / `PermissionKeys` — **`@fe`** (verify doc: `fe-tasks-verify.md`)
- Manual Independent Tests on `system.local` — **`@qa`**

## Notes for `@fe`

- CamelCase contracts: `publicUuid`, `publicPath`, `activeCount`, `lifetimeReceived`, `bodyEncoding`, `isBinary`, `captureStatus`
- `bodyEncoding`: **`utf-8`** (text or empty body) or **`base64`** (binary); never `"none"`
- Display URL: `{origin}{publicPath}` → `/tools/webhooks/{uuid}`
- List query: `method`, `q`, `page`, `limit`
- Soft-delete detail returns 404 (no trash browser); DELETE itself is idempotent `200`
- RBAC: gate page/menu with `webhooks:view`; regenerate/delete with `webhooks:modify` (`PermissionKeys.webhooks`)
- After BE restart/seed, keys appear in Roles UI; admin role gets both via seeder

## Sign-off (BE)

- [x] All claimed `[BE]` tasks done or explicitly deferred above
- [x] Tests listed above passed
- [x] Matches `plan.md` BE sections and `contracts/*`
