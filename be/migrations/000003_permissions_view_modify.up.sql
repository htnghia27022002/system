-- Migrate RBAC catalog to view/modify permission model.
-- Safe to run on dev DBs; re-assign role permissions in admin UI if needed.

INSERT INTO permissions (id, key, name, "group", description) VALUES
  ('10000001-0000-4000-8000-000000000011', 'dashboard:view', 'View dashboard', 'dashboard', 'View the admin dashboard overview'),
  ('10000001-0000-4000-8000-000000000012', 'users:view', 'View users', 'users', 'List and view user accounts'),
  ('10000001-0000-4000-8000-000000000013', 'users:modify', 'Modify users', 'users', 'Create, update, and delete user accounts'),
  ('10000001-0000-4000-8000-000000000014', 'roles:view', 'View roles', 'roles', 'List and view roles'),
  ('10000001-0000-4000-8000-000000000015', 'roles:modify', 'Modify roles', 'roles', 'Create, update, and delete roles'),
  ('10000001-0000-4000-8000-000000000016', 'permissions:view', 'View permissions', 'permissions', 'View the permissions catalog')
ON CONFLICT (key) DO UPDATE SET
  name = EXCLUDED.name,
  "group" = EXCLUDED."group",
  description = EXCLUDED.description;

DELETE FROM role_permissions
WHERE permission_id IN (
  SELECT id FROM permissions
  WHERE key LIKE '%:read'
     OR key LIKE '%:create'
     OR key LIKE '%:update'
     OR key LIKE '%:delete'
);

DELETE FROM permissions
WHERE key LIKE '%:read'
   OR key LIKE '%:create'
   OR key LIKE '%:update'
   OR key LIKE '%:delete';
