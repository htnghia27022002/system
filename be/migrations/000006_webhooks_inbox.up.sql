-- Webhooks tool: one inbox per user + captured inbound requests (005-webhooks-tool)
CREATE TABLE IF NOT EXISTS webhook_inboxes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    public_uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    lifetime_received INTEGER NOT NULL DEFAULT 0,
    active_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inbox_id UUID NOT NULL REFERENCES webhook_inboxes(id) ON DELETE CASCADE,
    method VARCHAR(16) NOT NULL,
    url TEXT NOT NULL,
    client_ip VARCHAR(64) NOT NULL DEFAULT '',
    headers JSONB NOT NULL DEFAULT '{}'::jsonb,
    query JSONB NOT NULL DEFAULT '{}'::jsonb,
    form JSONB NOT NULL DEFAULT '{}'::jsonb,
    body BYTEA,
    content_type VARCHAR(255) NOT NULL DEFAULT '',
    body_truncated BOOLEAN NOT NULL DEFAULT FALSE,
    capture_status VARCHAR(32) NOT NULL DEFAULT 'ok',
    soft_deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_inboxes_public_uuid ON webhook_inboxes(public_uuid);
CREATE INDEX IF NOT EXISTS idx_webhook_requests_inbox_created
    ON webhook_requests(inbox_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_webhook_requests_inbox_active_created
    ON webhook_requests(inbox_id, soft_deleted_at, created_at DESC);
