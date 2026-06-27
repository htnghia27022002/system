# System Monorepo Constitution

Governance for Speckit-driven development in this workspace.

## Core principles

### I. Package independence

- `be/` and `fe/` are standalone packages with their own module boundaries.
- Feature documentation lives in `docs/features/` and must not import or duplicate application source.
- Root workspace (`docker/`, `docs/`, `.cursor/`) coordinates; it does not own business logic.

### II. Role-owned artifacts

- **BA** → `spec.md`
- **All roles** → `tasks.md` (decomposed from spec)
- **BE** → `plan.md`, `be-implement.md`
- **FE** → `fe-implement.md`
- **QA** → `qa-checklist.md`

All files sit **flat** in `docs/features/NNN-name/` — no role subfolders.

No role overwrites another role's primary artifact without explicit user approval.

### III. Spec before code (new features)

Mandatory phase order:

1. **Analyze** — `/speckit-specify` → `spec.md`
2. **Decompose** — `/speckit-tasks` → `tasks.md` (from spec)
3. **Design** — `/speckit-plan` → `plan.md`, `be-implement.md`, `fe-implement.md`
4. **Implement** — `/speckit-implement` → `be/`, `fe/`
5. **Test** — `/speckit-checklist`, `make test` → `qa-checklist.md`

Skipping phases requires documented user approval.

### IV. API contract alignment

- JSON responses: **camelCase** (BE matches FE types)
- API prefix: `/api`
- Auth: JWT + RBAC permissions
- Contract changes must update plan and both `be-implement.md` / `fe-implement.md`

### V. English documentation

All files under `docs/` and all Speckit-generated feature artifacts are written in **English**.

- User chat may be Vietnamese; file output is still English.
- Exception: explicit multilingual UI requests → `fe/src/locales/` only.

## Quality gates

| Phase | Gate | Blocker for next phase? |
|-------|------|-------------------------|
| 1 | `spec.md` complete, no open clarifications | Yes → blocks tasks |
| 2 | `tasks.md` with `[BE]`/`[FE]`/`[QA]` items | Yes → blocks plan |
| 3 | `plan.md` + `be-implement.md` + `fe-implement.md` | Yes → blocks implement |
| 4 | Code matches tasks + plan | Yes → blocks sign-off |
| 5 | `qa-checklist.md` + tests pass | Release gate |

## Governance

- This constitution supersedes ad-hoc workflow shortcuts.
- Amendments: update this file + notify in `docs/workflow/overview.md`.
- Agents must read this file during `/speckit-plan` constitution check.

**Version**: 1.1.0 | **Ratified**: 2026-06-27 | **Last Amended**: 2026-06-27
