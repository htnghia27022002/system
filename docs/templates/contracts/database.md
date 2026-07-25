# Database contract: [FEATURE NAME]

**Feature**: `NNN-short-name`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Draft | Ready for implement

> Authoritative persistence contract for this feature. Implementers follow this file; `plan.md` links here.

## Scope

- **New tables**: …
- **Altered tables**: … (or none)
- **N/A**: … (only if this feature adds no persistence — keep this file and explain why)

## Tables

### `[table_name]`

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | uuid | NO | | PK |
| … | … | … | … | … |

**Constraints / indexes**

- PK: …
- UNIQUE: …
- FK: `…` → `…` (ON DELETE …)
- INDEX: …

## Relationships

```text
EntityA (1) ---- (*) EntityB
```

## Migration

- Up/down: `be/migrations/NNNNNN_short_name.{up,down}.sql` (or N/A)
- Notes: …

## Retention / soft-delete (if any)

- …
