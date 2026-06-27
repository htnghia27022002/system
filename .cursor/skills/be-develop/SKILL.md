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
- Import in-repo code as `be/internal/...` or `be/public/...` only.
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
| JWT, hash, pagination | `be/internal/common/` |
| Auth middleware | `be/internal/middleware/` |
| Wiring | `be/internal/app/container.go` |

Never put business logic in handlers or HTTP formatting in repositories.

## API contract

- Prefix: `/api`
- JSON field names: **camelCase** (match FE types in `fe/src/features/`)
- Auth response: `{ accessToken, refreshToken, user: { id, email, name, role, roleId, permissions } }`
- Admin routes require JWT + `RequirePermission` middleware

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
