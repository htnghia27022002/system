-- RBAC seed data (reference migration; runtime seed in internal/database/seed.go is preferred for dev)

INSERT INTO permissions (id, key, name, "group", description) VALUES
  ('10000001-0000-4000-8000-000000000001', 'dashboard:read', 'View dashboard', 'dashboard', 'Access the admin dashboard overview'),
  ('10000001-0000-4000-8000-000000000002', 'users:read', 'View users', 'users', 'List and view user accounts'),
  ('10000001-0000-4000-8000-000000000003', 'users:create', 'Create users', 'users', 'Create new user accounts'),
  ('10000001-0000-4000-8000-000000000004', 'users:update', 'Update users', 'users', 'Edit existing user accounts'),
  ('10000001-0000-4000-8000-000000000005', 'users:delete', 'Delete users', 'users', 'Remove user accounts'),
  ('10000001-0000-4000-8000-000000000006', 'roles:read', 'View roles', 'roles', 'List and view roles'),
  ('10000001-0000-4000-8000-000000000007', 'roles:create', 'Create roles', 'roles', 'Create new roles'),
  ('10000001-0000-4000-8000-000000000008', 'roles:update', 'Update roles', 'roles', 'Edit existing roles'),
  ('10000001-0000-4000-8000-000000000009', 'roles:delete', 'Delete roles', 'roles', 'Remove roles'),
  ('10000001-0000-4000-8000-000000000010', 'permissions:read', 'View permissions', 'permissions', 'View the permissions catalog')
ON CONFLICT DO NOTHING;

INSERT INTO roles (id, name, slug, description) VALUES
  ('11111111-1111-4111-8111-111111111111', 'Administrator', 'admin', 'Full system access'),
  ('22222222-2222-4222-8222-222222222222', 'Member', 'user', 'Standard member access')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT '11111111-1111-4111-8111-111111111111', id FROM permissions
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) VALUES
  ('22222222-2222-4222-8222-222222222222', '10000001-0000-4000-8000-000000000001')
ON CONFLICT DO NOTHING;
