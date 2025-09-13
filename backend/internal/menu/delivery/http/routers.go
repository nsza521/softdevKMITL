package http

import (
	"github.com/gin-gonic/gin"
	// "backend/internal/middleware"
)

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/:restaurantID/items", h.ListByRestaurant)
	g.POST("/:restaurantID/items", h.Create) // middleware.RequireRole("restaurant_owner"),
	g.PATCH("/items/:itemID", h.Update)      // middleware.RequireRole("restaurant_owner"),
	g.DELETE("/items/:itemID", h.Delete)    // middleware.RequireRole("restaurant_owner"),

}
