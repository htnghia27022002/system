package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

func RegisterAuthRoutes(r *gin.RouterGroup, c *app.Container) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", c.AuthHandler.Login)
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/refresh", c.AuthHandler.Refresh)
		auth.POST("/logout", c.AuthHandler.Logout)
		auth.GET("/oauth/providers", c.AuthHandler.OAuthProviders)
		auth.GET("/oauth/:provider/start", c.AuthHandler.OAuthStart)
		auth.POST("/oauth/:provider/callback", c.AuthHandler.OAuthCallback)

		protected := auth.Group("")
		protected.Use(middleware.Auth(c.JWT, c.RoleRepo))
		protected.GET("/me", c.AuthHandler.Me)
	}
}
