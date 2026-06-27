DROP INDEX IF EXISTS idx_roles_slug;

ALTER TABLE roles DROP COLUMN IF EXISTS slug;

ALTER TABLE permissions DROP COLUMN IF EXISTS "group";

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'permissions' AND column_name = 'key'
  ) AND NOT EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_schema = 'public' AND table_name = 'permissions' AND column_name = 'code'
  ) THEN
    ALTER TABLE permissions RENAME COLUMN key TO code;
  END IF;
END $$;
