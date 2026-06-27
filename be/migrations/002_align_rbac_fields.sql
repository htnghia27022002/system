-- Align RBAC fields with frontend contract

ALTER TABLE permissions RENAME COLUMN code TO key;

ALTER TABLE permissions ADD COLUMN IF NOT EXISTS "group" VARCHAR(50) NOT NULL DEFAULT '';

ALTER TABLE roles ADD COLUMN IF NOT EXISTS slug VARCHAR(50);

UPDATE roles SET slug = LOWER(REPLACE(name, ' ', '-')) WHERE slug IS NULL OR slug = '';

ALTER TABLE roles ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_slug ON roles(slug);
