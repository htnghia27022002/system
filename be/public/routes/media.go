package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
)

func RegisterMediaRoutes(r *gin.RouterGroup, c *app.Container) {
	media := r.Group("/media")
	{
		media.GET("/avatars/:filename", c.MediaHandler.ServeAvatar)
	}
}
