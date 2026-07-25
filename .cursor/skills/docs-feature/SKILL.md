---
name: docs-feature
description: Feature docs in docs/features/. Mandatory order spec → tasks → plan → implement → verify → test. Use for all Speckit feature work.
---

# Docs Feature Skill

## Invocation rule

User writes **`@agent` + description** only. Agent MUST read and execute the matching Speckit skill from `.cursor/skills/speckit-*/SKILL.md` — never ask the user to run slash commands.

## Phase order (mandatory)

```text
1. spec.md              ← /speckit-specify (@ba)
2. tasks.md             ← /speckit-tasks (@technical-architect)
3. plan.md              ← /speckit-plan (@technical-architect)
   contracts/database.md
   contracts/endpoints.md
   contracts/permissions.md
4. be/ fe/ code         ← /speckit-implement (@be @fe)
   be-tasks-verify.md   ← @be after verify
   fe-tasks-verify.md   ← @fe after verify
5. qa-checklist.md      ← /speckit-checklist + make test (@qa)
```

Never implement before phases 1–3 are complete (including `contracts/` when the feature has design). Never hand off to `@qa` before `*-tasks-verify.md` exists for the roles that implemented.

## Language (mandatory)

- User prompts may be Vietnamese; **every file under `docs/features/` must be English**.
- Includes: user stories, acceptance criteria, tasks, API notes, contracts, verify reports, QA checklists.
- Multilingual UI copy belongs only in `fe/src/locales/` when explicitly requested.

## Feature layout

```text
docs/features/NNN-name/
├── spec.md
├── tasks.md
├── plan.md
├── contracts/
│   ├── database.md      # tables / columns / indexes / FKs (or N/A)
│   ├── endpoints.md     # HTTP API contracts (or N/A)
│   └── permissions.md   # RBAC view/modify + admin menu (or N/A)
├── be-tasks-verify.md
├── fe-tasks-verify.md
└── qa-checklist.md
```

No role subfolders (`ba/`, `be/`, …). `contracts/` is the only nested design folder — owned by `@technical-architect` during `speckit-plan`.

## Role ownership by phase

| Phase | Agent | Files |
|-------|-------|-------|
| 1 Analyze | `@ba` | `spec.md` |
| 2 Decompose | `@technical-architect` | `tasks.md` |
| 3 Design | `@technical-architect` | `plan.md`, `contracts/database.md`, `contracts/endpoints.md`, `contracts/permissions.md` |
| 4 Implement + verify | `@be`, `@fe` | code in `be/` / `fe/` + `be-tasks-verify.md` / `fe-tasks-verify.md` |
| 5 Test | `@qa` | `qa-checklist.md` |

## Speckit input rules (project override)

| Command | Required inputs |
|---------|-----------------|
| `/speckit-tasks` | `spec.md` |
| `/speckit-plan` | `spec.md`, `tasks.md` → also write `contracts/*` (database, endpoints, permissions) |
| `/speckit-implement` | `tasks.md`, `plan.md`, `contracts/*` |

## Task prefixes

`[BE]` `[FE]` `[QA]` `[BA]` in `tasks.md`

## New feature permissions

Admin/protected features must define RBAC keys and FE menu gates — see [`.cursor/rules/feature-permissions.mdc`](../../.cursor/rules/feature-permissions.mdc).

## Templates

[`docs/templates/`](../../docs/templates/) · contracts: [`docs/templates/contracts/`](../../docs/templates/contracts/)

Rule: [`.cursor/rules/feature-contracts.mdc`](../../.cursor/rules/feature-contracts.mdc)
