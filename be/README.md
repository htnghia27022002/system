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
│   ├── app/
│   │   ├── container.go      # DI orchestrator (public Container API)
│   │   └── dependency/       # Service resolvers (repos + wiring per domain)
│   ├── config/
│   ├── database/
│   ├── handlers/             # Queue publisher + subscribers (not HTTP)
│   ├── queue/                # NATS JetStream (constants, nats.json, client)
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── repository/
│   ├── search/               # Elasticsearch client
│   ├── services/
│   └── common/
├── cmd/
│   ├── migrate/
│   ├── seed/
│   ├── queue/                # NATS consumer worker
│   └── reindex/              # Bulk search reindex CLI
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

Go 1.22 · Gin · GORM · PostgreSQL · JWT · OAuth2 · Elasticsearch · NATS JetStream · golangci-lint

## API routes

Base URL: `/api` (via nginx `http://localhost:8080/api` in Docker, or `http://localhost:8080/api` local)

| Group | Paths |
|-------|-------|
| Auth | `POST /auth/login`, `/register`, `/refresh`, `/logout`; `GET /auth/me`; OAuth under `/auth/oauth/:provider/*` |
| Admin | CRUD `/admin/users`, `/admin/roles`; `GET /admin/permissions`; `GET /admin/search` (+ reindex / outbox admin) |

JSON responses use **camelCase**. Admin routes require Bearer JWT + permissions.

Full route table and permissions: see existing handlers in `public/routes/`.

## Seed data

| Email | Password | Role |
|-------|----------|------|
| `admin@example.com` | `admin1234` | Administrator |
| `demo@example.com` | `password123` | Member |

Seeded on first startup via `internal/database/seeders/` (Laravel-style) after [golang-migrate](https://github.com/golang-migrate/migrate) SQL migrations.

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
make migrate-down     # roll back one migration
make migrate-version  # print current migration version
make migrate-create name=add_users_index  # scaffold up/down files
make seed             # run DatabaseSeeder
make seed-class class=PermissionSeeder  # run one seeder class
go run ./cmd/queue        # NATS queue worker (consumers; separate from API)
go run ./cmd/reindex      # one-shot Elasticsearch bulk reindex
```

## Cache

Redis dial: `be/pkg/redis.Connect(url)`; app wrappers in `internal/database` (`ConnectRedis` / `ConnectRedisURL`).

Cache stores live in `be/pkg/cache` (Options + drivers). App Init/Default wiring stays in `internal/common/cache` and maps `config.CacheConfig` → pkg Options. Disabled by default (`cache.enabled: false`).

| Driver | Config | Notes |
|--------|--------|-------|
| `file` | `cache.file.dir` (default `storage/cache`) | Local JSON files under `be/storage/cache` |
| `redis` | `cache.redis.url` or top-level `redis.url` | Uses Redis DB from URL |

Package API: `cache.Init(cfg.Cache, redisClient)` at startup, then `cache.Get`, `cache.Set`, `cache.Delete`, `cache.Purge`, `cache.Close`. Redis client comes from `database.ConnectRedisURL` in `main.go`.

## Queue (NATS JetStream)

Search index sync uses a **transactional outbox** in PostgreSQL and NATS JetStream for async processing.

| Layer | Path | Role |
|-------|------|------|
| Infra | `internal/queue/` | Connect, stream/consumer bootstrap, publish |
| Constants | `internal/queue/constants.go` | Stream, subject, consumer, handler names |
| Options | `internal/queue/nats.json` | Retention, ack policy, etc. (not names) |
| Publish | `internal/handlers/publisher/` | Outbound messages after DB write |
| Consume | `internal/handlers/subscribers/` | Inbound handlers (`process_search`, …) |
| Worker | `cmd/queue/` | Runs consumers (Docker service `queue`) |

The API process publishes only. **`cmd/queue`** calls `EnsureInfrastructure` and `StartConsumers` on startup.

Env: `NATS_ENABLED`, `NATS_URL`, optional `QUEUE_CONFIG_FILE` (default `internal/queue/nats.json`).

When NATS is disabled, the publisher no-ops and `cmd/queue` polls the outbox table.

Design reference: [`docs/features/002-elasticsearch-search/be-implement.md`](../docs/features/002-elasticsearch-search/be-implement.md).

From monorepo root: `make test-be`, `make test-be-integration`, `make test-be-e2e`, `make test-be-all`.

See [`test/README.md`](test/README.md) for layout and conventions.

## Dependency injection

Wiring lives under `internal/app/`. **`container.go`** is a thin orchestrator; each domain resolver in **`internal/app/dependency/`** constructs its repositories and related services.

```text
internal/app/
├── container.go              # NewContainer — exposes handlers + services to routes/cmd
└── dependency/
    ├── infra.go              # Queue, JWT, publisher, Elasticsearch client
    ├── repositories.go       # newAuthRepository, newUserRepository, …
    ├── search.go             # Outbox, index processor, search query service
    ├── auth_service.go       # Auth + OAuth
    ├── user_service.go       # UserService (+ user repo)
    ├── role_service.go       # RoleService (+ role repo for auth middleware)
    ├── permission_service.go
    └── handlers.go           # HTTP handler constructors
```

**Resolve order in `NewContainer`:** `Infra` → `SearchStack` (shared outbox) → domain services → `HTTPHandlers`.

Add a new domain service by adding `dependency/<feature>_service.go` with a `New…` constructor, then register it in `container.go`. Keep business logic in `internal/services/<feature>/`; resolvers only wire dependencies.

## Conventions

- HTTP in `public/handlers`, routes in `public/routes`
- Queue publish/consume in `internal/handlers/publisher` and `internal/handlers/subscribers`
- NATS infra in `internal/queue/` — subject names from `constants.go` only
- Business logic in `internal/services/<feature>`
- DTOs separate from models
- Errors via `internal/common/errors` + `response.HandleError`

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
