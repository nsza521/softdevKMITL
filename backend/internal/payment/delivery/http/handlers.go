package http

import (
	// "github.com/gin-gonic/gin"

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
