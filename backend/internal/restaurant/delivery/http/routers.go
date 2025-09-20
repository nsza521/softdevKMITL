package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/restaurant/interfaces"
)

func MapRestaurantRoutes(restaurantGroup *gin.RouterGroup, restaurantHandler interfaces.RestaurantHandler) {
	restaurantGroup.POST("/register", restaurantHandler.Register())
	restaurantGroup.GET("/all", middleware.AuthMiddleware(), restaurantHandler.GetAll())
	restaurantGroup.POST("/upload_pic", middleware.AuthMiddleware(), restaurantHandler.UploadProfilePicture())
	restaurantGroup.PATCH("/status", middleware.AuthMiddleware(), restaurantHandler.ChangeStatus()) // No auth middleware here, should be added later
}