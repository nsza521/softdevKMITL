package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/user/interfaces"
	"backend/internal/middleware"
)

func MapUserRoutes(userGroup *gin.RouterGroup, userHandler interfaces.UserHandler) {
	userGroup.POST("/login", userHandler.Login())
	userGroup.POST("/logout", middleware.AuthMiddleware(), userHandler.Logout())
	// userGroup.PUT("/change_password", middleware.AuthMiddleware(), userHandler.ChangePassword())
}