package repository

import (
	"gorm.io/gorm"

	"backend/internal/db_model"
)

type RestaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

func (r *RestaurantRepository) Create(restaurant *models.Restaurant) (*models.Restaurant, error) {

	result := r.db.Create(restaurant)
	if result.Error != nil {
		return nil, result.Error
	}

	restaurant, err := r.GetByUsername(restaurant.Username)
	if err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (r *RestaurantRepository) IsRestaurantExists(username string, email string) (bool, error) {

	var restaurant models.Restaurant

	// Check if name exists
	result := r.db.First(&restaurant, "username = ?", username)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	// Check if email exists
	result = r.db.First(&restaurant, "email = ?", email)
	if result.Error == nil {
		return true, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error
	}

	return false, nil
}

func (r *RestaurantRepository) GetByUsername(username string) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	result := r.db.First(&restaurant, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *RestaurantRepository) CreateBankAccount(bankAccount *models.BankAccount) error {
	return r.db.Create(bankAccount).Error
}