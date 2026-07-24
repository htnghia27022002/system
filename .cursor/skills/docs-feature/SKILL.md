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
4. be/ fe/ code         ← /speckit-implement (@be @fe)
   be-tasks-verify.md   ← @be after verify
   fe-tasks-verify.md   ← @fe after verify
5. qa-checklist.md      ← /speckit-checklist + make test (@qa)
```

Never implement before phases 1–3 are complete. Never hand off to `@qa` before `*-tasks-verify.md` exists for the roles that implemented.

## Language (mandatory)

- User prompts may be Vietnamese; **every file under `docs/features/` must be English**.
- Includes: user stories, acceptance criteria, tasks, API notes, verify reports, QA checklists.
- Multilingual UI copy belongs only in `fe/src/locales/` when explicitly requested.

## Flat layout

```text
docs/features/NNN-name/
├── spec.md
├── tasks.md
├── plan.md
├── be-tasks-verify.md
├── fe-tasks-verify.md
└── qa-checklist.md
```

## Role ownership by phase

| Phase | Agent | Files |
|-------|-------|-------|
| 1 Analyze | `@ba` | `spec.md` |
| 2 Decompose | `@technical-architect` | `tasks.md` |
| 3 Design | `@technical-architect` | `plan.md` |
| 4 Implement + verify | `@be`, `@fe` | code in `be/` / `fe/` + `be-tasks-verify.md` / `fe-tasks-verify.md` |
| 5 Test | `@qa` | `qa-checklist.md` |

## Speckit input rules (project override)

| Command | Required inputs |
|---------|-----------------|
| `/speckit-tasks` | `spec.md` |
| `/speckit-plan` | `spec.md`, `tasks.md` |
| `/speckit-implement` | `tasks.md`, `plan.md` |

## Task prefixes

`[BE]` `[FE]` `[QA]` `[BA]` in `tasks.md`

## Templates

[`docs/templates/`](../../docs/templates/)
