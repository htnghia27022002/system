# BE — Go + Gin API

Go backend with JWT auth, RBAC, user management, and OAuth2 SSO — the **core** of the
base. Optional features plug in as **modules** (shipped: `webhooks`); see
[`docs/modules.md`](../docs/modules.md). This folder builds and deploys on its own
(`go.mod`, `Dockerfile`, `.env`).

## Architecture

```text
route → handler → service → repository interface → repository → database
```

## Source layout

```text
be/
├── config.yaml           # Public config (committed)
├── main.go
├── common/               # App-shared helpers: jwt, errors, response, httpx, rbac, cache, utils
├── pkg/                  # Reusable infra: postgres, query, repo, hash, cache, redis
├── public/
│   ├── api.go
│   ├── handlers/
│   └── routes/
├── internal/
│   ├── app/
│   │   ├── container.go      # Core DI + Kernel + enabled modules
│   │   └── dependency/       # Core service resolvers (repos + wiring per domain)
│   ├── module/               # Module contract (Module, Kernel, capabilities)
│   ├── modules/              # registry.go + one folder per optional module (webhooks/)
│   ├── config/
│   ├── database/
│   ├── handlers/publisher/   # Queue publisher (not HTTP)
│   ├── queue/                # NATS JetStream client + Stream/Consumer topology
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── repository/
│   ├── services/
│   └── common/
├── cmd/
│   ├── migrate/
│   ├── seed/
│   └── queue/                # Queue worker (module streams / pollers)
├── migrations/               # Core schema (module SQL lives in each module)
├── test/                     # All unit, integration, and e2e tests
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   └── testutil/
├── Dockerfile / Dockerfile.dev
├── Makefile
└── .golangci.yml
```

## Stack

Go 1.22 · Gin · pgx · squirrel · PostgreSQL · JWT · OAuth2 · golangci-lint · optional Redis cache and NATS JetStream

## API routes

Base URL: `/api` (via nginx `http://localhost:8080/api` in Docker, or `http://localhost:8080/api` local)

| Group | Paths |
|-------|-------|
| Auth | `POST /auth/login`, `/register`, `/refresh`, `/logout`; `GET /auth/me`; OAuth under `/auth/oauth/:provider/*` |
| Admin | CRUD `/admin/users`, `/admin/roles`; `GET /admin/permissions` |
| App | `GET /app/info` → `{"modules": [...]}` (public) |
| Module: webhooks | `ANY /webhooks/capture/:uuid` (public); owner inbox under `/webhooks/inbox*` (`webhooks:view` / `webhooks:modify`) |

JSON responses use **camelCase**. Admin routes require Bearer JWT + permissions.

Full route table and permissions: see existing handlers in `public/routes/`.

## Seed data

| Email | Password | Role |
|-------|----------|------|
| `admin@example.com` | `admin1234` | Administrator |
| `demo@example.com` | `password123` | Member |

