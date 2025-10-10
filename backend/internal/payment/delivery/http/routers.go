package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/payment/interfaces"
	"backend/internal/middleware"
)

func MapPaymentRoutes(paymentGroup *gin.RouterGroup, paymentHandler interfaces.PaymentHandler) {
	paymentGroup.Use(middleware.AuthMiddleware())
	paymentGroup.GET("/methods/topup", paymentHandler.GetTopupPaymentMethods())
}