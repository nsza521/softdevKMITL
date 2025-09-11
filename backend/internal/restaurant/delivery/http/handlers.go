package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/restaurant/interfaces"
	"backend/internal/restaurant/dto"
	user "backend/internal/user/dto"
)

type RestaurantHandler struct {
	restaurantUsecase interfaces.RestaurantUsecase
}

func NewRestaurantHandler(restaurantUsecase interfaces.RestaurantUsecase) interfaces.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restaurantUsecase,
	}
}

func (h *RestaurantHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest *dto.RegisterRestaurantRequest
		if err := c.ShouldBindJSON(&registerRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := h.restaurantUsecase.Register(registerRequest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
}

func (h *RestaurantHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request user.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		token, err := h.restaurantUsecase.Login(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
	}
}