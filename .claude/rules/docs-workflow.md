---
description: "Feature workflow order spec → tasks → plan (+ contracts) → implement → verify → test."
paths:
  - "docs/**/*"
  - ".specify/**/*"
---

# Docs + Speckit workflow

## Agent auto-skills

| Agent | Skills invoked automatically |
|-------|------------------------------|
| `@agent-ba` | `speckit-specify`, `speckit-clarify` |
| `@agent-technical-architect` | `speckit-tasks`, `speckit-plan`, `speckit-analyze` |
| `@agent-be` | `speckit-implement` → `make test-be` → `be-tasks-verify.md` |
| `@agent-fe` | `speckit-implement` → `make test-fe` → `fe-tasks-verify.md` |
| `@agent-qa` | `speckit-checklist`, `speckit-analyze` |

## Language

Chat in any language; **docs/features/** output **English only** (`english-only-file-edits.md`).

## Phase order (do not skip)

1. **spec.md** — `@agent-ba` `/speckit-specify`
2. **tasks.md** — `@agent-technical-architect` `/speckit-tasks`
3. **plan.md** + **contracts/** — `@agent-technical-architect` `/speckit-plan`
4. **implement + verify** — `@agent-be` `@agent-fe` `/speckit-implement` → `be-tasks-verify.md` / `fe-tasks-verify.md`
5. **test** — `@agent-qa` `/speckit-checklist`, `make test`

## Feature folder

```text
spec.md → tasks.md → plan.md + contracts/{database,endpoints,permissions}.md → (code) → *-tasks-verify.md → qa-checklist.md
```

No role subfolders. Nested **`contracts/`** (database, endpoints, permissions). See `feature-contracts.md` and `feature-permissions.md`.

Prompts: [`docs/workflow/agent-prompts.md`](../../docs/workflow/agent-prompts.md)
