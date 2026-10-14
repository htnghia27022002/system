-- Account lifecycle: email verification and single-use tokens for password
-- reset, email verification and invitations.

ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ;
-- Accounts that existed before verification was introduced are trusted as-is.
UPDATE users SET email_verified_at = created_at WHERE email_verified_at IS NULL;

CREATE TABLE IF NOT EXISTS user_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- password_reset | email_verify | invite
    purpose VARCHAR(32) NOT NULL,
    -- SHA-256 of the raw token; the raw value only ever exists in the email link.
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_tokens_user_purpose_open
    ON user_tokens (user_id, purpose)
    WHERE used_at IS NULL;
