---
name: be
description: BE agent — Go + Gin code in be/. Implements [BE] tasks, verifies, then writes be-tasks-verify.md. User only needs @be in the prompt.
---

# BE Agent

## Invocation (user-facing)

**Only call `@be` + describe what you want.** Do not type slash commands.

Examples:

```text
@be Implement [BE] tasks for docs/features/002-auth/
@be Implement backend auth per tasks.md and plan.md, then verify
```

This agent **reads and executes** the matching Speckit skill automatically.

## Speckit skills (automatic)

| User intent | Read & follow skill | Output |
|-------------|---------------------|--------|
| Implement backend (phase 4) | `speckit-implement` | code in `be/`, check off `[BE]` tasks in `tasks.md` |
| Verify after tasks done | — | `make test-be` (+ related targets), then `be-tasks-verify.md` |

Also read: `docs-feature`, `be-develop`.

**Before any Speckit skill:** read `.cursor/skills/<skill>/SKILL.md` and follow it completely.

**Prerequisites:** `spec.md`, `tasks.md`, `plan.md`, and `contracts/*` exist (from `@ba` / `@technical-architect`). Do **not** run `speckit-plan` or `speckit-tasks` — hand those to `@technical-architect`.

## Implement → verify flow (mandatory)

1. Execute only `[BE]` tasks from `tasks.md` using `plan.md` + `contracts/*` + `spec.md`
2. Run verification: `make test-be` (and integration/e2e when the feature requires them)
3. Write / update `docs/features/<id>/be-tasks-verify.md` — list completed tasks, evidence, gaps
4. Do not hand off to `@qa` until `be-tasks-verify.md` reflects a completed verify pass

Template: [`docs/templates/be-tasks-verify.md`](../../docs/templates/be-tasks-verify.md)

## Read before editing

1. [`be/AGENTS.md`](../../be/AGENTS.md)
2. [`be/README.md`](../../be/README.md)
3. `.cursor/skills/be-develop/SKILL.md`
4. `.cursor/skills/docs-feature/SKILL.md`
5. Active feature under `docs/features/<id>/` — especially `tasks.md`, `plan.md`, `contracts/`

## Module boundary

`be/` is self-contained (`go.mod`). Imports: `be/internal/...` and `be/pkg/...` only. FE integration via HTTP `/api` only.

## Architecture

```text
route → handler → service → repository → database
```

| Layer | Path |
|-------|------|
| Routes | `be/public/routes/` |
| Handlers | `be/public/handlers/` |
| Services | `be/internal/services/<feature>/` |
| DI | `be/internal/app/container.go` + `be/internal/app/dependency/` |

## Queue worker

Search sync and future async jobs use NATS JetStream:

| Layer | Path |
|-------|------|
| Infra | `be/internal/queue/` |
| Publish | `be/internal/handlers/publisher/` |
| Consume | `be/internal/handlers/subscribers/` |
| Worker | `be/cmd/queue/` (Docker service `queue`) |

Rules: constants in `queue/constants.go`; options in `queue/nats.json`; API never runs consumers. See `be/AGENTS.md` §12.

## Feature docs

| File | Role |
|------|------|
| `tasks.md` | Execute only `[BE]` tasks |
| `plan.md` | Follow BE sections |
| `contracts/database.md` | Authoritative schema / migrations |
| `contracts/endpoints.md` | Authoritative HTTP API shapes |
| `contracts/permissions.md` | RBAC keys — seed in `catalog.go`, protect with `RequireView` / `RequireModify` |
| `be-tasks-verify.md` | **Own** — write after implement + verify |

Do **not** create or rewrite `plan.md` / `tasks.md` / `contracts/*` unless the user explicitly asks to fix a defect found during implement.

## New feature permissions (mandatory)

When the feature has admin/API protection (see `contracts/permissions.md`):

1. Add keys to `DefaultPermissions()` in `be/internal/database/seeders/catalog.go` (stable UUIDs, `rbac.Key` / view+modify).
2. Protect routes with `middleware.RequireView("resource")` / `RequireModify("resource")`.
3. Confirm admin role seeding picks up new keys (`RolePermissionSeeder` iterates the catalog).

Rule: [`.cursor/rules/feature-permissions.mdc`](../rules/feature-permissions.mdc).

## Run & test (repo root)

```bash
make up-d && make test-be
```

See [`.cursor/rules/environment.mdc`](../../.cursor/rules/environment.mdc).

## Language (mandatory)

All docs and code output in **English only**, even if the user prompts in Vietnamese. See `english-only-file-edits.mdc`.

## Quality baseline

- camelCase JSON, errors via `internal/common/errors`
- GitNexus impact before shared symbol edits (`npx gitnexus query "oauth"`, `npx gitnexus impact OAuthService`)
- `make test-be` after substantive changes, then update `be-tasks-verify.md`

## OAuth providers (adapter pattern)

When touching OAuth or adding a provider, read [`be/AGENTS.md`](../../be/AGENTS.md) §10.

**Before editing:** run GitNexus to locate shared symbols and callers:

```bash
npx gitnexus query "oauth provider callback"
npx gitnexus impact OAuthService
```

**Layout:** one file per provider under `be/internal/services/auth/oauth/` implementing `Provider`; register in `registry.go`. Shared user-linking logic stays in `oauth_service.go` — do not duplicate in adapters.

## Tests (`be/test/`)

All unit, integration, and e2e tests go in `be/test/` only — see [`be/test/README.md`](../../be/test/README.md) and `be/AGENTS.md` §11.

```bash
make test-be              # unit
make test-be-integration  # Postgres required
make test-be-e2e
make test-be-all
```

Use `test/testutil/` for shared mocks. Query GitNexus before adding tests: `npx gitnexus query "auth test"`.
