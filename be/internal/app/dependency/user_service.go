package dependency

import (
	"be/internal/services/media"
	usersvc "be/internal/services/user"
	"be/pkg/storage"
)

func NewUserService(infra *Infra, mediaSvc *media.Service) *usersvc.Service {
	userRepo := newUserRepository(infra.DB)
	authRepo := newAuthRepository(infra.DB)
	return usersvc.NewService(userRepo, authRepo, mediaSvc)
}

func NewMediaService(store storage.Storage) *media.Service {
	return media.NewService(store)
}
