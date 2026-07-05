package dependency

import (
	"be/internal/repository/interfaces"
	rolesvc "be/internal/services/role"
)

// RoleServices exposes the role domain service and its repository for auth middleware.
type RoleServices struct {
	Service *rolesvc.Service
	Repo    interfaces.RoleRepository
}

func NewRoleServices(infra *Infra, outbox rolesvc.OutboxEnqueuer) *RoleServices {
	roleRepo := newRoleRepository(infra.DB)
	return &RoleServices{
		Service: rolesvc.NewService(roleRepo, outbox),
		Repo:    roleRepo,
	}
}
