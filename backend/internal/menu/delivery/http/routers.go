package http

import (
	"github.com/gin-gonic/gin"
	// "backend/internal/middleware"
)

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/:restaurantID/items", h.ListByRestaurant)
}
