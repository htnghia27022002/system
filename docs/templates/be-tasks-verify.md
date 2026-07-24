# BE tasks verify: [FEATURE NAME]

**Feature:** `docs/features/NNN-short-name/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@be`  
**Date:** YYYY-MM-DD

## Summary

[1–3 sentences: what was implemented and overall verify result]

## Tasks completed

Mark each `[BE]` task from `tasks.md` that this pass covered.

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| Txxx | | Done / Partial / Blocked | PR path, test name, or note |

## Verification commands

```bash
make test-be
# make test-be-integration
# make test-be-e2e
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | Pass / Fail | |

## Acceptance coverage (BE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| | | Pass / Fail / N/A |

## Gaps / follow-ups

- [ ] None — ready for `@qa`
- [ ] …

## Sign-off (BE)

- [ ] All claimed `[BE]` tasks done or explicitly deferred above
- [ ] Tests listed above passed
- [ ] Matches `plan.md` BE sections (or deviations documented)
