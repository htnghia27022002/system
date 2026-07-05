package dependency

import (
	permissionsvc "be/internal/services/permission"
)

func NewPermissionService(infra *Infra) *permissionsvc.Service {
	permissionRepo := newPermissionRepository(infra.DB)
	return permissionsvc.NewService(permissionRepo)
}
