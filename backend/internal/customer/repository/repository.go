package repository

import (
	// "fmt"
	"gorm.io/gorm"

	"backend/internal/db_model"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	result := r.db.Create(customer)
	return result.Error
}

func (r *CustomerRepository) IsCustomerExists(username string, email string) (bool, error) {

	var customer models.Customer

	// Check if username exists
	result := r.db.First(&customer, "username = ?", username)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	// Check if email exists
	result = r.db.First(&customer, "email = ?", email)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	return false, nil
}

func (r *CustomerRepository) GetByUsername(username string) (*models.Customer, error) {

	var customer models.Customer
	
	result := r.db.First(&customer, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &customer, nil
}