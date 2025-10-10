package interfaces

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/payment/dto"
	"backend/internal/db_model"

)

type PaymentHandler interface {
	GetTopupPaymentMethods() gin.HandlerFunc
}

type PaymentRepository interface {
	GetTopupPaymentMethods() ([]models.PaymentMethod, error)
}

type PaymentUsecase interface {
	GetTopupPaymentMethods(userID uuid.UUID) ([]dto.PaymentMethodDetail, error)
}