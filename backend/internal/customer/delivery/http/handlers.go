package http

import (
	"time"
	"strings"

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

func getCustomerIDAndValidateRole(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	role, exist := c.Get("role")
	if !exist || role.(string) != "customer" {
		c.JSON(401, gin.H{"error": "customer unauthorized"})
		return uuid.Nil, false
	}

	customerID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "invalid user id"})
		return uuid.Nil, false
	}

	return customerID, true
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
		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		profile, err := h.customerUsecase.GetProfile(customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, profile)
	}
}

func (h *CustomerHandler) EditProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		var request *dto.EditProfileRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.customerUsecase.EditProfile(customerID, request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "profile updated successfully"})
	}
}

func (h *CustomerHandler) GetFullnameByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}
		
		var request *dto.GetFullnameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fullname, err := h.customerUsecase.GetFullnameByUsername(customerID, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, fullname)
	}
}

func (h *CustomerHandler) GetFirstNameByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}
		
		var request *dto.GetFullnameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		firstname, err := h.customerUsecase.GetFirstnameByUsername(customerID, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, firstname)
	}
}

func (h *CustomerHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, ok := getCustomerIDAndValidateRole(c)
		if !ok {
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

		err := h.customerUsecase.Logout(tokenString, expiry)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "logged out successfully"})
	}
}