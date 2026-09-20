package dependency

import (
	addresssvc "be/internal/services/address"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	usersvc "be/internal/services/user"
	webhooksvc "be/internal/services/webhook"
	"be/public/handlers"
)

// HTTPHandlers groups HTTP handler constructors wired from resolved services.
type HTTPHandlers struct {
	Auth       *handlers.AuthHandler
	User       *handlers.UserHandler
	Role       *handlers.RoleHandler
	Permission *handlers.PermissionHandler
	Search     *handlers.SearchHandler
	Media      *handlers.MediaHandler
	Webhook    *handlers.WebhookHandler
	Maps       *handlers.MapsHandler
	Address    *handlers.AddressHandler
}

func NewHTTPHandlers(
	auth *AuthServices,
	user *usersvc.Service,
	role *RoleServices,
	permission *permissionsvc.Service,
	search *SearchStack,
	mediaSvc *media.Service,
	webhook *webhooksvc.Service,
	maps *MapsServices,
	address *addresssvc.Service,
) *HTTPHandlers {
	return &HTTPHandlers{
		Auth:       handlers.NewAuthHandler(auth.Auth, auth.OAuth),
		User:       handlers.NewUserHandler(user),
		Role:       handlers.NewRoleHandler(role.Service),
		Permission: handlers.NewPermissionHandler(permission),
		Search:     handlers.NewSearchHandler(search.Service, search.Processor, search.Outbox),
		Media:      handlers.NewMediaHandler(mediaSvc),
		Webhook:    handlers.NewWebhookHandler(webhook),
		Maps:       handlers.NewMapsHandler(maps.Places, maps.Locations, maps.Sources, maps.Ingest, maps.Search),
		Address:    handlers.NewAddressHandler(address),
	}
}
