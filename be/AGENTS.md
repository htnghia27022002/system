# BE Agent Rules (Go + Gin)

This file defines strict implementation rules for AI agents working in `be/`.
Follow these rules before creating, editing, moving, or deleting files.

## 1) Self-contained module

`be/` is a **standalone Go module**. It can be moved out of this workspace without changing internal import paths.

- Module path: declared in `be/go.mod` (`module be`)
- Internal imports: `be/internal/...`, `be/public/...` — these refer to the **module name**, not the monorepo folder layout
- Run all Go commands from `be/` (`go run .`, `make test`, etc.)
- Do **not** import from `fe/`, repo root, or any path outside this module

When relocating the backend, move the entire `be/` directory (including `go.mod`). Update only external wiring (Docker, CI, env files) at the new location.

## 2) Source of truth

- Use `README.md` in this folder for API routes, layers, and commands
- Use `internal/app/container.go` for dependency wiring

## 3) Stack

Go 1.22 · Gin · GORM · PostgreSQL · JWT · OAuth2 · golangci-lint

## 4) Source layout

```text
be/
├── main.go
├── go.mod              ← module boundary; all internal imports start with `be/`
├── public/
│   ├── api.go
│   ├── handlers/
│   └── routes/
├── internal/
│   ├── app/container.go
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── dto/
│   ├── repository/
│   ├── services/
│   └── common/
└── migrations/
```

## 5) Layer rules

```text
route → handler → service → repository interface → repository → database
```

| Layer | Path |
|-------|------|
| HTTP routes | `public/routes/` |
| HTTP handlers | `public/handlers/` |
| Bootstrap | `public/api.go` |
| DI container | `internal/app/container.go` |
| Business logic | `internal/services/<feature>/` |
| Persistence contracts | `internal/repository/interfaces/` |
| Persistence impl | `internal/repository/` |
| Models / DTOs | `internal/models/`, `internal/dto/` |
| Shared helpers | `internal/common/` |
| Middleware | `internal/middleware/` |

Never put business logic in handlers or HTTP formatting in repositories.

## 6) Import rules

- Use the module path for all in-repo imports: `be/internal/...`, `be/public/...`
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

Integration with the frontend is **HTTP only**. Do not import FE code or share types across repos.

## 8) Configuration (hybrid)

| Layer | Source | Contents |
|-------|--------|----------|
| Public | `be/config.yaml` | Port, JWT TTL, CORS, OAuth provider list, non-secret defaults |
| Secrets | env only | `DB_PASSWORD`, `JWT_SECRET`, `OAUTH_GOOGLE_CLIENT_ID`, `OAUTH_GOOGLE_CLIENT_SECRET` |
| Override | env optional | `DB_HOST`, `REDIS_URL`, `CORS_ORIGINS`, `CONFIG_FILE` |

Priority: **env > config.yaml > built-in default**. Never commit secrets in YAML.

**Docker:** root `docker-compose.yml` — see [`docker/README.md`](../docker/README.md).

**Local:** secrets in `be/.env`; public settings in `config.yaml`.

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
4. Add env fields in `internal/config/config.go` and `be/.env.example`.
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
