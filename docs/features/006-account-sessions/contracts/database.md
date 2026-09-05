# Database contract: Account Sessions

**Feature**: `006-account-sessions`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative persistence contract for this feature. Implementers follow this file; `plan.md` links here.

## Scope

- **New tables**: none
- **Altered tables**: `refresh_tokens` (add `last_used_at`)
- **N/A**: not applicable — feature alters existing auth persistence
- **Unused in P1**: `refresh_tokens.device_id` already exists — do **not** populate, index for sessions, or expose it

## Tables

### `refresh_tokens`

Existing table from `be/migrations/000001_init_schema.up.sql`. P1 treats each **active** row as a session.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK; **session id** returned as `sessionId` on auth payloads |
| `user_id` | UUID | NO | | FK → `users(id)` ON DELETE CASCADE; owner scope |
| `token_hash` | VARCHAR(255) | NO | | UNIQUE; SHA-256 of raw refresh token; **never** returned on session list |
| `expires_at` | TIMESTAMPTZ | NO | | Active only while `expires_at > now()` |
| `revoked_at` | TIMESTAMPTZ | YES | | NULL = not ended; set on logout / revoke / refresh predecessor |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | Session **started**; copy from predecessor on refresh rotation |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | Existing column |
| `ip_address` | VARCHAR(64) | YES | | Already exists; populate on create/renew; empty → JSON `null` |
| `user_agent` | TEXT | YES | | Already exists; raw UA in P1; empty → JSON `null` |
| `device_id` | VARCHAR(255) | YES | | Already exists; **do not use in P1** |
| `last_used_at` | TIMESTAMPTZ | NO | | **New**; backfill from `created_at`; set on create; update on renew |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `token_hash`
- FK: `user_id` → `users(id)` ON DELETE CASCADE
- INDEX (existing): none required beyond PK/unique for P1 small lists
- INDEX (optional, recommended): `idx_refresh_tokens_user_active_last_used` on `(user_id, last_used_at DESC)` WHERE `revoked_at IS NULL`

**Active session predicate (application + queries)**

```text
user_id = :jwtUserId
AND revoked_at IS NULL
AND expires_at > now()
```

## Relationships

```text
users (1) ---- (*) refresh_tokens
```

Ended and expired rows remain in the table (audit leftover from 001-auth) but **MUST NOT** appear in the P1 list.

## Migration

- Up/down: `be/migrations/000008_refresh_token_last_used.{up,down}.sql`
- Notes:
  - `ADD COLUMN last_used_at TIMESTAMPTZ`
  - Backfill: `UPDATE refresh_tokens SET last_used_at = created_at WHERE last_used_at IS NULL`
  - Then `SET NOT NULL` if the column was added nullable
  - Down: `DROP COLUMN last_used_at`
  - Do **not** add columns for geo, friendly device name, or fingerprint
  - Next version after `000007_webhook_request_is_read`

## Retention / soft-delete (if any)

- Revoke sets `revoked_at` (existing pattern). No hard-delete of sessions in P1.
- Expired rows stay until some future cleanup; they are omitted from list queries only.

## Repository methods (AuthRepository)

Reuse `be/pkg/repo` on table `refresh_tokens`. Do **not** add these to `UserRepository`.

| Method | Behavior |
|--------|----------|
| `ListActiveByUserID` | Active predicate; return all matching rows (no pagination) |
| `RevokeByID` | Set `revoked_at` where `id` AND `user_id` match an active row; 0 rows → not found |
| `RevokeAllExcept` | Set `revoked_at` on all active rows for `user_id` except `keepID` |

Existing: `CreateRefreshToken`, `FindRefreshTokenByHash`, `RevokeRefreshToken` (hash) remain for issue / refresh / logout.
