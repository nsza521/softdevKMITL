package http

import (
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

	"backend/internal/user/dto"
	"backend/internal/user/interfaces"
	customerInterfaces "backend/internal/customer/interfaces"
	restaurantInterfaces "backend/internal/restaurant/interfaces"
)

type UserHandler struct {
	userUsecase         interfaces.UserUsecase
	customerUsecase     customerInterfaces.CustomerUsecase
	restaurantUsecase   restaurantInterfaces.RestaurantUsecase
}

func NewUserHandler(userUsecase interfaces.UserUsecase, customerusecase customerInterfaces.CustomerUsecase,
	restaurantusecase restaurantInterfaces.RestaurantUsecase) interfaces.UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		customerUsecase: customerusecase,
		restaurantUsecase: restaurantusecase,
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request *dto.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		token, _ := h.customerUsecase.Login(request)
		if token == "" {
			token, _ = h.restaurantUsecase.Login(request)
			if token == "" {
				c.JSON(401, gin.H{"error": "invalid username or password"})
				return
			}
		}
		
		c.JSON(200, gin.H{"token": token})
	}
}

func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		expClaim, exists := c.Get("exp")
		if !exists {
			c.JSON(500, gin.H{"error": "token expiry not found"})
			return
		}

		expFloat, ok := expClaim.(float64)
		if !ok {
			c.JSON(500, gin.H{"error": "invalid token expiry format"})
			return
		}

		expiry := time.Unix(int64(expFloat), 0)

		err := h.userUsecase.Logout(tokenString, expiry)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "logged out successfully"})
	}
}