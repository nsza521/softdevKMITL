package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models "backend/internal/db_model"
	iface "backend/internal/menu/interfaces"
)

type menuRepo struct{ db *gorm.DB }

func NewMenuRepository(db *gorm.DB) iface.MenuRepository { return &menuRepo{db: db} }

// internal/menu/repository/repository.go
func (r *menuRepo) ListMenuByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := r.db.WithContext(ctx).
		Model(&models.MenuItem{}).
		Select("menu_items.*").
		Joins("JOIN menu_tags  mt  ON mt.menu_item_id = menu_items.id").
		Joins("JOIN menu_types mty ON mty.id = mt.menu_type_id").
		Where("mty.restaurant_id = ?", restaurantID).
		Group("menu_items.id").
		// ⬇️ สำคัญ: โหลดเฉพาะ MenuTypes ของร้านนี้เท่านั้น
		Preload("MenuTypes", "menu_types.restaurant_id = ?", restaurantID).
		Find(&items).Error
	return items, err
}

func (r *menuRepo) RestaurantExists(ctx context.Context, restaurantID uuid.UUID) error {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Restaurant{}).
		Where("id = ?", restaurantID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

