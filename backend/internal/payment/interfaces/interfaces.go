package interfaces

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/payment/dto"
	"backend/internal/db_model"

)

type PaymentHandler interface {
	GetTopupPaymentMethods() gin.HandlerFunc
	TopupToWallet() gin.HandlerFunc
	GetAllTransactions() gin.HandlerFunc
}

type PaymentRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	CreatePaymentMethod(method *models.PaymentMethod) error
	GetPaymentMethodByID(paymentMethodID uuid.UUID) (*models.PaymentMethod, error)
	GetPaymentMethodsByType(methodType string) ([]models.PaymentMethod, error)
	GetAllTransactionsByUserID(userID uuid.UUID) ([]models.Transaction, error)
	GetAllPaymentMethods() ([]models.PaymentMethod, error)
}

type PaymentUsecase interface {
	GetTopupPaymentMethods(userID uuid.UUID) ([]dto.PaymentMethodDetail, error)
	TopupToWallet(userID uuid.UUID, request *dto.TopupRequest) error
	GetAllTransactions(userID uuid.UUID) ([]dto.TransactionDetail, error)
}