Seeded on first startup via `internal/database/seeders/` (Laravel-style) after [golang-migrate](https://github.com/golang-migrate/migrate) SQL migrations. Migrations come in sets: core (`migrations/`, table `schema_migrations`) and one per enabled module (table `schema_migrations_<module>`).

## Modules

`MODULES` in `.env` picks the optional modules: unset = registry defaults (`webhooks`), `MODULES=` = core only, a comma list = exactly those. Disabled modules mount no routes, run no migrations and seed no permissions. Add/remove guide: [`docs/modules.md`](../docs/modules.md).

## Environment

| Run mode | Env file |
|----------|----------|
| **Docker / queue / prod** | `be/.env` only (`env_file`) |
| **Local host** | same `be/.env` — set `DB_HOST=localhost`, local service URLs |

**Public config:** [`config.yaml`](config.yaml). **Everything else:** `be/.env`.

## Commands

```bash
make run              # start server
make test             # unit tests (test/unit/)
make test-integration # integration tests (needs Postgres)
make test-e2e         # HTTP e2e tests (needs Postgres)
make test-all         # all test tiers
make lint             # golangci-lint
make migrate-up       # apply SQL migrations (golang-migrate)
make migrate-down     # roll back one core migration (module=<name> for a module set)
make migrate-version  # print the version of core and every enabled module
make migrate-create name=add_users_index [module=<name>]  # scaffold up/down files
make seed             # run DatabaseSeeder
make seed-class class=PermissionSeeder  # run one seeder class
go run ./cmd/queue        # queue worker (module streams / pollers; idles when none)
```

## Cache

Redis dial: `be/pkg/redis.Connect(url)`; app wrappers in `internal/database` (`ConnectRedis` / `ConnectRedisURL`).

Cache stores live in `be/pkg/cache` (Options + drivers). App Init/Default wiring stays in `common/cache` and maps `config.CacheConfig` → pkg Options. Disabled by default (`cache.enabled: false`).

| Driver | Config | Notes |
|--------|--------|-------|
| `file` | `cache.file.dir` (default `storage/cache`) | Local JSON files under `be/storage/cache` |
| `redis` | `cache.redis.url` or top-level `redis.url` | Uses Redis DB from URL |

Package API: `cache.Init(cfg.Cache, redisClient)` at startup, then `cache.Get`, `cache.Set`, `cache.Delete`, `cache.Purge`, `cache.Close`. Redis client comes from `database.ConnectRedisURL` in `main.go`.

## Queue (NATS JetStream) — optional

Off by default (`NATS_ENABLED=false`). Modules declare JetStream streams and consumers (`module.StreamProvider`) and may offer a `Poller` for NATS-disabled setups.

| Layer | Path | Role |
|-------|------|------|
| Infra | `internal/queue/` | Connect, ensure declared streams/consumers, publish |
| Options | `internal/queue/nats.json` | Connection + optional per-stream overrides (defaults live in modules) |
| Publish | `internal/handlers/publisher/` | `PublishJSON` after the DB write (no-op when disabled) |
| Worker | `cmd/queue/` | Consumes module streams, polls module pollers, or idles |

Env: `NATS_ENABLED`, `NATS_URL`, optional `QUEUE_CONFIG_FILE` (default `internal/queue/nats.json`). Docker: `--profile queue`.

From monorepo root: `make test-be`, `make test-be-integration`, `make test-be-e2e`, `make test-be-all`.

See [`test/README.md`](test/README.md) for layout and conventions.

## Dependency injection

Wiring lives under `internal/app/`. **`container.go`** is a thin orchestrator; each domain resolver in **`internal/app/dependency/`** constructs its repositories and related services.

```text
internal/app/
├── container.go              # NewContainer — exposes handlers + services to routes/cmd
└── dependency/
    ├── infra.go              # Queue, JWT, publisher
    ├── repositories.go       # newAuthRepository, newUserRepository, …
    ├── auth_service.go       # Auth + OAuth
    ├── user_service.go       # UserService (+ user repo)
    ├── role_service.go       # RoleService (+ role repo for auth middleware)
    ├── permission_service.go
    └── handlers.go           # HTTP handler constructors
```

**Resolve order in `NewContainer`:** `Infra` → core services → `HTTPHandlers` → `Kernel` → enabled modules (`internal/modules.Build`).

Core services get a `dependency/<feature>_service.go` resolver registered in `container.go`. Optional features do their own wiring in `internal/modules/<name>/module.go`.

## Conventions

- HTTP in `public/handlers`, routes in `public/routes`
- Optional features in `internal/modules/<name>/` (same layers inside); core never imports a module
- Queue publish via `internal/handlers/publisher`; stream/subject names are constants in the owning module
- Business logic in `internal/services/<feature>`
- DTOs separate from models
- Errors via `common/errors` + `response.HandleError`

## Module imports

This directory is a standalone Go module (`module be` in `go.mod`). Imports use the module path:

```go
import (
    "be/pkg/hash"
    authmodel "be/internal/models/auth"
)
```

`be/pkg/...` is reusable infra (stdlib + third-party only; never imports `be/internal`). `be/` here is the **module name**, not a reference to the monorepo root. See [`AGENTS.md`](AGENTS.md) for full agent rules.

## Related docs

- Monorepo overview: [`README.md`](../README.md)
- Docker stack: [`docker/README.md`](../docker/README.md)
- Agent rules: [`AGENTS.md`](../AGENTS.md)
