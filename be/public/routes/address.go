package routes

import (
	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/middleware"
)

// RegisterAddressRoutes mounts country and division catalogs under /address.
func RegisterAddressRoutes(r *gin.RouterGroup, c *app.Container) {
	addresses := r.Group("/address")
	addresses.Use(middleware.Auth(c.JWT, c.RoleRepo, c.UserRepo))
	{
		addresses.GET("/countries", c.AddressHandler.ListCountries)
		addresses.GET("/divisions", c.AddressHandler.ListDivisions)
	}
}
