package app

import (
	jwtmanager "be/common/jwt"
	"be/internal/app/dependency"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/middleware"
	"be/internal/module"
	"be/internal/modules"
	"be/internal/queue"
	"be/internal/repository/interfaces"
	authsvc "be/internal/services/auth"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	rolesvc "be/internal/services/role"
	usersvc "be/internal/services/user"
	"be/pkg/postgres"
	"be/public/handlers"
)

// Container holds the core services and handlers plus the enabled modules.
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
	RoleRepo          interfaces.RoleRepository
	UserRepo          interfaces.UserRepository
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleHandler       *handlers.RoleHandler
	PermissionHandler *handlers.PermissionHandler
	MediaHandler      *handlers.MediaHandler

	// Kernel is the core infrastructure handed to modules.
	Kernel *module.Kernel
	// Modules are the enabled optional modules (internal/modules/registry.go).
	Modules module.Set
}

// NewContainer wires the core and builds the enabled modules. It fails only when
// MODULES names an unknown module.
func NewContainer(cfg config.Config, db *postgres.Postgres) (*Container, error) {
	infra := dependency.NewInfra(cfg, db)
	mediaSvc := dependency.NewMediaService(infra)
	userRepo := dependency.NewUserRepository(infra)
	authServices := dependency.NewAuthServices(infra, mediaSvc)
	userService := dependency.NewUserService(infra, mediaSvc)
	roleServices := dependency.NewRoleServices(infra)
	permissionService := dependency.NewPermissionService(infra)
	httpHandlers := dependency.NewHTTPHandlers(
		infra.Config,
		authServices,
		userService,
		roleServices,
		permissionService,
		mediaSvc,
	)

	kernel := &module.Kernel{
		Config:    infra.Config,
		DB:        infra.DB,
		JWT:       infra.JWT,
		Publisher: infra.Publisher,
		Auth:      middleware.Auth(infra.JWT, roleServices.Repo, userRepo),
	}
	mods, err := modules.Build(kernel)
	if err != nil {
		infra.Close()
		return nil, err
	}

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
		RoleRepo:          roleServices.Repo,
		UserRepo:          userRepo,
		AuthHandler:       httpHandlers.Auth,
		UserHandler:       httpHandlers.User,
		RoleHandler:       httpHandlers.Role,
		PermissionHandler: httpHandlers.Permission,
		MediaHandler:      httpHandlers.Media,
		Kernel:            kernel,
		Modules:           mods,
	}, nil
}

func (c *Container) Close() {
	c.Modules.Close()
	if c.Publisher != nil {
		c.Publisher.Close()
	}
}
