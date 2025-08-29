package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/customer/interfaces"
	// "backend/internal/customer/dto"
)

type CustomerHandler struct {
	customerUsecase interfaces.CustomerUsecase
}

func NewCustomerHandler(customerUsecase interfaces.CustomerUsecase) interfaces.CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(201, gin.H{"message": "customer registered"})
	}
}

func (h *CustomerHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		token, err := h.customerUsecase.Login(loginRequest.Username, loginRequest.Password)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
	}
}