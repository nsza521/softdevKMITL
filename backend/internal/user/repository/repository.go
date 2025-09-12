package repository

import (
	// "fmt"
	"gorm.io/gorm"

	"backend/internal/db_model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetCustomerByUsername(username string) (*models.Customer, error) {

	var customer models.Customer
	
	result := r.db.First(&customer, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &customer, nil
}

func (r *UserRepository) GetRestaurantByUsername(username string) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	result := r.db.First(&restaurant, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}