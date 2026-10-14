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
| Account (email links) | `POST /auth/password/forgot` (always 202), `/auth/password/reset`, `/auth/email/verify`, `/auth/email/verification` (resend, signed in), `/auth/invite/inspect`, `/auth/invite/accept`; admin `POST /admin/users/invite`, `/admin/users/:id/invite` (resend) |
| Admin | CRUD `/admin/users`, `/admin/roles`; `GET /admin/permissions` |
| Admin system | `GET /admin/dashboard/overview` (`dashboard:view`), `GET /admin/audit-logs` (`audit:view`; filters `search`, `action`, `actorId`, `resourceType`, `resourceId`, `from`, `to`), `GET`/`PATCH /admin/settings` (`settings:view` / `settings:modify`) |
| Media | `GET /media/*key` (public) streams local uploads with a sandboxing CSP |
| App | `GET /app/info` → `{"modules": [...], "settings": {public settings}}` (public); `GET /health` (liveness), `GET /ready` (DB + enabled Redis/NATS, 503 when down) |
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

## Security & operations

| Concern | Behavior | Env |
|---------|----------|-----|
| Logging | `log/slog`, one structured access line per request (no query string), panics logged and answered with a generic 500 | `LOG_FORMAT=text\|json`, `LOG_LEVEL` |
| Request ID | `X-Request-ID` reused when well-formed, else generated; echoed on the response and attached to error logs | — |
| Security headers | `nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: no-referrer`, COOP; HSTS only over HTTPS | — |
| Client IP | `X-Forwarded-For` / `X-Forwarded-Proto` are honored only from trusted proxies (used for rate limits, sessions, HTTPS detection) | `TRUSTED_PROXIES` (default loopback + private ranges) |
| Rate limit | Per IP + route, fixed window, in-memory (per instance) on login, register, OAuth callback, change-password → 429 + `Retry-After` | `AUTH_RATE_LIMIT_PER_MINUTE` (0 = off) |
| Account lockout | DB-backed (`login_attempts`): N wrong passwords lock password login for a duration → 429 | `AUTH_MAX_FAILED_LOGINS` (0 = off), `AUTH_LOCKOUT_DURATION` |
| Refresh token | httpOnly cookie `refresh_token`, path `/api/auth`; not in JSON. Cookie-backed `/refresh` and `/logout` require `X-Requested-With` (CSRF guard). Body tokens still work for non-browser clients; the new token returns on the same channel | `AUTH_REFRESH_COOKIE`, `AUTH_COOKIE_SAMESITE` (`none` when FE is on another site; forces Secure), `AUTH_COOKIE_SECURE`, `AUTH_COOKIE_DOMAIN` |

### Audit, settings, storage

- **Audit trail** (`internal/audit`, table `audit_logs`): services call `Recorder.Record(ctx, audit.Entry{Action: "user.update", ...})`; actor, IP, user agent and request ID come from the request context (`middleware.AuditContext` + `middleware.Auth`). Never put secrets in `Metadata`. Retention: setting `audit.retentionDays` (janitor runs every 6 h, also prunes dead tokens).
- **Settings** (`internal/settings`, table `app_settings`): declare `settings.Definition`s in code (core: `CoreDefinitions`, modules: `SettingsProvider`); the DB only holds overrides; reads are cached 30 s per process. Unknown keys and wrong types are rejected.
- **Storage** (`pkg/storage`): `STORAGE_DRIVER=local` (files under `UPLOAD_DIR`, URLs `STORAGE_PUBLIC_URL` + `/key`, default `/api/media`) or `s3` (`S3_ENDPOINT`, `S3_BUCKET`, `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`, `S3_PUBLIC_URL`, `S3_REGION`, `S3_USE_SSL`, `S3_FORCE_PATH_STYLE`). Objects are public; keep private files elsewhere. With an FE on another origin, set an absolute `STORAGE_PUBLIC_URL`.

### Mail & account lifecycle

`internal/mail` sends transactional mail through `MAIL_DRIVER=log` (default — the message, including its link, is written to the BE log) or `smtp` (`SMTP_HOST`, `SMTP_PORT`, `SMTP_TLS=starttls|tls|none`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `MAIL_FROM`). Sending is asynchronous; a bad driver config stops startup. Links point at `APP_URL` (the FE origin): `/reset-password`, `/verify-email`, `/accept-invite`, each with a single-use token (stored hashed in `user_tokens`; reset 1 h, verify 48 h, invite 7 days; only the newest link of a kind works). Modules can send mail via `Kernel.Mailer`. Dev: `make up-d PROFILES="dev mail"` starts Mailpit (UI on `:8025`).

Probe `/api/ready` from orchestrators (compose healthchecks already do); `/api/health` only says the process is up.

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
