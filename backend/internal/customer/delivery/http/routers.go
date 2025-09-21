package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/customer/interfaces"
	"backend/internal/middleware"
)

func MapCustomerRoutes(customerGroup *gin.RouterGroup, customerHandler interfaces.CustomerHandler) {
	customerGroup.POST("/register", customerHandler.Register())
	customerGroup.POST("/login", customerHandler.Login())
	customerGroup.GET("/profile", middleware.AuthMiddleware(), customerHandler.GetProfile())
	customerGroup.PATCH("/profile", middleware.AuthMiddleware(), customerHandler.EditProfile())
	customerGroup.GET("/fullname", middleware.AuthMiddleware(), customerHandler.GetFullnameByUsername())
	// customerGroup.POST("/logout", middleware.AuthMiddleware(), customerHandler.Logout())
	// customerGroup.PUT("/change_password", middleware.AuthMiddleware(), customerHandler.ChangePassword())
}