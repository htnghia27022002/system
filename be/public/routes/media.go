package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
)

func RegisterMediaRoutes(r *gin.RouterGroup, c *app.Container) {
	media := r.Group("/media")
	{
		// Keys look like avatars/<uuid>.png; old /api/media/avatars/... URLs still resolve.
		media.GET("/*key", c.MediaHandler.Serve)
	}
}
