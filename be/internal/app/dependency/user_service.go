package dependency

import (
	"be/internal/services/media"
	usersvc "be/internal/services/user"
)

func NewUserService(infra *Infra, mediaSvc *media.Service) *usersvc.Service {
	userRepo := newUserRepository(infra.DB)
	authRepo := newAuthRepository(infra.DB)
	return usersvc.NewService(userRepo, authRepo, mediaSvc)
}

func NewMediaService(infra *Infra) *media.Service {
	return media.NewService(infra.Config.UploadDir)
}
