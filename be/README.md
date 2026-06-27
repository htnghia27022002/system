# BE — Go + Gin API

Go backend with JWT auth, RBAC, user management, and OAuth2 SSO.

## Architecture

```text
route → handler → service → repository interface → repository → database
```

## Source layout

```text
be/
├── config.yaml           # Public config (committed)
├── main.go
├── public/
│   ├── api.go
│   ├── handlers/
│   └── routes/
├── internal/
│   ├── app/container.go      # DI wiring
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── repository/
│   ├── services/
│   └── common/
├── migrations/
├── test/                     # All unit, integration, and e2e tests
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   └── testutil/
├── Makefile
└── .golangci.yml
```

## Stack

Go 1.22 · Gin · GORM · PostgreSQL · JWT · OAuth2 · golangci-lint

## API routes

Base URL: `/api` (via nginx `http://localhost:8080/api` in Docker, or `http://localhost:8080/api` local)

| Group | Paths |
|-------|-------|
| Auth | `POST /auth/login`, `/register`, `/refresh`, `/logout`; `GET /auth/me`; OAuth under `/auth/oauth/:provider/*` |
| Admin | CRUD `/admin/users`, `/admin/roles`; `GET /admin/permissions` |

JSON responses use **camelCase**. Admin routes require Bearer JWT + permissions.

Full route table and permissions: see existing handlers in `public/routes/`.

## Seed data

| Email | Password | Role |
|-------|----------|------|
| `admin@example.com` | `admin1234` | Administrator |
| `demo@example.com` | `password123` | Member |

Seeded on first startup via `internal/database/seed.go` (GORM AutoMigrate).

## Environment

**Public config:** [`config.yaml`](config.yaml) — port, JWT TTL, CORS, OAuth provider list, non-secret defaults (safe to commit).

**Secrets & overrides:** copy [`be/.env.example`](.env.example) to `be/.env` for local host runs, or set env in Docker compose. Env **always wins** over `config.yaml`.

| Source | Examples |
|--------|----------|
| `config.yaml` | `server.port`, `jwt.accessTtl`, `cors.origins`, `oauth.allowedProviders` |
| Env (secrets) | `DB_PASSWORD`, `JWT_SECRET`, `OAUTH_GOOGLE_CLIENT_ID`, `OAUTH_GOOGLE_CLIENT_SECRET` |
| Env (deploy override) | `DB_HOST=postgres`, `REDIS_URL`, `CORS_ORIGINS`, `CONFIG_FILE` |

**Docker:** configured in root `docker-compose.yml` — see [`docker/README.md`](../docker/README.md).

## Commands

```bash
make run              # start server
make test             # unit tests (test/unit/)
make test-integration # integration tests (needs Postgres)
make test-e2e         # HTTP e2e tests (needs Postgres)
make test-all         # all test tiers
make lint             # golangci-lint
make migrate-up       # SQL migration instructions
```

From monorepo root: `make test-be`, `make test-be-integration`, `make test-be-e2e`, `make test-be-all`.

See [`test/README.md`](test/README.md) for layout and conventions.

## Conventions

- HTTP in `public/handlers`, routes in `public/routes`
- Business logic in `internal/services/<feature>`
- DTOs separate from models
- Errors via `internal/common/errors` + `response.HandleError`

## Module imports

This directory is a standalone Go module (`module be` in `go.mod`). Internal imports use the module path:

```go
import (
    "be/internal/common/hash"
    authmodel "be/internal/models/auth"
)
```

`be/` here is the **module name**, not a reference to the monorepo root. The whole `be/` folder can be relocated; internal imports stay valid as long as `go.mod` moves with it. See [`AGENTS.md`](AGENTS.md) for full agent rules.

## Related docs

- Monorepo overview: [`README.md`](../README.md)
- Docker stack: [`docker/README.md`](../docker/README.md)
- Agent rules: [`AGENTS.md`](../AGENTS.md)
