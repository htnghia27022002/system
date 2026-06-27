---
name: be
description: BE agent — Go + Gin in be/ and Speckit BE docs. Auto-runs speckit-plan / speckit-implement for backend. User only needs @be in the prompt.
---

# BE Agent

## Invocation (user-facing)

**Only call `@be` + describe what you want.** Do not type slash commands.

Examples:

```text
@be Design auth backend from docs/features/002-auth/spec.md and tasks.md → be-implement.md
@be Implement [BE] tasks for docs/features/002-auth/
```

This agent **reads and executes** the matching Speckit skill automatically.

## Speckit skills (automatic)

| User intent | Read & follow skill | Output |
|-------------|---------------------|--------|
| Technical design (phase 3) | `speckit-plan` | `plan.md` (BE), `be-implement.md` |
| Implement backend (phase 4) | `speckit-implement` | code in `be/`, update `be-implement.md` |
| Review / gap analysis | `speckit-analyze` | report only |

Also read: `docs-feature`, `be-develop`.

**Before any Speckit skill:** read `.cursor/skills/<skill>/SKILL.md` and follow it completely.

**Prerequisites:** `spec.md` exists; for plan → `tasks.md` too; for implement → `plan.md` + `be-implement.md`.

## Read before editing

1. [`be/AGENTS.md`](../../be/AGENTS.md)
2. [`be/README.md`](../../be/README.md)
3. `.cursor/skills/be-develop/SKILL.md`
4. `.cursor/skills/docs-feature/SKILL.md`
5. Active feature under `docs/features/<id>/`

## Module boundary

`be/` is self-contained (`go.mod`). Imports: `be/internal/...` only. FE integration via HTTP `/api` only.

## Architecture

```text
route → handler → service → repository → database
```

| Layer | Path |
|-------|------|
| Routes | `be/public/routes/` |
| Handlers | `be/public/handlers/` |
| Services | `be/internal/services/<feature>/` |
| DI | `be/internal/app/container.go` |

## Feature docs (BE-owned)

| File | Phase |
|------|-------|
| `be-implement.md` | 3 — design |
| `plan.md` | 3 — BE sections |

Execute only `[BE]` tasks from `tasks.md` during `speckit-implement`.

## Run & test (repo root)

```bash
make up-d && make test-be
```

See [`.cursor/rules/environment.mdc`](../../.cursor/rules/environment.mdc).

## Language (mandatory)

All docs and code output in **English only**, even if the user prompts in Vietnamese. See `english-only-file-edits.mdc`.

## Quality baseline

- camelCase JSON, errors via `internal/common/errors`
- GitNexus impact before shared symbol edits (`npx gitnexus query "oauth"`, `npx gitnexus impact OAuthService`)
- `make test-be` after substantive changes

## OAuth providers (adapter pattern)

When touching OAuth or adding a provider, read [`be/AGENTS.md`](../../be/AGENTS.md) §10.

**Before editing:** run GitNexus to locate shared symbols and callers:

```bash
npx gitnexus query "oauth provider callback"
npx gitnexus impact OAuthService
```

**Layout:** one file per provider under `be/internal/services/auth/oauth/` implementing `Provider`; register in `registry.go`. Shared user-linking logic stays in `oauth_service.go` — do not duplicate in adapters.

## Tests (`be/test/`)

All unit, integration, and e2e tests go in `be/test/` only — see [`be/test/README.md`](../../be/test/README.md) and `be/AGENTS.md` §11.

```bash
make test-be              # unit
make test-be-integration  # Postgres required
make test-be-e2e
make test-be-all
```

Use `test/testutil/` for shared mocks. Query GitNexus before adding tests: `npx gitnexus query "auth test"`.
