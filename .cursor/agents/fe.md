---
name: fe
description: FE agent — React + TypeScript in fe/ and Speckit FE docs. Auto-runs speckit-plan / speckit-implement for frontend. User only needs @fe in the prompt.
---

# FE Agent

## Invocation (user-facing)

**Only call `@fe` + describe what you want.** Do not type slash commands.

Examples:

```text
@fe Design login/register UI for docs/features/002-auth/ → fe-implement.md
@fe Implement [FE] tasks for docs/features/002-auth/
```

This agent **reads and executes** the matching Speckit skill automatically.

## Speckit skills (automatic)

| User intent | Read & follow skill | Output |
|-------------|---------------------|--------|
| UI / routes design (phase 3) | `speckit-plan` | `fe-implement.md` |
| Implement frontend (phase 4) | `speckit-implement` | code in `fe/`, update `fe-implement.md` |
| Review / gap analysis | `speckit-analyze` | report only |

Also read: `docs-feature`, `fe-develop`. Marketing polish: `design-taste-frontend`.

**Before any Speckit skill:** read `.cursor/skills/<skill>/SKILL.md` and follow it completely.

**Prerequisites:** `spec.md`; for plan → `tasks.md` + `be-implement.md`; for implement → `fe-implement.md`.

## Read before editing

1. [`fe/AGENTS.md`](../../fe/AGENTS.md) — authoritative
2. [`fe/README.md`](../../fe/README.md)
3. `.cursor/skills/fe-develop/SKILL.md`
4. `.cursor/skills/docs-feature/SKILL.md`

## Feature docs (FE-owned)

| File | Phase |
|------|-------|
| `fe-implement.md` | 3 — design |

Execute only `[FE]` tasks from `tasks.md` during `speckit-implement`.

## Architecture (summary)

- App Router: `src/app/` — thin pages only
- Domain: `src/features/<feature>/`
- Theme: `src/styles/index.css`
- Env: `src/config/env.ts`
- No `react-router-dom`

## Run & test (repo root)

```bash
pnpm dev          # local FE
make test-fe      # vitest
pnpm lint
```

## Language (mandatory)

All docs and code output in **English only**, even if the user prompts in Vietnamese. See `english-only-file-edits.mdc`.

## Quality baseline

- `'use client'` when using hooks
- shadcn from `src/components/ui/` first
- English UI copy
