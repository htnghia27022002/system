# Agent prompts

**Call `@agent` + describe what you want** (Vietnamese OK in chat). **All files written: English only.** No `/speckit-*` slash commands needed.

## Phase order

```text
@ba → spec.md → tasks.md → @be @fe → plan + *-implement → @be @fe → code → @qa → test
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
@ba Break down tasks from the auth spec
```

---

## Phase 3 — Design

```text
@be Design auth backend from spec + tasks → be-implement.md
@fe Design login/register UI from spec + be-implement → fe-implement.md
```

---

## Phase 4 — Implement

```text
@be Implement backend auth per tasks.md
@fe Implement frontend auth per tasks.md
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
| `@ba` | `speckit-specify`, `speckit-clarify`, `speckit-tasks` |
| `@be` | `speckit-plan`, `speckit-implement`, `speckit-analyze` |
| `@fe` | `speckit-plan`, `speckit-implement`, `speckit-analyze` |
| `@qa` | `speckit-checklist`, `speckit-analyze`, `speckit-converge` |

Skill files: `.cursor/skills/speckit-*/SKILL.md`
