package dependency

import (
	"be/internal/repository"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
)

func newAuthRepository(db *postgres.Postgres) interfaces.AuthRepository {
	return repository.NewAuthRepository(db)
}

func newUserRepository(db *postgres.Postgres) interfaces.UserRepository {
	return repository.NewUserRepository(db)
}

// NewUserRepository exposes the user repository for auth middleware wiring.
func NewUserRepository(infra *Infra) interfaces.UserRepository {
	return newUserRepository(infra.DB)
}

func newRoleRepository(db *postgres.Postgres) interfaces.RoleRepository {
	return repository.NewRoleRepository(db)
}

func newPermissionRepository(db *postgres.Postgres) interfaces.PermissionRepository {
	return repository.NewPermissionRepository(db)
}
