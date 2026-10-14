DROP TABLE IF EXISTS user_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
