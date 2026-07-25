ALTER TABLE users
    DROP COLUMN IF EXISTS social_links,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS birthday,
    DROP COLUMN IF EXISTS general;
