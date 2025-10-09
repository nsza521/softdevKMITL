package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/restaurant/interfaces"
)

func MapRestaurantRoutes(restaurantGroup *gin.RouterGroup, restaurantHandler interfaces.RestaurantHandler) {
	restaurantGroup.POST("/register", restaurantHandler.Register())
	restaurantGroup.POST("/login", restaurantHandler.Login())

	restaurantGroup.GET("/all", middleware.AuthMiddleware(), restaurantHandler.GetAll())
	restaurantGroup.POST("/upload_pic", middleware.AuthMiddleware(), restaurantHandler.UploadProfilePicture())
	restaurantGroup.PATCH("/status", restaurantHandler.ChangeStatus())
	restaurantGroup.POST("/logout", restaurantHandler.Logout())
	restaurantGroup.PUT("/edit/:id", middleware.AuthMiddleware(), restaurantHandler.EditInfo())
	restaurantGroup.PATCH("/editname/:id", middleware.AuthMiddleware(), restaurantHandler.UpdateName())
}

