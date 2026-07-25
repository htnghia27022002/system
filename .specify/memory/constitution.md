# System Monorepo Constitution

Governance for Speckit-driven development in this workspace.

## Core principles

### I. Package independence

- `be/` and `fe/` are standalone packages with their own module boundaries.
- Feature documentation lives in `docs/features/` and must not import or duplicate application source.
- Root workspace (`docker/`, `docs/`, `.cursor/`) coordinates; it does not own business logic.

### II. Role-owned artifacts

- **BA** → `spec.md`
- **Technical Architect** → `tasks.md`, `plan.md`, `contracts/database.md`, `contracts/endpoints.md`, `contracts/permissions.md`
- **BE** → application code in `be/` + `be-tasks-verify.md`
- **FE** → application code in `fe/` + `fe-tasks-verify.md`
- **QA** → `qa-checklist.md`

Feature docs sit in `docs/features/NNN-name/`. No role subfolders (`ba/`, `be/`, …). The only nested design folder is **`contracts/`** (architect-owned during plan).

No role overwrites another role's primary artifact without explicit user approval.

### III. Spec before code (new features)

Mandatory phase order:

1. **Analyze** — `@ba` `/speckit-specify` → `spec.md`
2. **Decompose** — `@technical-architect` `/speckit-tasks` → `tasks.md` (from spec)
3. **Design** — `@technical-architect` `/speckit-plan` → `plan.md` + `contracts/*`
4. **Implement + verify** — `@be` / `@fe` `/speckit-implement` → code + `be-tasks-verify.md` / `fe-tasks-verify.md`
5. **Test** — `@qa` `/speckit-checklist`, `make test` → `qa-checklist.md`

Skipping phases requires documented user approval.

### IV. API contract alignment

- JSON responses: **camelCase** (BE matches FE types)
- API prefix: `/api`
- Auth: JWT + RBAC permissions (`{resource}:view` / `{resource}:modify`)
- Authoritative schema, HTTP, and RBAC contracts live in `contracts/database.md`, `contracts/endpoints.md`, and `contracts/permissions.md`
- New admin features must seed permissions (BE catalog) and gate admin **menu** + pages (FE `hasPermission` / `PermissionGuard`)
- Contract changes must update `contracts/*` (and be noted in `plan.md`) and be reflected in `*-tasks-verify.md` when implemented

### V. English documentation

All files under `docs/` and all Speckit-generated feature artifacts are written in **English**.

- User chat may be Vietnamese; file output is still English.
- Exception: explicit multilingual UI requests → `fe/src/locales/` only.

## Quality gates

| Phase | Gate | Blocker for next phase? |
|-------|------|-------------------------|
| 1 | `spec.md` complete, no open clarifications | Yes → blocks tasks |
| 2 | `tasks.md` with `[BE]`/`[FE]`/`[QA]` items | Yes → blocks plan |
| 3 | `plan.md` + `contracts/{database,endpoints,permissions}.md` complete | Yes → blocks implement |
| 4 | Code matches tasks + plan + contracts; `be-tasks-verify.md` / `fe-tasks-verify.md` written | Yes → blocks sign-off |
| 5 | `qa-checklist.md` + tests pass | Release gate |

## Governance

- This constitution supersedes ad-hoc workflow shortcuts.
- Amendments: update this file + notify in `docs/workflow/overview.md`.
- Agents must read this file during `/speckit-plan` constitution check.

**Version**: 1.5.0 | **Ratified**: 2026-06-27 | **Last Amended**: 2026-07-25
