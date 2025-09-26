package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/restaurant/interfaces"
)

func MapRestaurantRoutes(restaurantGroup *gin.RouterGroup, restaurantHandler interfaces.RestaurantHandler) {
	restaurantGroup.POST("/register", restaurantHandler.Register())
	restaurantGroup.POST("/login", restaurantHandler.Login())

	restaurantGroup.Use(middleware.AuthMiddleware())
	restaurantGroup.GET("/all", restaurantHandler.GetAll())
	restaurantGroup.POST("/upload_pic", restaurantHandler.UploadProfilePicture())
	restaurantGroup.PATCH("/status", restaurantHandler.ChangeStatus())
	restaurantGroup.POST("/logout", restaurantHandler.Logout())
}