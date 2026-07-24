# FE tasks verify: [FEATURE NAME]

**Feature:** `docs/features/NNN-short-name/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@fe`  
**Date:** YYYY-MM-DD

## Summary

[1–3 sentences: what was implemented and overall verify result]

## Tasks completed

Mark each `[FE]` task from `tasks.md` that this pass covered.

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| Txxx | | Done / Partial / Blocked | route, test name, or note |

## Verification commands

```bash
make test-fe
pnpm lint
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass / Fail | |
| `pnpm lint` | Pass / Fail / Skipped | |

## Acceptance coverage (FE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| | | Pass / Fail / N/A |

## Gaps / follow-ups

- [ ] None — ready for `@qa`
- [ ] …

## Sign-off (FE)

- [ ] All claimed `[FE]` tasks done or explicitly deferred above
- [ ] Tests / lint listed above passed (or skipped with reason)
- [ ] Matches `plan.md` FE sections (or deviations documented)
