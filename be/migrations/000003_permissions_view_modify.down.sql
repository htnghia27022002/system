-- Data migration rollback is best-effort; legacy CRUD keys are not restored automatically.

DELETE FROM permissions
WHERE key IN (
  'dashboard:view',
  'users:view',
  'users:modify',
  'roles:view',
  'roles:modify',
  'permissions:view'
);
