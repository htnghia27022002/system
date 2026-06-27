# Project documentation

Feature documentation for this monorepo. **Code lives in `be/` and `fe/`** — `docs/` only holds specs and implementation notes.

## Layout

```text
docs/features/NNN-short-name/
├── README.md
├── spec.md           # 1. BA — feature analysis (requirements)
├── tasks.md          # 2. Work breakdown from spec
├── plan.md           # 3. Technical design (last doc before code)
├── be-implement.md   # 3. BE design detail
├── fe-implement.md   # 3. FE design detail
└── qa-checklist.md   # 5. Test cases & sign-off (after implement)
```

Feature root: `.specify/init-options.json` → `"feature_root": "docs/features"`.

## Workflow order (mandatory)

| Phase | Output | Agent | Speckit skill (auto) |
|-------|--------|-------|----------------------|
| **1. Analyze** | `spec.md` | `@ba` | `speckit-specify`, `speckit-clarify` |
| **2. Decompose** | `tasks.md` | `@ba` | `speckit-tasks` |
| **3. Design** | `plan.md`, `*-implement.md` | `@be` `@fe` | `speckit-plan` |
| **4. Implement** | `be/`, `fe/` | `@be` `@fe` | `speckit-implement` |
| **5. Test** | `qa-checklist.md` | `@qa` | `speckit-checklist`, `make test` |

User prompt: **`@agent` + description** (any language in chat). **All generated files: English only.**

**Do not skip phases.** No plan or code before `spec.md`. No implement before `tasks.md` + `plan.md`.

## Quick start

```text
@ba Auth: email/password, Google login, JWT + refresh, modern login/register UI
@ba Break down tasks from the auth spec
@be Design backend → be-implement.md
@fe Design UI → fe-implement.md
@be Implement backend
@fe Implement frontend
@qa Test and sign-off
```

Details: [`workflow/overview.md`](workflow/overview.md) · Prompts: [`workflow/agent-prompts.md`](workflow/agent-prompts.md).

## Language

All feature docs under `docs/features/` are written in **English**, regardless of the language used in `@agent` prompts.
