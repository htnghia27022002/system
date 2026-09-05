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

func newRoleRepository(db *postgres.Postgres) interfaces.RoleRepository {
	return repository.NewRoleRepository(db)
}

func newPermissionRepository(db *postgres.Postgres) interfaces.PermissionRepository {
	return repository.NewPermissionRepository(db)
}

func newWebhookInboxRepository(db *postgres.Postgres) interfaces.WebhookInboxRepository {
	return repository.NewWebhookInboxRepository(db)
}

func newWebhookRequestRepository(db *postgres.Postgres) interfaces.WebhookRequestRepository {
	return repository.NewWebhookRequestRepository(db)
}
