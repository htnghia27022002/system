package app

import (
	jwtmanager "be/internal/common/jwt"
	"be/internal/app/dependency"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/queue"
	"be/internal/repository/interfaces"
	searchpkg "be/internal/search"
	authsvc "be/internal/services/auth"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	rolesvc "be/internal/services/role"
	searchsvc "be/internal/services/search"
	usersvc "be/internal/services/user"
	webhooksvc "be/internal/services/webhook"
	"be/public/handlers"

	"gorm.io/gorm"
)

type Container struct {
	Config            config.Config
	Queue             queue.Config
	QueueClient       *queue.Client
	DB                *gorm.DB
	JWT               *jwtmanager.Manager
	Publisher         *publisher.Publisher
	MediaService      *media.Service
	AuthService       *authsvc.Service
	OAuthService      *authsvc.OAuthService
	UserService       *usersvc.Service
	RoleService       *rolesvc.Service
	PermissionService *permissionsvc.Service
	SearchService     *searchsvc.Service
	SearchProcessor   *searchsvc.IndexProcessor
	OutboxService     *searchsvc.OutboxService
	WebhookService    *webhooksvc.Service
	SearchClient      *searchpkg.Client
	RoleRepo          interfaces.RoleRepository
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleHandler       *handlers.RoleHandler
	PermissionHandler *handlers.PermissionHandler
	SearchHandler     *handlers.SearchHandler
	MediaHandler      *handlers.MediaHandler
	WebhookHandler    *handlers.WebhookHandler
}

func NewContainer(cfg config.Config, db *gorm.DB) *Container {
	infra := dependency.NewInfra(cfg, db)
	searchStack := dependency.NewSearchStack(infra)
	mediaSvc := dependency.NewMediaService(infra)
	authServices := dependency.NewAuthServices(infra, mediaSvc)
	userService := dependency.NewUserService(infra, searchStack.Outbox, mediaSvc)
	roleServices := dependency.NewRoleServices(infra, searchStack.Outbox)
	permissionService := dependency.NewPermissionService(infra)
	webhookService := dependency.NewWebhookService(infra)
	httpHandlers := dependency.NewHTTPHandlers(
		authServices,
		userService,
		roleServices,
		permissionService,
		searchStack,
		mediaSvc,
		webhookService,
	)

	return &Container{
		Config:            infra.Config,
		Queue:             infra.Queue,
		QueueClient:       infra.QueueClient,
		DB:                infra.DB,
		JWT:               infra.JWT,
		Publisher:         infra.Publisher,
		MediaService:      mediaSvc,
		AuthService:       authServices.Auth,
		OAuthService:      authServices.OAuth,
		UserService:       userService,
		RoleService:       roleServices.Service,
		PermissionService: permissionService,
		SearchService:     searchStack.Service,
		SearchProcessor:   searchStack.Processor,
		OutboxService:     searchStack.Outbox,
		WebhookService:    webhookService,
		SearchClient:      infra.SearchClient,
		RoleRepo:          roleServices.Repo,
		AuthHandler:       httpHandlers.Auth,
		UserHandler:       httpHandlers.User,
		RoleHandler:       httpHandlers.Role,
		PermissionHandler: httpHandlers.Permission,
		SearchHandler:     httpHandlers.Search,
		MediaHandler:      httpHandlers.Media,
		WebhookHandler:    httpHandlers.Webhook,
	}
}

func (c *Container) Close() {
	if c.Publisher != nil {
		c.Publisher.Close()
	}
}
