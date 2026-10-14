# BE Agent Rules (Go + Gin)

This file defines strict implementation rules for Claude when working in `be/`. Claude Code loads it automatically for any file under `be/`; project-wide rules live in `.claude/rules/`, skills in `.claude/skills/` (`be-develop`), and subagents in `.claude/agents/` (`be`, `be-test`).
Follow these rules before creating, editing, moving, or deleting files.

## 1) Self-contained module

`be/` is a **standalone Go module**. It can be moved out of this workspace without changing internal import paths.

- Module path: declared in `be/go.mod` (`module be`)
- In-module imports: `be/pkg/...`, `be/common/...`, `be/internal/...`, `be/public/...` — these refer to the **module name**, not the monorepo folder layout
- `be/pkg/...` is reusable public infra and must **never** import `be/internal/...`
- Run all Go commands from `be/` (`go run .`, `make test`, etc.)
- Do **not** import from `fe/`, repo root, or any path outside this module

When relocating the backend, move the entire `be/` directory (including `go.mod`). Update only external wiring (Docker, CI, env files) at the new location.

## 2) Source of truth

- Use `README.md` in this folder for API routes, layers, and commands
- Use `internal/app/container.go` for the DI entrypoint; per-domain core wiring in `internal/app/dependency/`
- Optional features are **modules** — contract in `internal/module/`, list in `internal/modules/registry.go`, one folder per module. Full guide: [`docs/modules.md`](../docs/modules.md)

## 3) Stack

Go 1.22 · Gin · pgx · squirrel · PostgreSQL · JWT · OAuth2 · golangci-lint · golang-migrate

## 3.1) Migrations

