# Project documentation

Feature documentation for this monorepo. **Code lives in `be/` and `fe/`** — `docs/` only holds specs and implementation notes.

## Layout

```text
docs/features/NNN-short-name/
├── README.md
├── spec.md              # 1. BA — feature analysis (requirements)
├── tasks.md             # 2. Technical architect — work breakdown from spec
├── plan.md              # 3. Technical architect — technical design (last doc before code)
├── contracts/
│   ├── database.md      # 3. Persistence contract (tables / indexes / FKs)
│   ├── endpoints.md     # 3. HTTP API contract (method / path / request / response)
│   └── permissions.md   # 3. RBAC keys + admin menu gate
├── be-tasks-verify.md   # 4. BE — task completion + verify evidence
├── fe-tasks-verify.md   # 4. FE — task completion + verify evidence
└── qa-checklist.md      # 5. Test cases & sign-off (after implement)
```

Feature root: `.specify/init-options.json` → `"feature_root": "docs/features"`.

## Workflow order (mandatory)

| Phase | Output | Agent | Speckit skill (auto) |
|-------|--------|-------|----------------------|
| **1. Analyze** | `spec.md` | `@ba` | `speckit-specify`, `speckit-clarify` |
| **2. Decompose** | `tasks.md` | `@technical-architect` | `speckit-tasks` |
| **3. Design** | `plan.md`, `contracts/*` | `@technical-architect` | `speckit-plan` |
| **4. Implement + verify** | `be/`, `fe/`, `*-tasks-verify.md` | `@be` `@fe` | `speckit-implement` + package tests |
| **5. Test** | `qa-checklist.md` | `@qa` | `speckit-checklist`, `make test` |

User prompt: **`@agent` + description** (any language in chat). **All generated files: English only.**

**Do not skip phases.** No plan or code before `spec.md`. No implement before `tasks.md` + `plan.md` + `contracts/`. No `@qa` before role verify docs.

## Quick start

```text
@ba Auth: email/password, Google login, JWT + refresh, modern login/register UI
@technical-architect Break down tasks from the auth spec
@technical-architect Design full stack → plan.md
@be Implement backend, then be-tasks-verify.md
@fe Implement frontend, then fe-tasks-verify.md
@qa Test and sign-off
```

Details: [`workflow/overview.md`](workflow/overview.md) · Prompts: [`workflow/agent-prompts.md`](workflow/agent-prompts.md).

## Language

All feature docs under `docs/features/` are written in **English**, regardless of the language used in `@agent` prompts.
