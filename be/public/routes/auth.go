package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

func RegisterAuthRoutes(r *gin.RouterGroup, c *app.Container) {
	// Per-IP limit on endpoints that verify a password or create accounts
	// (AUTH_RATE_LIMIT_PER_MINUTE; 0 disables). Refresh is not limited: it already
	// requires a valid refresh token, and many users can share one NAT IP.
	limited := middleware.NewRateLimiter(c.Config.Auth.RateLimitPerMinute, time.Minute).Middleware()

	auth := r.Group("/auth")
	{
		auth.POST("/login", limited, c.AuthHandler.Login)
		auth.POST("/register", limited, c.AuthHandler.Register)
		auth.POST("/refresh", c.AuthHandler.Refresh)
		auth.POST("/logout", c.AuthHandler.Logout)
		auth.GET("/oauth/providers", c.AuthHandler.OAuthProviders)
		auth.GET("/oauth/:provider/start", c.AuthHandler.OAuthStart)
		auth.POST("/oauth/:provider/callback", limited, c.AuthHandler.OAuthCallback)

		protected := auth.Group("")
		protected.Use(c.Kernel.Auth)
		protected.GET("/me", c.AuthHandler.Me)
		protected.PATCH("/profile", c.AuthHandler.UpdateProfile)
		protected.POST("/profile/avatar", c.AuthHandler.UploadAvatar)
		protected.POST("/change-password", limited, c.AuthHandler.ChangePassword)
		protected.GET("/sessions", c.AuthHandler.ListSessions)
		protected.POST("/sessions/revoke-others", c.AuthHandler.RevokeOtherSessions)
		protected.DELETE("/sessions/:id", c.AuthHandler.RevokeSession)
	}
}
