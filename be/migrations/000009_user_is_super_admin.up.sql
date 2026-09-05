ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE users
SET is_super_admin = TRUE
WHERE deleted_at IS NULL
  AND (
    id = 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa'
    OR email = 'admin@example.com'
  );

CREATE INDEX IF NOT EXISTS idx_users_is_super_admin
    ON users (is_super_admin)
    WHERE is_super_admin = TRUE AND deleted_at IS NULL;
