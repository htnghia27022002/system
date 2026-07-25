package dependency

import (
	"be/internal/repository"
	"be/internal/repository/interfaces"

	"gorm.io/gorm"
)

func newAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return repository.NewAuthRepository(db)
}

func newUserRepository(db *gorm.DB) interfaces.UserRepository {
	return repository.NewUserRepository(db)
}

func newRoleRepository(db *gorm.DB) interfaces.RoleRepository {
	return repository.NewRoleRepository(db)
}

func newPermissionRepository(db *gorm.DB) interfaces.PermissionRepository {
	return repository.NewPermissionRepository(db)
}

func newWebhookInboxRepository(db *gorm.DB) interfaces.WebhookInboxRepository {
	return repository.NewWebhookInboxRepository(db)
}

func newWebhookRequestRepository(db *gorm.DB) interfaces.WebhookRequestRepository {
	return repository.NewWebhookRequestRepository(db)
}
