package dependency

import (
	"be/internal/config"
	accountsvc "be/internal/services/account"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	usersvc "be/internal/services/user"
	"be/public/handlers"
)

// HTTPHandlers groups core HTTP handler constructors wired from resolved services.
// Module handlers are wired inside each module (internal/modules/<name>).
type HTTPHandlers struct {
	Auth       *handlers.AuthHandler
	User       *handlers.UserHandler
	Role       *handlers.RoleHandler
	Permission *handlers.PermissionHandler
	Media      *handlers.MediaHandler
	Account    *handlers.AccountHandler
}

func NewHTTPHandlers(
	cfg config.Config,
	auth *AuthServices,
	user *usersvc.Service,
	role *RoleServices,
	permission *permissionsvc.Service,
	mediaSvc *media.Service,
	account *accountsvc.Service,
) *HTTPHandlers {
	return &HTTPHandlers{
		Auth: handlers.NewAuthHandler(auth.Auth, auth.OAuth, handlers.RefreshCookieConfig{
			Enabled:  cfg.Auth.RefreshCookie,
			Domain:   cfg.Auth.CookieDomain,
			SameSite: cfg.Auth.CookieSameSite,
			Secure:   cfg.Auth.CookieSecure,
			MaxAge:   cfg.JWTRefreshTTL,
		}),
		User:       handlers.NewUserHandler(user),
		Role:       handlers.NewRoleHandler(role.Service),
		Permission: handlers.NewPermissionHandler(permission),
		Media:      handlers.NewMediaHandler(mediaSvc.Store()),
		Account:    handlers.NewAccountHandler(account, user),
	}
}
