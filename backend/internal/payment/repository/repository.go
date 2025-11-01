package repository

import (
	"gorm.io/gorm"
	"github.com/google/uuid"

	"backend/internal/db_model"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// Transactions
func (r *PaymentRepository) CreateTransaction(transaction *models.Transaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepository) GetAllTransactionsByUserID(userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// Payment Methods
func (r *PaymentRepository) GetPaymentMethodByID(paymentMethodID uuid.UUID) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	if err := r.db.First(&method, "id = ?", paymentMethodID).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *PaymentRepository) GetPaymentMethods(methodType string) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	if err := r.db.Where("type IN ?", []string{methodType, "all"}).Find(&methods).Error; err != nil {
		return nil, err
	}
	return methods, nil
}