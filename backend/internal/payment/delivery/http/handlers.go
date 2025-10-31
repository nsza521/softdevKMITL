package http

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/payment/interfaces"
	// "backend/internal/customer/dto"
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