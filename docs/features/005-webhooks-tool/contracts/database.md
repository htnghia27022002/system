# Database contract: Tools → Webhooks

**Feature**: `005-webhooks-tool`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement (as-shipped)

> Authoritative persistence contract for this feature. Implementers follow this file; `plan.md` links here.  
> Aligned with shipped migration `be/migrations/000006_webhooks_inbox.up.sql` and models under `be/internal/models/webhook/`.

## Scope

- **New tables**: `webhook_inboxes`, `webhook_requests`
- **Altered tables**: none
- **N/A**: not applicable — feature owns PostgreSQL persistence

## Tables

### `webhook_inboxes`

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `user_id` | UUID | NO | | UNIQUE; FK → `users(id)` ON DELETE CASCADE; one inbox per account |
| `public_uuid` | UUID | NO | `gen_random_uuid()` | UNIQUE; public capture identity; rotated on regenerate |
| `lifetime_received` | INTEGER | NO | `0` | Count of currently stored request rows (decrements on hard-purge) |
| `active_count` | INTEGER | NO | `0` | Non-soft-deleted request count |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `user_id`, `public_uuid`
- FK: `user_id` → `users(id)` ON DELETE CASCADE
- INDEX: `idx_webhook_inboxes_public_uuid` on `(public_uuid)` (capture lookup)

### `webhook_requests`

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `inbox_id` | UUID | NO | | FK → `webhook_inboxes(id)` ON DELETE CASCADE |
| `method` | VARCHAR(16) | NO | | HTTP method as received (GET, POST, …) |
| `url` | TEXT | NO | | Request URI / path + query as received |
| `client_ip` | VARCHAR(64) | NO | `''` | Best-effort; may be empty |
| `headers` | JSONB | NO | `'{}'::jsonb` | Map of header name → string or string[] |
| `query` | JSONB | NO | `'{}'::jsonb` | Query params map |
| `form` | JSONB | NO | `'{}'::jsonb` | Form fields when urlencoded/multipart |
| `body` | BYTEA | YES | | Raw body; null/empty for oversize markers |
| `content_type` | VARCHAR(255) | NO | `''` | |
| `body_truncated` | BOOLEAN | NO | `FALSE` | Truncation / oversize marker flag |
| `capture_status` | VARCHAR(32) | NO | `'ok'` | Values: `ok`, `oversized`, `error` (see model constants) |
| `is_read` | BOOLEAN | NO | `FALSE` | Owner read state; set true when detail is opened |
| `soft_deleted_at` | TIMESTAMPTZ | YES | | NULL = active; set on soft-delete |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | List sort: newest first |

**Constraints / indexes**

- PK: `id`
- FK: `inbox_id` → `webhook_inboxes(id)` ON DELETE CASCADE
- INDEX: `idx_webhook_requests_inbox_created` on `(inbox_id, created_at DESC)`
- INDEX: `idx_webhook_requests_inbox_active_created` on `(inbox_id, soft_deleted_at, created_at DESC)`
- INDEX: `idx_webhook_requests_inbox_active_read_created` on `(inbox_id, soft_deleted_at, is_read, created_at DESC)`

## Relationships

```text
users (1) ---- (1) webhook_inboxes
webhook_inboxes (1) ---- (*) webhook_requests
```

## Migration

- Up/down: `be/migrations/000006_webhooks_inbox.{up,down}.sql`
- Read flag: `be/migrations/000007_webhook_request_is_read.{up,down}.sql`
- Notes:
  - Soft-deleted rows **count toward** the retention cap of **200** stored requests per inbox.
  - Hard-purge deletes oldest rows when count exceeds 200; counters adjusted in application service (`MaxStoredRequests = 200`).
  - Body size limit (**1 MiB**) is enforced in application code, not a DB CHECK.

## Retention / soft-delete

- Soft-delete: set `soft_deleted_at`; decrement `active_count` only; idempotent.
- Hard-purge: delete oldest rows beyond 200; adjust `lifetime_received` / `active_count` per service rules.
- Soft-deleted detail is not returned on GET (404 for owner detail in P1; no trash browser).