- SQL files: `migrations/{version}_{name}.up.sql` and `.down.sql` via [golang-migrate](https://github.com/golang-migrate/migrate), embedded with `go:embed`
- Sets: **core** (`be/migrations`, table `schema_migrations`, starts at `000001_core`) and one set per module (`internal/modules/<name>/migrations`, table `schema_migrations_<name>`). Never put a module table in core migrations.
- Startup: `database.RunMigrations(cfg, modules.MigrationSets()...)` + seeders in `public/api.go`, gated by `AUTO_MIGRATE` (default `true`). Multi-replica deploys set `AUTO_MIGRATE=false` and run `cmd/migrate` + `cmd/seed` as a release step.
- CLI: `go run ./cmd/migrate up|version|drop` (all sets) and `go run ./cmd/migrate -module <core|name> down|force|steps`
- Dev seed data: `internal/database/seeders/` — `DatabaseSeeder`, `PermissionSeeder`, `RoleSeeder`, … (idempotent, runs after migrations); modules add permissions/seeders through `seeders.Contribution`
- CLI: `go run ./cmd/seed` or `go run ./cmd/seed --class=PermissionSeeder`

## 4) Source layout

```text
be/
├── main.go
├── go.mod              ← module boundary; all internal imports start with `be/`
├── common/             # App-shared helpers (JWT, errors, response, httpx, rbac, cache, utils)
├── pkg/                # Public reusable infra (stdlib + third-party only)
│   ├── redis/
│   ├── hash/
│   ├── cache/
│   ├── postgres/
│   ├── query/
│   └── repo/           # Generic pgx + squirrel CRUD
├── public/
│   ├── api.go
│   ├── handlers/
│   └── routes/
├── internal/
│   ├── app/
│   │   ├── container.go    # core DI + Kernel + enabled modules
│   │   └── dependency/     # Core service resolvers (infra, repos, handlers)
│   ├── module/             # Module contract (Module, Kernel, capabilities, Set)
│   ├── modules/
│   │   ├── registry.go     # the only list of modules (+ MODULES env)
│   │   └── webhooks/       # reference module: handlers/ services/ repository/ models/ dto/ migrations/ module.go
│   ├── config/
│   ├── database/
│   ├── handlers/           # Queue side (not HTTP)
│   │   └── publisher/      # Outbound NATS publish (PublishJSON)
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── queue/              # NATS JetStream infra (client, nats.json, Stream/Consumer topology)
│   ├── repository/
│   │   └── interfaces/     # Persistence contracts (feature-specific)
│   ├── services/
│   └── common/             # App-private shared helpers (JWT, errors, response, httpx, …)
├── cmd/
│   ├── migrate/
│   ├── seed/
│   └── queue/              # Queue worker (module streams / pollers)
└── migrations/
```

### 4.1) `pkg/` vs `internal/`

| Path | Role |
|------|------|
| **`be/pkg/...`** | Public reusable infrastructure libraries (connection dialers, cache stores, query DSL, crypto helpers). May import **stdlib and third-party packages only**. Must **never** import `be/internal/...`. |
| **`be/internal/...`** | App-private code (config, DI, domain services, product JWT/RBAC, queue wiring, modules). May import `be/pkg/...`. |

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
| Persistence impl | `internal/repository/` (embed `pkg/repo.Repository[T]`; query via `pkg/query`) |
| Models / DTOs | `internal/models/`, `internal/dto/` (paginated lists: `query.Page[T]`) |
| Reusable infra | `pkg/` (no `internal` imports) |
| App-shared helpers | `common/` (HTTP bind/error helpers in `httpx`; format helpers in `utils`) |
| Middleware | `internal/middleware/` |
| NATS infra | `internal/queue/` |
| Queue publish | `internal/handlers/publisher/` |
| Queue consume | module `Streams()` consumers, run by `cmd/queue` |
| Optional features | `internal/modules/<name>/` (same layers inside the module folder) |

Never put business logic in HTTP handlers or HTTP formatting in repositories.

### 5.1) SOLID and reusable helpers

Keep the layered layout above. Do **not** introduce extra top-level trees (`biz`/`data`/`usecase`). The one exception is **`internal/modules/<name>/`**: an optional feature lives there as a vertical slice that repeats the same layers (`handlers/`, `services/`, `repository/`, `models/`, `dto/`, `migrations/`) so it can be deleted as one folder. Core features (auth, users, roles, permissions, media) stay in the shared layers.

| Principle | Rule |
|-----------|------|
| SRP | Handlers bind HTTP and map errors; services own rules; repositories own SQL. Feature repos embed `pkg/repo.Repository[T]` and add **only** feature queries. |
| OCP / DIP | Services depend on `internal/repository/interfaces` (and small ports like `media.AvatarStorage`), never on concrete `*Repository` types. New OAuth providers implement `oauth.Provider` — no `switch` in `OAuthService`. |
| ISP | Repository interfaces stay narrow. User rows go through `UserRepository`; tokens/OAuth accounts through `AuthRepository`. |
| Reuse | Prefer an existing helper over a copy. Domain-agnostic code belongs in `pkg/`; Gin/JWT/RBAC helpers belong in `common/`. Time/string format helpers belong in `common/utils`. |

**`pkg/` helpers (no `internal` imports)**

| Package | Use for |
|---------|---------|
| `pkg/repo` | Generic Find/Insert/Update/Delete/Paginate over pgx + squirrel |
| `pkg/query` | Filter DSL, pagination clamp, `Page[T]` list JSON |
| `pkg/postgres` | Pool, `WithTx`, `QueryStrings` |
| `pkg/hash` | Password bcrypt, `SHA256Hex` |
| `pkg/cache`, `pkg/redis` | Cache stores |

**`common/httpx`** — `BindJSON` / `BindQuery` / `RequireUserID` / `OK` / `Created` / `NoContent` / `FormFile`. Do not duplicate bind-or-400 in handlers.

**`common/utils`** — `FormatRFC3339` / `FormatRFC3339Ptr` / `FormatDate` / `StringPtr`. Do not copy these helpers in services.

When adding a **core** feature: DTO in `internal/dto/<feature>/`, model in `internal/models/<feature>/`, interface + thin repo, service constructor taking interfaces, handler using `httpx`. When adding an **optional/business** feature: make it a module (`internal/modules/<name>/`, see [`docs/modules.md`](../docs/modules.md)). Modules may import core packages; core must never import a module (only `internal/modules/registry.go` does).

## 6) Import rules

- Use the module path for all in-repo imports: `be/pkg/...`, `be/common/...`, `be/internal/...`, `be/public/...`
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
- App info: `GET /api/app/info` → `{"modules": [...]}` (public; FE hides UI of modules the BE does not run)
- Modules mount on `module.Routes{API, Admin}` (e.g. webhooks: `/api/webhooks/*`)

### 7.1) New feature permissions (mandatory)

When adding a protected admin/API feature:

1. Core resource: add `{resource}:view` / `{resource}:modify` to `DefaultPermissions()` in `internal/database/seeders/catalog.go`. Module resource: return them from the module's `Permissions()`. Use new stable UUIDs; duplicate keys/IDs abort seeding.
2. Protect routes with `middleware.RequireView("resource")` / `RequireModify("resource")`.
3. Admin role receives new keys automatically via `RolePermissionSeeder` (iterates catalog).
4. Align with `docs/features/.../contracts/permissions.md`.

Do not ship admin routes without catalog entries. See `.claude/rules/feature-permissions.md`.

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
| **Docker stack** | `be/.env` via compose `env_file` (root `.env` only holds nginx port / image tag) |
| **Local host / own server** | `be/.env` only — copy from `be/.env.example` |

## 9) Before finishing any task

1. Is the file placed in the correct layer folder?
2. Are imports scoped to the `be` module only?
3. Are DTOs separate from models?
4. Are errors handled via `common/errors` and `response.HandleError`?
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

## 12) Queue / NATS (JetStream) — optional

Async work uses **NATS JetStream**, off by default (`NATS_ENABLED=false`; compose profile `queue`). The HTTP API only publishes; consumers run in **`cmd/queue`**.

```text
internal/queue/              # Infrastructure — client, Stream/Consumer topology, nats.json overrides
internal/handlers/publisher/ # Outbound — PublishJSON(ctx, subject, payload); no-op when disabled
internal/modules/<name>/     # A module declares Streams() (+ consumer Handlers) and/or a Poller
cmd/queue/                   # Worker — EnsureInfrastructure + StartConsumers for module streams,
                             # or polls module Pollers when NATS is disabled, or idles
```

| Rule | Detail |
|------|--------|
| Names | Stream/subject/consumer names are constants in the module that owns them — no string literals elsewhere |
| JSON config | `internal/queue/nats.json` holds the connection and optional per-stream option **overrides** keyed by `Stream.Key`; defaults live in the module |
| Env | `NATS_ENABLED`, `NATS_URL`, optional `QUEUE_CONFIG_FILE` |
| HTTP vs worker | `be/main.go` (API) wires the publisher only; **never** call `StartConsumers` in the API process |
| Services | After the DB write, publish through the kernel `Publisher` — services never import consumer code |

**Adding a consumer (inside a module):** constants → `Streams()` with a `queue.Consumer{Handler: …}` → publish with `Publisher.PublishJSON` → optional `Poll` fallback for NATS-disabled setups.

```bash
go run ./cmd/queue      # queue worker
```
