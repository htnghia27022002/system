package app

import (
	jwtmanager "be/common/jwt"
	"be/internal/app/dependency"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/mail"
	"be/internal/middleware"
	"be/internal/module"
	"be/internal/modules"
	"be/internal/queue"
	"be/internal/repository/interfaces"
	accountsvc "be/internal/services/account"
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
	AccountService    *accountsvc.Service
	Mailer            mail.Mailer
	RoleRepo          interfaces.RoleRepository
	UserRepo          interfaces.UserRepository
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleHandler       *handlers.RoleHandler
	PermissionHandler *handlers.PermissionHandler
	MediaHandler      *handlers.MediaHandler
	AccountHandler    *handlers.AccountHandler

	// Kernel is the core infrastructure handed to modules.
	Kernel *module.Kernel
	// Modules are the enabled optional modules (internal/modules/registry.go).
	Modules module.Set
}

// Option overrides a dependency (tests).
type Option func(*options)

type options struct {
	mailer mail.Mailer
}

// WithMailer replaces the MAIL_DRIVER transport, e.g. with a recorder in tests.
func WithMailer(m mail.Mailer) Option {
	return func(o *options) { o.mailer = m }
}

// NewContainer wires the core and builds the enabled modules. It fails when
// MODULES names an unknown module or the mail transport is misconfigured.
func NewContainer(cfg config.Config, db *postgres.Postgres, opts ...Option) (*Container, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	infra := dependency.NewInfra(cfg, db)
	mailer := o.mailer
	if mailer == nil {
		var err error
		if mailer, err = dependency.NewMailer(infra); err != nil {
			infra.Close()
			return nil, err
		}
	}
	mediaSvc := dependency.NewMediaService(infra)
	userRepo := dependency.NewUserRepository(infra)
	authServices := dependency.NewAuthServices(infra, mediaSvc)
	userService := dependency.NewUserService(infra, mediaSvc)
	roleServices := dependency.NewRoleServices(infra)
	permissionService := dependency.NewPermissionService(infra)
	accountService := dependency.NewAccountService(infra, authServices, mailer)
	httpHandlers := dependency.NewHTTPHandlers(
		infra.Config,
		authServices,
		userService,
		roleServices,
		permissionService,
		mediaSvc,
		accountService,
	)

	kernel := &module.Kernel{
		Config:    infra.Config,
		DB:        infra.DB,
		JWT:       infra.JWT,
		Publisher: infra.Publisher,
		Mailer:    mailer,
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
		AccountService:    accountService,
		Mailer:            mailer,
		RoleRepo:          roleServices.Repo,
		UserRepo:          userRepo,
		AuthHandler:       httpHandlers.Auth,
		UserHandler:       httpHandlers.User,
		RoleHandler:       httpHandlers.Role,
		PermissionHandler: httpHandlers.Permission,
		MediaHandler:      httpHandlers.Media,
		AccountHandler:    httpHandlers.Account,
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
