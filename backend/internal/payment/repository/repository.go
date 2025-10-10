package repository

import (
	"gorm.io/gorm"

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

func (r *PaymentRepository) GetTopupPaymentMethods() ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	if err := r.db.Where("type = ?", "topup OR all").Find(&methods).Error; err != nil {
		return nil, err
	}
	return methods, nil
}