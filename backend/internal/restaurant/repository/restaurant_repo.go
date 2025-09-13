package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/google/uuid"

	models "backend/internal/db_model"
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

func (r *RestaurantRepository) GetByID(id uuid.UUID) (*models.Restaurant, error) {

	var restaurant models.Restaurant

	result := r.db.First(&restaurant, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &restaurant, nil
}

func (r *RestaurantRepository) GetAll() ([]*models.Restaurant, error) {

	var restaurants []*models.Restaurant

	result := r.db.Find(&restaurants)
	if result.Error != nil {
		return nil, result.Error
	}
	return restaurants, nil
}

func (r *RestaurantRepository) CreateBankAccount(bankAccount *models.BankAccount) error {
	return r.db.Create(bankAccount).Error
}

func (r *RestaurantRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Restaurant{}).
		Where("id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RestaurantRepository) IsOwner(ctx context.Context, restaurantID, userID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		// Table("restaurants").
		Model(&models.Restaurant{}).
		Where("id = ? AND owner_id = ?", restaurantID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RestaurantRepository) PartialUpdate(ctx context.Context, id string, changes map[string]any) error {
	if len(changes) == 0 {
		return nil // ไม่มีอะไรให้แก้
	}
	return r.db.WithContext(ctx).
		Model(&models.Restaurant{}).
		Where("id = ?", id).
		Updates(changes).Error
}