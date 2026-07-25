---
name: be-develop
description: Backend development skill for Go + Gin in be/. Use when implementing or modifying API routes, handlers, services, repositories, migrations, auth/RBAC, and Docker-related BE config.
---

# BE Develop Skill

## Mandatory context first

1. Read [`be/AGENTS.md`](../../../be/AGENTS.md) for module boundaries and import rules.
2. Read [`docs-feature/SKILL.md`](../docs-feature/SKILL.md) when writing feature docs.
3. Read [`be/README.md`](../../../be/README.md) for API routes and layer boundaries.
3. Read [`AGENTS.md`](../../../AGENTS.md) for workspace boundaries and run environment.
4. Run GitNexus impact analysis before editing shared symbols (see root `AGENTS.md`).

## Module boundary

- `be/` owns its Go module (`module be` in `go.mod`).
- Import in-repo code as `be/pkg/...`, `be/internal/...`, or `be/public/...` only.
- `be/pkg/...` must not import `be/internal/...`.
- File-local aliases (e.g. `authmodel "be/internal/models/auth"`) stay inside `be/`.
- Never import from `fe/` or the monorepo root.

## Layer rules

```text
public/routes → public/handlers → internal/services → internal/repository → DB
```

| Do | Location |
|----|----------|
| Route registration | `be/public/routes/` |
| Request validation, HTTP status | `be/public/handlers/` |
| Business rules | `be/internal/services/<feature>/` |
| DB access | `be/internal/repository/` |
| Request/response shapes | `be/internal/dto/<feature>/` |
| Entities | `be/internal/models/<feature>/` |
| Hash, cache, postgres dial, query DSL | `be/pkg/` |
| JWT, errors, response | `be/internal/common/` |
| Auth middleware | `be/internal/middleware/` |
| DI orchestrator | `be/internal/app/container.go` |
| DI resolvers | `be/internal/app/dependency/` (one file per domain, e.g. `user_service.go`) |

Never put business logic in handlers or HTTP formatting in repositories.

## Queue / NATS (JetStream)

Async messaging for search outbox sync (and future events). See `be/AGENTS.md` §12 and `docs/features/002-elasticsearch-search/be-implement.md`.

```text
internal/queue/                 # NATS infra, constants, nats.json options
internal/handlers/publisher/    # Publish* after DB writes
internal/handlers/subscribers/  # process_* consumers
cmd/queue/                      # Worker process (not in API)
```

- Subject/stream/consumer names → `internal/queue/constants.go` only
- API publishes; `cmd/queue` consumes
- Do not add in-process goroutine workers in `main.go`

## API contract

- Prefix: `/api`
- JSON field names: **camelCase** (match FE types in `fe/src/features/`)
- Auth response: `{ accessToken, refreshToken, user: { id, email, name, role, roleId, permissions } }`
- Admin routes require JWT + `RequirePermission` middleware

## New feature permissions (mandatory)

For admin/protected features (see `contracts/permissions.md` and `.cursor/rules/feature-permissions.mdc`):

1. Seed keys in `be/internal/database/seeders/catalog.go` → `DefaultPermissions()` (`rbac.Key`, view/modify, stable UUIDs).
2. Protect routes with `RequireView("resource")` / `RequireModify("resource")`.
3. Rely on `RolePermissionSeeder` to attach catalog keys to the admin role.

## Run and verify

**Docker:**

```bash
make up   # from repo root
curl http://localhost:8080/api/auth/login ...
```

**Docker (stack must be up for BE tests):**

```bash
make up-d    # repo root
make test-be
make test-fe
# or
make test
```

**Local BE only** (requires Go on host):

```bash
cd be && make run && make test
```

## Env (Docker compose overrides)

- `DB_HOST=postgres`
- `REDIS_URL=redis://redis:6379`
- `NATS_ENABLED`, `NATS_URL` — queue publish/consume (`internal/queue/nats.json` for JetStream options)
- `ELASTICSEARCH_ENABLED`, `ELASTICSEARCH_URL`
- `CORS_ORIGINS` must include nginx origin (`http://localhost:8080`)

## Quality checks

```bash
make test-be && make lint   # from repo root; Docker up for test-be
```

Preserve API contracts unless the user explicitly requests breaking changes.

## OAuth providers

Follow the adapter pattern in `be/internal/services/auth/oauth/` (see `be/AGENTS.md` §10).

- One struct per provider file implementing `oauth.Provider`
- Register in `oauth/registry.go`; no provider `switch` in `OAuthService`
- Before edits: `npx gitnexus query "oauth"` and `npx gitnexus impact OAuthService`
- Shared helpers: `ValidateRedirectURI`, `resolveOAuthUser`, `syncOAuthAccount` — extend these instead of copying logic into adapters

## Tests

All BE tests live in `be/test/` (unit / integration / e2e). Never add `*_test.go` under `internal/`. See `be/test/README.md`.

```bash
make test-be && make test-be-integration && make test-be-e2e   # from repo root
```
