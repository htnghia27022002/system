-- Align RBAC fields with frontend contract (idempotent for fresh and legacy databases)

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'permissions' AND column_name = 'code'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'permissions' AND column_name = 'key'
  ) THEN
    ALTER TABLE permissions RENAME COLUMN code TO key;
  END IF;
END $$;

ALTER TABLE permissions ADD COLUMN IF NOT EXISTS "group" VARCHAR(50) NOT NULL DEFAULT '';

ALTER TABLE roles ADD COLUMN IF NOT EXISTS slug VARCHAR(50);

UPDATE roles SET slug = LOWER(REPLACE(name, ' ', '-')) WHERE slug IS NULL OR slug = '';

ALTER TABLE roles ALTER COLUMN slug SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_slug ON roles(slug);
