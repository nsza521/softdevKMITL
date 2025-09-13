package http

import "github.com/gin-gonic/gin"

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/restaurant/:restaurantID/items", h.ListByRestaurant)
}
