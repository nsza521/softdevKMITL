package http

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/payment/interfaces"
	"backend/internal/payment/dto"
)

type PaymentHandler struct {
	paymentUsecase interfaces.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase interfaces.PaymentUsecase) interfaces.PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) GetTopupPaymentMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		customerID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(500, gin.H{"error": "invalid user id"})
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
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		customerID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(500, gin.H{"error": "invalid user id"})
			return
		}

		var request dto.TopupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		err = h.paymentUsecase.TopupToWallet(customerID, &request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "customer top-up successful"})
	}
}

func (h *PaymentHandler) GetAllTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		customerID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(500, gin.H{"error": "invalid user id"})
			return
		}

		transactions, err := h.paymentUsecase.GetAllTransactions(customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"transactions": transactions})
	}
}