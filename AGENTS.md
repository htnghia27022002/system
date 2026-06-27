<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **system** (1992 symbols, 4918 relationships, 160 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/system/context` | Codebase overview, check index freshness |
| `gitnexus://repo/system/clusters` | All functional areas |
| `gitnexus://repo/system/processes` | All execution flows |
| `gitnexus://repo/system/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->

## Monorepo overview

| Package | Path | Role |
|---------|------|------|
| Frontend | `fe/` | Next.js 15 App Router, React 19, TypeScript |
| Backend | `be/` | Go + Gin API, JWT auth, RBAC |
| Docker | `docker/` + `docker-compose.yml` | nginx, postgres, redis, dev/prod profiles |

Read [`README.md`](README.md) for quick start and doc index.

## Workspace boundaries

`be/` and `fe/` are **independent packages**. They do not import each other. Only root-level tooling knows both:

| Scope | Knows full stack? | Can move standalone? |
|-------|-------------------|----------------------|
| `be/` | No — HTTP API only | Yes — own `go.mod`, imports use module path `be/...` |
| `fe/` | No — calls API via env URL | Yes — own `package.json`, imports use `@/*` |
| Root (`docker/`, docs, `.cursor/`, rules) | Yes — wiring + architecture | N/A |

**Import aliases stay inside each package:**

- **BE:** `be/internal/...` = Go module path from `be/go.mod`, not a monorepo path. File-local aliases (e.g. `authmodel "be/internal/models/auth"`) are fine within `be/`.
- **FE:** `@/*` maps to `fe/src/*` only. Never import from `be/`.

When moving `be/` or `fe/` elsewhere, update Docker/CI/env at the destination; internal code and import paths stay the same.

## Run environments

### Docker (default for full stack)

```bash
cp .env.docker.example .env && make up-d
```

- Entry: `http://system.local:8080` (or port from `NGINX_HTTP_PORT` in `.env`)
- API: `http://system.local:8080/api`
- Env file: repo root `.env` only for compose
- **Agents:** run shell from **repo root**; use `make test`, `make test-be`, `make test-fe` — not `cd be && make test` on WSL host (Go may be missing)

See [`.cursor/rules/environment.mdc`](.cursor/rules/environment.mdc) for command matrix.

### Local (split processes)

- BE: `be/.env` → requires **Go on host** → `cd be && make run`
- FE: `fe/.env` → `pnpm dev` from **repo root** → port 3000
- FE must set `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api` when using real BE

Do not mix Docker FE env with host BE URLs without updating `NEXT_PUBLIC_*` to match.

## Agent routing

Invoke **`@agent` + natural language** — each agent auto-runs Speckit skills (see [`docs/workflow/agent-prompts.md`](docs/workflow/agent-prompts.md)).

| Task | Agent | Auto skills |
|------|-------|-------------|
| Feature spec | `@ba` | `speckit-specify`, `speckit-clarify` |
| Task breakdown | `@ba` | `speckit-tasks` |
| BE design / code | `@be` | `speckit-plan`, `speckit-implement` |
| FE design / code | `@fe` | `speckit-plan`, `speckit-implement` |
| Test & sign-off | `@qa` | `speckit-checklist`, `make test` |
| Docker | — | `docker/README.md` |

Skills: `docs-feature`, `fe-develop`, `be-develop`, `design-taste-frontend` (marketing/landing polish only), Speckit (`speckit-*`).

## Speckit + docs workflow

Feature documentation lives in **`docs/features/`** (configured via `.specify/init-options.json`).

**Phase order:** `spec.md` → `tasks.md` → `plan.md` + `*-implement.md` → implement → test

```text
docs/features/NNN-name/
├── spec.md
├── tasks.md
├── plan.md
├── be-implement.md
├── fe-implement.md
└── qa-checklist.md
```

Full workflow: [`docs/workflow/overview.md`](docs/workflow/overview.md)  
Prompt cheat sheet: [`docs/workflow/agent-prompts.md`](docs/workflow/agent-prompts.md)

## Frontend (`fe/`)

- Routing: `src/app/` App Router only. No `react-router-dom`, no `src/pages/`.
- Domain logic: `src/features/<feature>/`
- Env: `NEXT_PUBLIC_*` via `src/config/env.ts`
- Theme: `src/styles/index.css` (`:root` / `.dark` CSS variables)
- Full rules: [`fe/AGENTS.md`](fe/AGENTS.md)

## Backend (`be/`)

- Layers: `public/routes` → `public/handlers` → `internal/services` → `internal/repository` → DB
- DI: `internal/app/container.go`
- API prefix: `/api/auth`, `/api/admin`
- JSON: camelCase (matches FE contract)
- Module imports: `be/internal/...` (from `be/go.mod`, self-contained)
- Full rules: [`be/AGENTS.md`](be/AGENTS.md), [`be/README.md`](be/README.md)

## Cross-cutting rules

- Write new or modified file content in **English** — including `docs/features/**` from `@ba` `@be` `@fe` `@qa`. User prompts in Vietnamese do not change this.
- Do not commit `.env` files (use `.env.docker.example`, `fe/.env.example`).
- API contract: FE calls `/auth/*` and `/admin/*` under `NEXT_PUBLIC_API_BASE_URL`.
