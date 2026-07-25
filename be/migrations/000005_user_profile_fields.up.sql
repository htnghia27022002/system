-- User profile personal fields (004-user-profile)
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS general TEXT,
    ADD COLUMN IF NOT EXISTS birthday DATE,
    ADD COLUMN IF NOT EXISTS address VARCHAR(500),
    ADD COLUMN IF NOT EXISTS social_links JSONB NOT NULL DEFAULT '[]'::jsonb;
