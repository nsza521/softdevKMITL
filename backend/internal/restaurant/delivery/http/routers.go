package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/restaurant/interfaces"
)

func MapRestaurantRoutes(restaurantGroup *gin.RouterGroup, restaurantHandler interfaces.RestaurantHandler) {

	restaurantGroup.POST("/register", restaurantHandler.Register())
}