package dependency

import (
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
}

func NewHTTPHandlers(
	auth *AuthServices,
	user *usersvc.Service,
	role *RoleServices,
	permission *permissionsvc.Service,
	mediaSvc *media.Service,
) *HTTPHandlers {
	return &HTTPHandlers{
		Auth:       handlers.NewAuthHandler(auth.Auth, auth.OAuth),
		User:       handlers.NewUserHandler(user),
		Role:       handlers.NewRoleHandler(role.Service),
		Permission: handlers.NewPermissionHandler(permission),
		Media:      handlers.NewMediaHandler(mediaSvc),
	}
}
