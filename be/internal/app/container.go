package app

import (
	jwtmanager "be/common/jwt"
	"be/internal/app/dependency"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/queue"
	"be/internal/repository/interfaces"
	authsvc "be/internal/services/auth"
	ingestsvc "be/internal/services/maps/ingest"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	rolesvc "be/internal/services/role"
	searchsvc "be/internal/services/search"
	usersvc "be/internal/services/user"
	webhooksvc "be/internal/services/webhook"
	"be/pkg/postgres"
	"be/public/handlers"
)

type Container struct {
	Config            config.Config
	Queue             queue.Config
	QueueClient       *queue.Client
	DB                *postgres.Postgres
	JWT               *jwtmanager.Manager
	Publisher         *publisher.Publisher
	MediaService      *media.Service
	AuthService       *authsvc.Service
	OAuthService      *authsvc.OAuthService
	UserService       *usersvc.Service
	RoleService       *rolesvc.Service
	PermissionService *permissionsvc.Service
	SearchService     *searchsvc.Service
	WebhookService    *webhooksvc.Service
	IngestService     *ingestsvc.Service
	RoleRepo          interfaces.RoleRepository
	UserRepo          interfaces.UserRepository
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleHandler       *handlers.RoleHandler
	PermissionHandler *handlers.PermissionHandler
	SearchHandler     *handlers.SearchHandler
	MediaHandler      *handlers.MediaHandler
	WebhookHandler    *handlers.WebhookHandler
	MapsHandler       *handlers.MapsHandler
	AddressHandler    *handlers.AddressHandler
}

func NewContainer(cfg config.Config, db *postgres.Postgres) *Container {
	infra := dependency.NewInfra(cfg, db)
	searchService := dependency.NewSearchService(infra)
	mediaSvc := dependency.NewMediaService(infra)
	userRepo := dependency.NewUserRepository(infra)
	authServices := dependency.NewAuthServices(infra, mediaSvc)
	userService := dependency.NewUserService(infra, mediaSvc)
	roleServices := dependency.NewRoleServices(infra)
	permissionService := dependency.NewPermissionService(infra)
	webhookService := dependency.NewWebhookService(infra)
	addressService := dependency.NewAddressService(infra)
	mapsServices := dependency.NewMapsServices(infra, addressService)
	httpHandlers := dependency.NewHTTPHandlers(
		authServices,
		userService,
		roleServices,
		permissionService,
		searchService,
		mediaSvc,
		webhookService,
		mapsServices,
		addressService,
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
		SearchService:     searchService,
		WebhookService:    webhookService,
		IngestService:     mapsServices.Ingest,
		RoleRepo:          roleServices.Repo,
		UserRepo:          userRepo,
		AuthHandler:       httpHandlers.Auth,
		UserHandler:       httpHandlers.User,
		RoleHandler:       httpHandlers.Role,
		PermissionHandler: httpHandlers.Permission,
		SearchHandler:     httpHandlers.Search,
		MediaHandler:      httpHandlers.Media,
		WebhookHandler:    httpHandlers.Webhook,
		MapsHandler:       httpHandlers.Maps,
		AddressHandler:    httpHandlers.Address,
	}
}

func (c *Container) Close() {
	if c.Publisher != nil {
		c.Publisher.Close()
	}
}
