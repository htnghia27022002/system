-- Failed-login tracking for temporary account lockout (AUTH_MAX_FAILED_LOGINS,
-- AUTH_LOCKOUT_DURATION). Kept out of `users` so the user model stays unchanged.
CREATE TABLE IF NOT EXISTS login_attempts (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    failed_count INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
