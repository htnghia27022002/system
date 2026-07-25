-- Webhooks: read/unread for owner inbox UX (005-webhooks-tool)
ALTER TABLE webhook_requests
    ADD COLUMN IF NOT EXISTS is_read BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_webhook_requests_inbox_active_read_created
    ON webhook_requests(inbox_id, soft_deleted_at, is_read, created_at DESC);
