// internal/menu/repository/menutype_gorm.go
package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	models "backend/internal/db_model"
	iface "backend/internal/menu/interfaces"
)

type menuTypeRepo struct{ db *gorm.DB }

func NewMenuTypeRepository(db *gorm.DB) iface.MenuTypeRepository {
	return &menuTypeRepo{db: db}
}

func (r *menuTypeRepo) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuType, error) {
	var out []models.MenuType
	err := r.db.WithContext(ctx).
		Where("restaurant_id = ?", restaurantID).
		Order("type ASC").
		Find(&out).Error
	return out, err
}

func (r *menuTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.MenuType, error) {
	var mt models.MenuType
	if err := r.db.WithContext(ctx).First(&mt, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

func (r *menuTypeRepo) FindByName(ctx context.Context, restaurantID uuid.UUID, name string) (*models.MenuType, error) {
    var mt models.MenuType
    err := r.db.WithContext(ctx).
        Where("restaurant_id = ? AND type = ?", restaurantID, name).
        First(&mt).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &mt, err
}


func (r *menuTypeRepo) Create(ctx context.Context, mt *models.MenuType) error {
	return r.db.WithContext(ctx).Create(mt).Error
}

func (r *menuTypeRepo) Update(ctx context.Context, mt *models.MenuType) error {
	return r.db.WithContext(ctx).Save(mt).Error
}

func (r *menuTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.MenuType{}, "id = ?", id).Error
}

func (r *menuTypeRepo) GetMenuTypeRestaurantID(ctx context.Context, typeID uuid.UUID) (uuid.UUID, error) {
	type row struct{ RestaurantID uuid.UUID }
	var out row
	err := r.db.WithContext(ctx).
		Table("menu_types").
		Select("restaurant_id").
		Where("id = ?", typeID).
		Scan(&out).Error
	if err != nil {
		return uuid.Nil, err
	}
	return out.RestaurantID, nil
}
