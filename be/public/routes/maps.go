package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

// RegisterMapsRoutes mounts Maps admin APIs under /admin/maps.
func RegisterMapsRoutes(admin *gin.RouterGroup, c *app.Container) {
	maps := admin.Group("/maps")
	{
		maps.GET("/search", middleware.RequireView("maps"), c.MapsHandler.Search)
		maps.GET("/places", middleware.RequireView("maps"), c.MapsHandler.ListPlaces)
		maps.POST("/places", middleware.RequireModify("maps"), c.MapsHandler.CreatePlace)
		maps.GET("/places/:id", middleware.RequireView("maps"), c.MapsHandler.GetPlace)
		maps.PATCH("/places/:id", middleware.RequireModify("maps"), c.MapsHandler.PatchPlace)
		maps.DELETE("/places/:id", middleware.RequireModify("maps"), c.MapsHandler.DeletePlace)

		maps.GET("/locations", middleware.RequireView("maps"), c.MapsHandler.ListLocations)
		maps.POST("/locations", middleware.RequireModify("maps"), c.MapsHandler.CreateLocation)
		maps.GET("/locations/:id", middleware.RequireView("maps"), c.MapsHandler.GetLocation)
		maps.PATCH("/locations/:id", middleware.RequireModify("maps"), c.MapsHandler.PatchLocation)
		maps.DELETE("/locations/:id", middleware.RequireModify("maps"), c.MapsHandler.DeleteLocation)

		maps.GET("/categories", middleware.RequireView("maps"), c.MapsHandler.ListCategories)

		maps.GET("/sources", middleware.RequireView("maps"), c.MapsHandler.ListSources)
		maps.POST("/sources/probe", middleware.RequireModify("maps"), c.MapsHandler.ProbeSource)
		maps.GET("/sources/:id", middleware.RequireView("maps"), c.MapsHandler.GetSource)
		maps.POST("/sources", middleware.RequireModify("maps"), c.MapsHandler.CreateSource)
		maps.PATCH("/sources/:id", middleware.RequireModify("maps"), c.MapsHandler.PatchSource)
		maps.DELETE("/sources/:id", middleware.RequireModify("maps"), c.MapsHandler.DeleteSource)

		maps.POST("/ingest", middleware.RequireModify("maps"), c.MapsHandler.StartIngest)
		maps.GET("/ingest", middleware.RequireView("maps"), c.MapsHandler.LatestIngest)
		maps.GET("/ingest/:id", middleware.RequireView("maps"), c.MapsHandler.GetIngest)
	}
}
