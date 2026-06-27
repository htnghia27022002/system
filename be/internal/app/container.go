package app

import (
	jwtmanager "be/internal/common/jwt"
	"be/internal/config"
	"be/internal/repository"
	"be/internal/repository/interfaces"
	authsvc "be/internal/services/auth"
	permissionsvc "be/internal/services/permission"
	rolesvc "be/internal/services/role"
	usersvc "be/internal/services/user"
	"be/public/handlers"

	"gorm.io/gorm"
)

type Container struct {
	Config            config.Config
	DB                *gorm.DB
	JWT               *jwtmanager.Manager
	AuthService       *authsvc.Service
	OAuthService      *authsvc.OAuthService
	UserService       *usersvc.Service
	RoleService       *rolesvc.Service
	PermissionService *permissionsvc.Service
	RoleRepo          interfaces.RoleRepository
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleHandler       *handlers.RoleHandler
	PermissionHandler *handlers.PermissionHandler
}

func NewContainer(cfg config.Config, db *gorm.DB) *Container {
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	jwtManager := jwtmanager.NewManager(cfg)
	authService := authsvc.NewService(authRepo, userRepo, roleRepo, jwtManager, cfg.JWTRefreshTTL)
	oauthService := authsvc.NewOAuthService(cfg, authRepo, roleRepo, authService)
	userService := usersvc.NewService(userRepo)
	roleService := rolesvc.NewService(roleRepo)
	permissionService := permissionsvc.NewService(permissionRepo)

	return &Container{
		Config:            cfg,
		DB:                db,
		JWT:               jwtManager,
		AuthService:       authService,
		OAuthService:      oauthService,
		UserService:       userService,
		RoleService:       roleService,
		PermissionService: permissionService,
		RoleRepo:          roleRepo,
		AuthHandler:       handlers.NewAuthHandler(authService, oauthService),
		UserHandler:       handlers.NewUserHandler(userService),
		RoleHandler:       handlers.NewRoleHandler(roleService),
		PermissionHandler: handlers.NewPermissionHandler(permissionService),
	}
}
