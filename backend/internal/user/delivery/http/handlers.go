package http

import (
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

	"backend/internal/user/dto"
	"backend/internal/user/interfaces"
)

type UserHandler struct {
	userUsecase interfaces.UserUsecase
}

func NewUserHandler(userUsecase interfaces.UserUsecase) interfaces.UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request *dto.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		token, err := h.userUsecase.Login(request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
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