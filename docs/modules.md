# Core and modules

This repository is a **base**: clone it for a new project, keep the core, add the
features that project needs as modules. The core holds no business logic.

| Layer | What it contains | Can be removed? |
|-------|------------------|-----------------|
| **Core** | Auth (password + OAuth, refresh sessions), users, roles, permissions (RBAC view/modify), profile + avatar, admin shell, dashboard, design system, i18n, mock API | No |
| **Modules** | Optional features plugged in through a registry. Shipped: **webhooks** (BE) / **tools** (FE) — a public capture URL + owner inbox | Yes — delete the folder + one registry line |
| **Infrastructure** | Postgres (required); Redis cache, NATS queue + `cmd/queue` worker (optional, off by default) | Optional parts are compose profiles |

`be/` and `fe/` are independent apps: each has its own lockfile/`go.mod`, Dockerfile
and `.env`. They talk only over HTTP. Copy one folder to its own repo or server and
set its `.env` — see [Deploying apps separately](#deploying-apps-separately).

---

## How a module plugs in

### Backend (`be/`)

```text
be/internal/module/            # contract: Module, Kernel, optional capabilities, Set helpers
be/internal/modules/registry.go  # the ONLY list of modules
be/internal/modules/<name>/      # one vertical slice per module
├── module.go                    # New(*module.Kernel) module.Module + capabilities
├── handlers/  services/  repository/{,interfaces/}  models/  dto/
└── migrations/                  # own SQL, own version table schema_migrations_<name>
```

A module implements `module.Module` (`Name`, `RegisterRoutes`) and any optional
capability it needs:

| Capability | Used for |
|------------|----------|
| `PermissionProvider` | `{resource}:view` / `{resource}:modify` rows; the admin role gets them automatically |
| `MigrationProvider` | SQL in its own migration set (never shares a version sequence with core) |
| `SeederProvider` | Seeders run after the core seeders |
| `MiddlewareProvider` | Engine-level middleware that must run before core CORS |
| `StreamProvider` | JetStream streams/consumers consumed by `cmd/queue` |
| `Poller` | Fallback work loop for `cmd/queue` when NATS is disabled |
| `Closer` | Cleanup on shutdown |
| `SettingsProvider` | Runtime settings (`<module>.<name>`, typed, with default) shown on the admin Settings page; `Public: true` exposes one on `/api/app/info` |

The `Kernel` gives a module `Config`, `DB`, `JWT`, `Publisher`, `Mailer`
(transactional email), `Storage` (uploads: local or S3), `Audit` (record
`<resource>.<verb>` entries to the audit trail), `Settings` (read at request
time, not in `New`) and the `Auth` middleware. `Routes.Admin` is `/api/admin` with auth already applied. A module's
`New` must only wire objects (no I/O) — CLIs build modules from a Config-only kernel.

`MODULES` in `be/.env` picks the set: unset = registry defaults, a comma list = exactly
those, empty = core only, an unknown name = startup error. `GET /api/app/info` returns
the enabled list and the public settings.

### Frontend (`fe/`)

```text
fe/src/types/modules.ts            # AppModule contract
fe/src/config/modules.ts           # the ONLY list of modules
fe/src/features/<name>/module.ts   # descriptor: nav, breadcrumbs, permissions, public links, sitemap
fe/src/features/<name>/permissions.ts  # definePermissions('<resource>')
fe/src/app/admin/<name>/…          # routes (Next.js file routing cannot be registered dynamically)
```

Core chrome — admin sidebar, breadcrumbs, public header links, sitemap, mock RBAC
catalog — is composed from `src/config/modules.ts`; core never imports a module. A
descriptor with `backendModule` only shows its admin entries when `/api/app/info`
lists that BE module, so a BE deployed without it never shows dead links.

---

## Add a module

**Fastest path — generate it:** `make module name=invoices [label="Invoices"] [singular=invoice]`
writes a working admin CRUD module on both sides and registers it:

| Side | Generated |
|------|-----------|
| BE | `internal/modules/<name>/` — `module.go` (permissions with stable IDs derived from the name, routes under `/api/admin/<name>`), `apidoc.go`, migration, model, DTOs, repository, service (records `<singular>.create/update/delete` to the audit log), handler; registry entry; `test/integration/<name>_module_integration_test.go` (CRUD + permission check) |
| FE | `src/features/<name>/` — `module.ts`, `permissions.ts`, typed API + in-memory mock, TanStack Query hooks, table + columns + form dialog (search, server paging, edit/delete gated by `<name>:modify`), mock test; `app/admin/<name>/page.tsx` behind `PermissionGuard`; registry line; EN/VI strings |

The resource starts with `name` + `description`; reshape the migration, model, DTOs,
FE types and form for real fields, then `make openapi` and run the tests. CI runs the
generator on a throwaway module (`scaffold` job) so the templates cannot rot.

**By hand:**

1. Write the spec/plan/contracts as usual (`docs/features/NNN-name/`).
2. **BE** — create `internal/modules/<name>/` (copy `webhooks/` as the template):
   `module.go` with `Name`, `New`, `RegisterRoutes`, `Permissions` (new stable UUIDs),
   `Migrations` (start at `000001_<name>.up.sql`), `APIDoc` (every route — a test
   enforces it). Append one `Entry` to `internal/modules/registry.go`. Tests go under
   `be/test/unit/<name>/`.
3. **FE** — create `src/features/<name>/` with `permissions.ts`, `module.ts`
   (`backendModule: '<be name>'`), components and hooks; add its `app/` routes wrapped in
   `PermissionGuard`; append the descriptor to `src/config/modules.ts`; add i18n keys.
4. Run `make test` (or `cd be && make test` + `cd fe && pnpm test:run`).

## Remove a module (e.g. in a cloned project)

| Side | Delete | Edit |
|------|--------|------|
| BE | `be/internal/modules/<name>/`, `be/test/unit/<name>/`, its mocks in `be/test/testutil/` | its `Entry` in `be/internal/modules/registry.go` |
| FE | `fe/src/features/<name>/`, its `fe/src/app/**` routes | its line in `fe/src/config/modules.ts`; its i18n keys |

For **tools/webhooks** also: the landing "Tools" section (`fe/src/features/landing/components/sections/tools-section.tsx`
and its use in `LandingPage.tsx`), the `/tools/webhooks/{uuid}` rewrite in
`fe/next.config.ts` and `docker/nginx/config/webhooks-capture.conf` (+ its `include`).

Existing databases keep the module's tables and its `schema_migrations_<name>` row;
drop them with the module's down migration first (`go run ./cmd/migrate -module <name> down`)
if you need a clean schema.

---

## Optional infrastructure

| Need | Turn on |
|------|---------|
| Redis cache | `make up-d PROFILES="dev cache"`, set `REDIS_URL` (+ `CACHE_ENABLED=true`, `CACHE_DRIVER=redis`) in `be/.env` |
| Async jobs (NATS JetStream) | `make up-d PROFILES="dev queue"`, set `NATS_ENABLED=true` in `be/.env`; a module declares streams via `StreamProvider` |
| Prod worker | `make prod PROD_PROFILES="prod prod-queue"` |

With nothing declared, `cmd/queue` starts and idles.

---

## Deploying apps separately

| Deploy | Copy | Configure |
|--------|------|-----------|
| Full stack (Docker) | whole repo | `be/.env`, `fe/.env`, root `.env` (nginx port, image tag) |
| BE only | `be/` | `be/.env` (`DB_*`, `CORS_ORIGINS`, `JWT_SECRET`, `MODULES`, …); build with `docker build be` |
| FE only | `fe/` | `fe/.env` (`NEXT_PUBLIC_API_BASE_URL` → the BE URL); build with `docker build fe` or `pnpm install && pnpm build` |

Migrations and seeders run on BE start (`AUTO_MIGRATE=true`) or as a release step
(`go run ./cmd/migrate up` + `go run ./cmd/seed`).

When FE and BE run on **different sites** (e.g. `app.example.com` and `api.other.com`),
the refresh-token cookie must be cross-site: set `AUTH_COOKIE_SAMESITE=none` (forces
`Secure`, so HTTPS is required) and list the FE origin in `CORS_ORIGINS`. On the same
site (`app.example.com` + `api.example.com`, or one host behind nginx) the default
`lax` works. Put a reverse proxy's address in `TRUSTED_PROXIES` so client IPs and
HTTPS detection are correct.
