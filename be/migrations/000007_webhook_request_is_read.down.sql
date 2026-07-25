DROP INDEX IF EXISTS idx_webhook_requests_inbox_active_read_created;
ALTER TABLE webhook_requests DROP COLUMN IF EXISTS is_read;
