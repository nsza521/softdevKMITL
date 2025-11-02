package repository

import (
	"context"
	// "fmt"
	"gorm.io/gorm"
	"github.com/google/uuid"
	
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

func (r *CustomerRepository) GetByID(id uuid.UUID) (*models.Customer, error) {

	var customer models.Customer
	
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	
	return &customer, nil
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) ListServedOrdersByCustomer(
    ctx context.Context,
    customerID string,
) ([]models.FoodOrder, error) {

    var orders []models.FoodOrder

    err := r.db.Debug().WithContext(ctx).
        Preload("Items.Options"). // <-- ปรับชื่อ Preload ให้ตรงกับ model คุณ
        Where("customer_id = ?", customerID).
        Where("status = ?", "paid").
        Order("order_date DESC").
        Find(&orders).Error

    if err != nil {
        return nil, err
    }

    return orders, nil
}
