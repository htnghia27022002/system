# System Monorepo

pnpm workspace with a **Next.js frontend** (`fe/`), **Go backend** (`be/`), and **Docker** stack for local development.

## Architecture

```text
Browser → nginx (:8080) → /api/* → be (Go + Gin)
                       → /*     → fe (Next.js 15)
be → postgres, redis
```

| Layer | Path | Stack |
|-------|------|-------|
| Frontend | [`fe/`](fe/) | Next.js 15, React 19, TypeScript, Tailwind v4 |
| Backend | [`be/`](be/) | Go 1.22, Gin, GORM, PostgreSQL, JWT, RBAC |
| Infrastructure | [`docker/`](docker/) | nginx, postgres, redis, compose profiles |

**Request flow (BE):** `route → handler → service → repository → database`

## Quick start (Docker — recommended)

```bash
cp .env.docker.example .env
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

# Frontend (separate terminal, repo root)
pnpm install
cp fe/.env.example fe/.env
pnpm dev                        # http://localhost:3000
```

Set `NEXT_PUBLIC_USE_MOCK_API=false` and `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api` in `fe/.env` to use the real API.

## Documentation map

| Doc | Purpose |
|-----|---------|
| [`docs/README.md`](docs/README.md) | Feature docs layout, role ownership |
| [`docs/workflow/overview.md`](docs/workflow/overview.md) | Speckit SDD cycle |
| [`docs/workflow/agent-prompts.md`](docs/workflow/agent-prompts.md) | Prompts per agent (BA/BE/FE/QA) |
| [`AGENTS.md`](AGENTS.md) | Root agent rules, GitNexus, monorepo policy |
| [`fe/AGENTS.md`](fe/AGENTS.md) | FE architecture and import rules (authoritative) |
| [`fe/README.md`](fe/README.md) | FE structure, routes, commands |
| [`be/README.md`](be/README.md) | BE API, layers, env, conventions |
| [`docker/README.md`](docker/README.md) | Docker compose, nginx routing, profiles |
| [`.env.docker.example`](.env.docker.example) | Shared Docker env template |

## Root commands

| Command | Description |
|---------|-------------|
| `make up-d` | Docker dev stack (background) |
| `make down` | Stop containers |
| `make test` | BE (Docker) + FE tests |
| `make test-be` | Go tests inside `be` container |
| `make test-fe` | Vitest (`pnpm test`) |
| `make prod` | Docker production profile |
| `pnpm dev` | FE only (host) |
| `pnpm docker:up` | Same as `make up` |

## Workspace layout

```text
system/
├── docs/               # Feature specs & role docs (Speckit)
├── fe/                 # Next.js app
├── be/                 # Go API
├── docker/             # Dockerfiles + nginx configs
├── docker-compose.yml
├── pnpm-workspace.yaml
├── Makefile
└── AGENTS.md
```
