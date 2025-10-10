package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/customer/interfaces"
	"backend/internal/middleware"
)

func MapCustomerRoutes(customerGroup *gin.RouterGroup, customerHandler interfaces.CustomerHandler) {
	customerGroup.POST("/register", customerHandler.Register())
	customerGroup.POST("/login", customerHandler.Login())

	customerGroup.Use(middleware.AuthMiddleware())
	customerGroup.GET("/profile", customerHandler.GetProfile())
	customerGroup.PATCH("/profile", customerHandler.EditProfile())
	customerGroup.POST("/fullname", customerHandler.GetFullnameByUsername())
	customerGroup.POST("/firstname", customerHandler.GetFirstNameByUsername())
	customerGroup.POST("/logout", customerHandler.Logout())
	// customerGroup.PUT("/change_password", customerHandler.ChangePassword())
}