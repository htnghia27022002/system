# System — core base

A base to clone for new projects: a **Next.js frontend** (`fe/`), a **Go backend** (`be/`)
and a **Docker** stack. It ships only the core (auth, users, roles/permissions, profile,
admin shell, design system) plus optional **modules** — currently the Webhooks tool.
No business logic lives in the core. How modules work, and how to add or remove one:
[`docs/modules.md`](docs/modules.md).

`be/` and `fe/` are independent apps (own `go.mod` / `pnpm-lock.yaml`, `Dockerfile`,
`.env`) that talk over HTTP only — run, build or deploy either one on its own.

## Architecture

```text
Browser → nginx (:8080) → /api/* → be (Go + Gin)
                       → /*     → fe (Next.js 15)
be → postgres                (required)
be → redis, nats             (optional: compose profiles "cache" / "queue")
queue (cmd/queue) → nats     (consumes streams declared by modules)
```

| Layer | Path | Stack |
|-------|------|-------|
| Frontend | [`fe/`](fe/) | Next.js 15, React 19, TypeScript, Tailwind v4 |
| Backend | [`be/`](be/) | Go 1.22, Gin, pgx, squirrel, PostgreSQL, JWT, RBAC |
| Infrastructure | [`docker/`](docker/) | nginx, postgres; optional redis, nats (compose profiles) |

**Request flow (BE):** `route → handler → service → repository → database`

## Quick start (Docker — recommended)

```bash
cp .env.example .env
make up
```

Open **http://system.local** or **http://localhost** (after `make up`).

Add to `/etc/hosts` once: `127.0.0.1 system.local`

API at **http://system.local/api**.

Seed admin: `admin@example.com` / `admin1234`

See [`docker/README.md`](docker/README.md) for commands, profiles, and env vars.

## Local dev (without Docker)

Requires PostgreSQL locally.

```bash
# Backend
cd be && cp .env.example .env   # or use be/.env
make run                        # http://localhost:8080/api

# Frontend (separate terminal)
cd fe && pnpm install
cp .env.example .env
pnpm dev                        # http://localhost:3000
```

Set `NEXT_PUBLIC_USE_MOCK_API=false` and `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api` in `fe/.env` to use the real API.

## Documentation map

| Doc | Purpose |
|-----|---------|
| [`docs/README.md`](docs/README.md) | Feature docs layout, role ownership |
| [`docs/workflow/overview.md`](docs/workflow/overview.md) | Speckit SDD cycle |
| [`docs/workflow/agent-prompts.md`](docs/workflow/agent-prompts.md) | Prompts per agent (BA/BE/FE/QA) |
| [`docs/new-project.md`](docs/new-project.md) | **Start here when cloning:** rename, choose optional pieces, env, first deploy |
| [`docs/modules.md`](docs/modules.md) | Core vs modules, add/remove a module (generator), optional infra, separate deploys |
| [`be/api/openapi.json`](be/api/openapi.json) | API contract (generated; also served at `/api/openapi.json`) |
| [`AGENTS.md`](AGENTS.md) | Root agent rules, GitNexus, monorepo policy |
| [`fe/CLAUDE.md`](fe/CLAUDE.md) | FE architecture and import rules (authoritative) |
| [`fe/README.md`](fe/README.md) | FE structure, routes, commands |
| [`be/README.md`](be/README.md) | BE API, layers, env, conventions |
| [`docker/README.md`](docker/README.md) | Docker compose, nginx routing, profiles |
| [`.env.example`](.env.example) | Docker stack env template (repo root) |
| [`be/.env.example`](be/.env.example) | Local BE only |
| [`fe/.env.example`](fe/.env.example) | Local FE only |

## Root commands

| Command | Description |
|---------|-------------|
| `make up-d` | Docker dev stack (background); add-ons: `PROFILES="dev queue cache"` |
| `make down` | Stop containers |
| `make test` | BE (Docker) + FE tests |
| `make test-be` | Go tests inside `be` container |
| `make test-fe` | Vitest (`pnpm -C fe test:run`) |
| `make prod` | Docker production profile |
| `make bootstrap name=… slug=…` | Rename a fresh clone into a project ([docs/new-project.md](docs/new-project.md)) |
| `make module name=…` | Scaffold a CRUD module on BE + FE |
| `make openapi` | Regenerate `be/api/openapi.json` and `fe/src/types/api.gen.ts` |
| `pnpm setup` / `pnpm dev` | Install / run FE only (host; delegates to `fe/`) |
| `pnpm docker:up` | Same as `make up` |

## CI

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) runs on every PR: BE (gofmt on changed files, vet, unit with `-race`, integration + e2e against a Postgres service), FE (API types in sync with the spec, lint, typecheck, Vitest, `next build`), a scaffold job that generates a throwaway module and builds both apps with it, and a Docker build of each app from its own folder (proves `be/` and `fe/` stay standalone).

## Workspace layout

```text
system/
├── docs/               # Feature specs & role docs (Speckit)
├── fe/                 # Next.js app (own pnpm-lock.yaml, Dockerfile, .env)
├── be/                 # Go API (own go.mod, Dockerfile, .env)
├── docker/             # nginx configs + data volumes
├── docker-compose.yml
├── Makefile
└── AGENTS.md
```
