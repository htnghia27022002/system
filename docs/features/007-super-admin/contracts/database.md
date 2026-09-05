# Database contract: Super Admin Flag

**Feature**: `007-super-admin`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative persistence contract. Implementers follow this file; `plan.md` links here.

## Scope

- **New tables**: none
- **Altered tables**: `users` (add `is_super_admin`)
- **N/A**: not applicable

## Tables

### `users`

Existing table from `be/migrations/000001_init_schema.up.sql`. P1 adds one column.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| *(existing columns unchanged)* | | | | |
| `is_super_admin` | BOOLEAN | NO | `FALSE` | **New**. Super admin bypass flag. Not a role. |

**Constraints / indexes**

- Partial index (optional, recommended): `idx_users_is_super_admin` on `(is_super_admin)` WHERE `is_super_admin = TRUE` AND `deleted_at IS NULL`
- Soft-deleted rows (`deleted_at IS NOT NULL`) do **not** count as Super admins for “last Super admin”

**Backfill**

```sql
UPDATE users
SET is_super_admin = TRUE
WHERE deleted_at IS NULL
  AND (
    id = 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa'
    OR email = 'admin@example.com'
  );
```

Seed insert for new installs must set `is_super_admin = TRUE` for that admin row and `FALSE` for demo member.

## Migration

- File: `be/migrations/000009_user_is_super_admin.up.sql` / `.down.sql`
- Down: drop index (if created), `ALTER TABLE users DROP COLUMN IF EXISTS is_super_admin`

## Out of scope

- New roles / permissions tables
- Audit table for Super admin grants
