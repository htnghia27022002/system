# Agent prompts

**Call `@agent` + describe what you want** (Vietnamese OK in chat). **All files written: English only.** No `/speckit-*` slash commands needed.

## Phase order

```text
@ba → spec.md → @technical-architect → tasks.md → plan.md → @be @fe → code + *-tasks-verify.md → @qa → test
```

---

## Phase 1 — Analyze

```text
@ba Auth: email/password register & login, Google OAuth, JWT access + refresh, modern friendly login/register UI
```

```text
@ba Clarify open questions in the auth spec
```

---

## Phase 2 — Decompose

```text
@technical-architect Break down tasks from the auth spec
```

---

## Phase 3 — Design

```text
@technical-architect Design full stack for docs/features/002-auth/ → plan.md + contracts/
```

Creates `plan.md`, `contracts/database.md`, `contracts/endpoints.md`, and `contracts/permissions.md` (RBAC keys + admin menu).

---

## Phase 4 — Implement + verify

```text
@be Implement backend auth per tasks.md, then write be-tasks-verify.md
@fe Implement frontend auth per tasks.md, then write fe-tasks-verify.md
```

---

## Phase 5 — Test

```text
@qa Write qa-checklist and run make test for auth feature
```

---

## Agent → Speckit skill map

| Agent | Skills (auto) |
|-------|----------------|
| `@ba` | `speckit-specify`, `speckit-clarify` |
| `@technical-architect` | `speckit-tasks`, `speckit-plan`, `speckit-analyze` |
| `@be` | `speckit-implement` → `make test-be` → `be-tasks-verify.md` |
| `@fe` | `speckit-implement` → `make test-fe` → `fe-tasks-verify.md` |
| `@qa` | `speckit-checklist`, `speckit-analyze`, `speckit-converge` |

Skill files: `.cursor/skills/speckit-*/SKILL.md`
