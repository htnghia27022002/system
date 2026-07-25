# BE Agent Rules (Go + Gin)

This file defines strict implementation rules for AI agents working in `be/`.
Follow these rules before creating, editing, moving, or deleting files.

## 1) Self-contained module

`be/` is a **standalone Go module**. It can be moved out of this workspace without changing internal import paths.

- Module path: declared in `be/go.mod` (`module be`)
- In-module imports: `be/pkg/...`, `be/internal/...`, `be/public/...` — these refer to the **module name**, not the monorepo folder layout
- `be/pkg/...` is reusable public infra and must **never** import `be/internal/...`
- Run all Go commands from `be/` (`go run .`, `make test`, etc.)
- Do **not** import from `fe/`, repo root, or any path outside this module

When relocating the backend, move the entire `be/` directory (including `go.mod`). Update only external wiring (Docker, CI, env files) at the new location.

## 2) Source of truth

- Use `README.md` in this folder for API routes, layers, and commands
- Use `internal/app/container.go` for the DI entrypoint; per-domain wiring in `internal/app/dependency/`

## 3) Stack

Go 1.22 · Gin · GORM · PostgreSQL · JWT · OAuth2 · golangci-lint · golang-migrate

## 3.1) Migrations

- SQL files: `migrations/{version}_{name}.up.sql` and `.down.sql` via [golang-migrate](https://github.com/golang-migrate/migrate)
- Startup: `database.RunMigrations(cfg)` in `public/api.go`
- CLI: `go run ./cmd/migrate up|down|version|force|steps|drop`
- Dev seed data: `internal/database/seeders/` — `DatabaseSeeder`, `PermissionSeeder`, `RoleSeeder`, … (idempotent, runs after migrations)
- CLI: `go run ./cmd/seed` or `go run ./cmd/seed --class=PermissionSeeder`

## 4) Source layout

```text
be/
├── main.go
├── go.mod              ← module boundary; all internal imports start with `be/`
├── pkg/                # Public reusable infra (stdlib + third-party only)
│   ├── redis/
│   ├── hash/
│   ├── cache/
│   ├── postgres/
│   └── query/
├── public/
│   ├── api.go
│   ├── handlers/
│   └── routes/
├── internal/
│   ├── app/
│   │   ├── container.go
│   │   └── dependency/     # Service resolvers (infra, repos, handlers)
│   ├── config/
│   ├── database/
│   ├── handlers/           # Queue handlers (not HTTP)
│   │   ├── publisher/      # Outbound NATS publish + payload types
│   │   └── subscribers/    # Inbound message handlers
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── queue/              # NATS JetStream infra (constants, nats.json, client)
│   ├── repository/
│   ├── search/             # Elasticsearch client
│   ├── services/
│   └── common/             # App-private shared helpers (JWT, errors, response, …)
├── cmd/
│   ├── migrate/
│   ├── seed/
│   ├── queue/              # Queue worker (consumers)
│   └── reindex/            # Bulk reindex CLI
└── migrations/
```

### 4.1) `pkg/` vs `internal/`

| Path | Role |
|------|------|
| **`be/pkg/...`** | Public reusable infrastructure libraries (connection dialers, cache stores, query DSL, crypto helpers). May import **stdlib and third-party packages only**. Must **never** import `be/internal/...`. |
| **`be/internal/...`** | App-private code (config, DI, domain services, product JWT/RBAC, queue/search product wiring). May import `be/pkg/...`. |

Keep product-specific config loaders, JWT claims, middleware, and DI in `internal`. Put domain-agnostic dial/store helpers in `pkg`.

## 5) Layer rules

```text
route → handler → service → repository interface → repository → database
```

| Layer | Path |
|-------|------|
| HTTP routes | `public/routes/` |
| HTTP handlers | `public/handlers/` |
| Bootstrap | `public/api.go` |
| DI orchestrator | `internal/app/container.go` |
| DI resolvers | `internal/app/dependency/` |
| Business logic | `internal/services/<feature>/` |
| Persistence contracts | `internal/repository/interfaces/` |
| Persistence impl | `internal/repository/` |
| Models / DTOs | `internal/models/`, `internal/dto/` |
| Reusable infra | `pkg/` (no `internal` imports) |
| App-private shared helpers | `internal/common/` |
| Middleware | `internal/middleware/` |
| NATS infra | `internal/queue/` |
| Queue publish | `internal/handlers/publisher/` |
| Queue consume | `internal/handlers/subscribers/` |

Never put business logic in HTTP handlers or HTTP formatting in repositories.

## 6) Import rules

- Use the module path for all in-repo imports: `be/pkg/...`, `be/internal/...`, `be/public/...`
- `be/pkg/...` must not import `be/internal/...`
- Optional **file-local** package aliases are allowed to avoid name clashes, e.g.:

  ```go
  authmodel "be/internal/models/auth"
  usermodel "be/internal/models/user"
  ```

  Keep aliases inside `be/` only. Do not introduce workspace-level or cross-package alias conventions.
- Standard library and third-party imports come first; module imports grouped after
- Do **not** use relative import paths (`../`) between packages

## 7) API contract (FE-aligned)

- Base path: `/api`
- JSON field names: **camelCase**
- Auth: `/api/auth/login|register|refresh|logout|me`, OAuth under `/api/auth/oauth/:provider/*`
- Admin: `/api/admin/users|roles|permissions` — JWT + permission middleware

### 7.1) New feature permissions (mandatory)

When adding a protected admin/API feature:

1. Add `{resource}:view` / `{resource}:modify` to `DefaultPermissions()` in `internal/database/seeders/catalog.go` (stable UUIDs via `rbac.Key`).
2. Protect routes with `middleware.RequireView("resource")` / `RequireModify("resource")`.
3. Admin role receives new keys automatically via `RolePermissionSeeder` (iterates catalog).
4. Align with `docs/features/.../contracts/permissions.md`.

Do not ship admin routes without catalog entries. See `.cursor/rules/feature-permissions.mdc`.

Integration with the frontend is **HTTP only**. Do not import FE code or share types across repos.

## 8) Configuration (hybrid)

| Layer | Source | Contents |
|-------|--------|----------|
| Public | `be/config.yaml` | Port, JWT TTL, CORS, OAuth provider list, non-secret defaults |
| Secrets | env only | `DB_PASSWORD`, `JWT_SECRET`, `OAUTH_GOOGLE_CLIENT_ID`, `OAUTH_GOOGLE_CLIENT_SECRET` |
| Override | env optional | `DB_HOST`, `REDIS_URL`, `CORS_ORIGINS`, `CONFIG_FILE` |

Priority: **env > config.yaml > built-in default**. Never commit secrets in YAML.

| Run mode | Env file |
|----------|----------|
| **Docker stack** | repo root `.env` — compose injects into BE container; `be/.env` is **not** read |
| **Local host** | `be/.env` only — copy from `be/.env.example` |

## 9) Before finishing any task

1. Is the file placed in the correct layer folder?
2. Are imports scoped to the `be` module only?
3. Are DTOs separate from models?
4. Are errors handled via `internal/common/errors` and `response.HandleError`?
5. Is content written in English?
6. Did you run `make test-be` from repo root (or `make test` in `be/` only when Go is installed locally)?

## 10) OAuth providers (adapter pattern)

OAuth lives under `internal/services/auth/oauth/`. Each external provider is a **separate adapter**; `OAuthService` orchestrates token exchange and user linking only.

```text
handler → OAuthService → oauth.Registry → oauth.Provider (GoogleProvider, …)
                              ↓
                    shared: resolveOAuthUser, syncOAuthAccount (oauth_service.go)
```

| Piece | Path | Responsibility |
|-------|------|----------------|
| `Provider` interface | `oauth/provider.go` | Contract per provider |
| Provider adapter | `oauth/<name>_provider.go` | Credentials, oauth2 config, profile fetch |
| Registry | `oauth/registry.go` | Register adapters; lookup by provider ID |
| Shared helpers | `oauth/redirect.go`, `oauth_service.go` | Redirect validation, user/account linking |
| Orchestrator | `oauth_service.go` | Allow-list, start URL, callback, JWT issuance |

**Adding a new provider**

1. Run GitNexus before editing shared symbols: `npx gitnexus query "oauth provider"` or `npx gitnexus impact OAuthService`.
2. Create `oauth/<provider>_provider.go` implementing `Provider`.
3. Register it in `oauth/registry.go` (`NewRegistry`).
4. Add env fields in `internal/config/config.go`; document in root `.env.example` (Docker) and `be/.env.example` (local host).
5. Append provider ID to `OAUTH_ALLOWED_PROVIDERS` (comma-separated).
6. Run `make test-be`.

Do **not** add provider-specific `switch` branches in `OAuthService`; keep provider logic in the adapter file.

## 11) Tests (`test/`)

All unit, integration, and e2e tests belong under `be/test/` — **never** under `internal/` or `public/`.

| Tier | Path | Command |
|------|------|---------|
| Unit | `test/unit/` | `make test` / `make test-be` |
| Integration | `test/integration/` (`//go:build integration`) | `make test-integration` / `make test-be-integration` |
| E2E | `test/e2e/` (`//go:build e2e`) | `make test-e2e` / `make test-be-e2e` |

Shared mocks and helpers: `test/testutil/`. Full guide: [`test/README.md`](test/README.md).

Before adding tests for a feature, run `npx gitnexus query "<feature> test"` to find existing coverage and callers.

## 12) Queue / NATS (JetStream)

Async work uses **NATS JetStream** with a strict four-layer layout under `internal/`. The HTTP API publishes only; consumers run in **`cmd/queue`** (Docker service `queue`).

```text
internal/queue/              # Infrastructure — connect, EnsureInfrastructure, Publish
internal/handlers/publisher/ # Outbound — payload types + Publish* methods
internal/handlers/subscribers/  # Inbound — process_* handlers + registry
cmd/queue/                   # Separate process — bootstrap streams/consumers + consume
```

| Rule | Detail |
|------|--------|
| Names | All stream, subject, consumer, handler names in `internal/queue/constants.go` — **no string literals** elsewhere |
| JSON config | `internal/queue/nats.json` holds JetStream **options** only; keys (`search`, `search_outbox`) map to constants via `config.go` |
| Env | `NATS_ENABLED`, `NATS_URL`, optional `QUEUE_CONFIG_FILE` |
| HTTP vs worker | `be/main.go` (API) wires publisher only; **never** call `StartConsumers` in the API process |
| Services | After DB write, call `handlers/publisher` — do not import `subscribers` from services |
| Subscribers | Import `publisher` for payload decode; call `internal/services/*` for business logic |

**Adding a consumer:** constants → `nats.json` options → `resolveStream`/`resolveConsumer` → publisher method → subscriber + registry → service hook. See `docs/features/002-elasticsearch-search/be-implement.md`.

**Commands:**

```bash
go run ./cmd/queue      # queue worker
go run ./cmd/reindex    # one-shot ES reindex
```

Docker: `queue` service runs `go run ./cmd/queue`. When NATS is disabled, `cmd/queue` polls `search_outbox` instead.
