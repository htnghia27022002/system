package dependency

import (
	permissionsvc "be/internal/services/permission"
	usersvc "be/internal/services/user"
	"be/public/handlers"
)

// HTTPHandlers groups HTTP handler constructors wired from resolved services.
type HTTPHandlers struct {
	Auth       *handlers.AuthHandler
	User       *handlers.UserHandler
	Role       *handlers.RoleHandler
	Permission *handlers.PermissionHandler
	Search     *handlers.SearchHandler
}

func NewHTTPHandlers(
	auth *AuthServices,
	user *usersvc.Service,
	role *RoleServices,
	permission *permissionsvc.Service,
	search *SearchStack,
) *HTTPHandlers {
	return &HTTPHandlers{
		Auth:       handlers.NewAuthHandler(auth.Auth, auth.OAuth),
		User:       handlers.NewUserHandler(user),
		Role:       handlers.NewRoleHandler(role.Service),
		Permission: handlers.NewPermissionHandler(permission),
		Search:     handlers.NewSearchHandler(search.Service, search.Processor, search.Outbox),
	}
}
