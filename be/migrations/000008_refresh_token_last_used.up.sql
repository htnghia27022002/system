-- Account sessions: last activity on refresh_tokens (006-account-sessions)
-- DEFAULT NOW() fills existing rows and concurrent inserts so SET NOT NULL cannot see NULLs.
ALTER TABLE refresh_tokens
    ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ DEFAULT NOW();

UPDATE refresh_tokens
SET last_used_at = created_at
WHERE created_at IS NOT NULL;

ALTER TABLE refresh_tokens
    ALTER COLUMN last_used_at SET NOT NULL;

ALTER TABLE refresh_tokens
    ALTER COLUMN last_used_at DROP DEFAULT;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_active_last_used
    ON refresh_tokens (user_id, last_used_at DESC)
    WHERE revoked_at IS NULL;
