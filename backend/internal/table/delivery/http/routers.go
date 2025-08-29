package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/customer/interfaces"
)

func MapCustomerRoutes(customerGroup *gin.RouterGroup, customerHandler interfaces.CustomerHandler) {

	customerGroup.POST("/login", customerHandler.Login())
}