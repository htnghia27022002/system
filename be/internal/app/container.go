package app

import (
	"context"
	"log/slog"

	jwtmanager "be/common/jwt"
	"be/internal/app/dependency"
	"be/internal/audit"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/mail"
	"be/internal/maintenance"
	"be/internal/middleware"
	"be/internal/module"
	"be/internal/modules"
	"be/internal/queue"
	"be/internal/repository"
	"be/internal/repository/interfaces"
	accountsvc "be/internal/services/account"
	auditlogsvc "be/internal/services/auditlog"
	authsvc "be/internal/services/auth"
	dashboardsvc "be/internal/services/dashboard"
	"be/internal/services/media"
	permissionsvc "be/internal/services/permission"
	rolesvc "be/internal/services/role"
	usersvc "be/internal/services/user"
	"be/internal/settings"
	"be/pkg/postgres"
	"be/pkg/storage"
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
	Storage           storage.Storage
	Audit             audit.Recorder
	Settings          *settings.Service
	SystemHandler     *handlers.SystemHandler
	// Janitor prunes old audit entries and dead tokens (started by public.Run).
	Janitor           *maintenance.Janitor
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
	store, err := dependency.NewStorage(infra)
	if err != nil {
		infra.Close()
		return nil, err
	}
	mediaSvc := dependency.NewMediaService(store)
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

	recorder := audit.NewRecorder(repository.NewAuditLogRepository(infra.DB), slog.Default())
	authServices.Auth.SetAuditor(recorder)
	userService.SetAuditor(recorder)
	roleServices.Service.SetAuditor(recorder)
	accountService.SetAuditor(recorder)

	kernel := &module.Kernel{
		Config:    infra.Config,
		DB:        infra.DB,
		JWT:       infra.JWT,
		Publisher: infra.Publisher,
		Mailer:    mailer,
		Storage:   store,
		Audit:     recorder,
		Auth:      middleware.Auth(infra.JWT, roleServices.Repo, userRepo),
	}
	mods, err := modules.Build(kernel)
	if err != nil {
		infra.Close()
		return nil, err
	}

	// Settings are declared by the core and by the enabled modules.
	settingsSvc, err := settings.NewService(repository.NewSettingRepository(infra.DB), settings.CoreDefinitions(), mods.Settings())
	if err != nil {
		mods.Close()
		infra.Close()
		return nil, err
	}
	settingsSvc.SetAuditor(recorder)
	kernel.Settings = settingsSvc
	authServices.Auth.SetRegistrationGate(func(ctx context.Context) bool {
		return settingsSvc.Bool(ctx, settings.KeyRegistrationEnabled)
	})
	auditLogRepo := repository.NewAuditLogRepository(infra.DB)
	systemHandler := handlers.NewSystemHandler(
		auditlogsvc.NewService(auditLogRepo),
		settingsSvc,
		dashboardsvc.NewService(repository.NewStatsRepository(infra.DB), auditLogRepo),
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
		AccountService:    accountService,
		Mailer:            mailer,
		Storage:           store,
		Audit:             recorder,
		Settings:          settingsSvc,
		SystemHandler:     systemHandler,
		Janitor:           maintenance.NewJanitor(infra.DB, auditLogRepo, settingsSvc),
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
