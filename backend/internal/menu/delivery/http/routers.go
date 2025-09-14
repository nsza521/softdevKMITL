package http

import (
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/:restaurantID/items", h.ListByRestaurant)
	g.POST("/:restaurantID/items", middleware.AuthMiddleware(), h.Create) // middleware.RequireRole("restaurant_owner"),
	g.PATCH("/items/:itemID", middleware.AuthMiddleware(), h.Update)      // middleware.RequireRole("restaurant_owner"),
	g.DELETE("/items/:itemID", middleware.AuthMiddleware(), h.Delete)    // middleware.RequireRole("restaurant_owner"),
	g.POST("/items/:itemID/upload_pic", middleware.AuthMiddleware(), h.UploadMenuItemPicture) // middleware.RequireRole("restaurant_owner"),

}
