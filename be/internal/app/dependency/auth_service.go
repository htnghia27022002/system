package dependency

import (
	"be/internal/repository"
	authsvc "be/internal/services/auth"
	"be/internal/services/media"
)

// AuthServices groups authentication-related services.
type AuthServices struct {
	Auth  *authsvc.Service
	OAuth *authsvc.OAuthService
}

func NewAuthServices(infra *Infra, mediaSvc *media.Service) *AuthServices {
	authRepo := newAuthRepository(infra.DB)
	userRepo := newUserRepository(infra.DB)
	roleRepo := newRoleRepository(infra.DB)

	authService := authsvc.NewService(
		authRepo,
		userRepo,
		roleRepo,
		infra.JWT,
		infra.Config.JWTRefreshTTL,
		mediaSvc,
	)

	authService.SetLoginLockout(
		repository.NewLoginAttemptRepository(infra.DB),
		infra.Config.Auth.MaxFailedLogins,
		infra.Config.Auth.LockoutDuration,
	)

	return &AuthServices{
		Auth:  authService,
		OAuth: authsvc.NewOAuthService(infra.Config, authRepo, userRepo, roleRepo, authService),
	}
}
