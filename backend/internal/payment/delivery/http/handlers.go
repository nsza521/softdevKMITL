package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/payment/dto"
	"backend/internal/payment/interfaces"
)

type PaymentHandler struct {
	paymentUsecase interfaces.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase interfaces.PaymentUsecase) interfaces.PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func getUserIDAndValidateRole(c *gin.Context, userRole string) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	role, exist := c.Get("role")
	if !exist || role.(string) != userRole {
		c.JSON(401, gin.H{"error": fmt.Sprintf("%s unauthorized", userRole)})
		return uuid.Nil, false
	}

	parseUserID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "invalid user id"})
		return uuid.Nil, false
	}

	return parseUserID, true
}

func (h *PaymentHandler) GetTopupPaymentMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, valid := getUserIDAndValidateRole(c, "customer")
		if !valid {
			return
		}

		paymentMethods, err := h.paymentUsecase.GetTopupPaymentMethods(customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"payment_methods": paymentMethods})
	}
}

func (h *PaymentHandler) TopupToWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, valid := getUserIDAndValidateRole(c, "customer")
		if !valid {
			return
		}

		var request dto.TopupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		err := h.paymentUsecase.TopupToWallet(customerID, &request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "customer top-up successful"})
	}
}

func (h *PaymentHandler) GetAllTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		userID, valid := getUserIDAndValidateRole(c, role.(string))
		if !valid {
			return
		}
		
		transactions, err := h.paymentUsecase.GetAllTransactions(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"transactions": transactions})
	}
}

func (h *PaymentHandler) PaidForFoodOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, valid := getUserIDAndValidateRole(c, "customer")
		if !valid {
			return
		}

		foodOrderIDParam := c.Param("food_order_id")
		foodOrderID, err := uuid.Parse(foodOrderIDParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid food order id"})
			return
		}

		response, err := h.paymentUsecase.PaidForFoodOrder(customerID, foodOrderID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": response})
	}
}

func (h *PaymentHandler) GetWithdrawPaymentMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID, valid := getUserIDAndValidateRole(c, "restaurant")
		if !valid {
			return
		}

		paymentMethods, err := h.paymentUsecase.GetWithdrawPaymentMethods(restaurantID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"payment_methods": paymentMethods})
	}
}

func (h *PaymentHandler) WithdrawFromWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID, valid := getUserIDAndValidateRole(c, "restaurant")
		if !valid {
			return
		}

		var request *dto.WithdrawRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("invalid request body: %v", err)})
			return
		}

		response, err := h.paymentUsecase.WithdrawFromWallet(restaurantID, request)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to withdraw from wallet: %v", err)})
			return
		}

		c.JSON(200, gin.H{"message": "restaurant withdrawal successful", "data": response})
	}
}