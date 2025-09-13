package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/customer/interfaces"
	"backend/internal/customer/dto"
	user "backend/internal/user/dto"
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

		var request *dto.RegisterCustomerRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.customerUsecase.Register(request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "registration successful"})
	}
}

func (h *CustomerHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request *user.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		token, err := h.customerUsecase.Login(request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(200, gin.H{"token": token})
	}
}

func (h *CustomerHandler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		parseCustomerID, err := uuid.Parse(customerID.(string))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user id"})
			return
		}

		profile, err := h.customerUsecase.GetProfile(parseCustomerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, profile)
	}
}

func (h *CustomerHandler) EditProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		parseCustomerID, err := uuid.Parse(customerID.(string))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid user id"})
			return
		}

		var request *dto.EditProfileRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.customerUsecase.EditProfile(parseCustomerID, request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "profile updated successfully"})
	}
}