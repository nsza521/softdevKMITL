package repository

import (
	"context"

	"backend/internal/db_model"

	"gorm.io/gorm"
)

type QueueRepository interface {
	ListPendingOrdersByRestaurant(ctx context.Context, restaurantID string) ([]models.FoodOrder, error)
	CountMenuItemBelongsToRestaurant(menuItemID, restaurantID string) (bool, error)
}

type queueRepository struct {
	db *gorm.DB
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) ListPendingOrdersByRestaurant(ctx context.Context, restaurantID string) ([]models.FoodOrder, error) {
	var orders []models.FoodOrder

	err := r.db.WithContext(ctx).
		Preload("Items.Options").
		Joins("JOIN food_order_items ON food_order_items.food_order_id = food_orders.id").
		Joins("JOIN menu_items ON menu_items.id = food_order_items.menu_item_id").
		Where("food_orders.status = ?", "pending").
		Where("menu_items.restaurant_id = ?", restaurantID).
		Group("food_orders.id").
		Order("food_orders.order_date ASC").
		Find(&orders).Error

	return orders, err
}

func (r *queueRepository) ListServedOrdersByRestaurant(ctx context.Context, restaurantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("food_orders").
		Where("status = ? AND restaurant_id = ?", "served", restaurantID).
		Count(&count).Error
	return count, err
}

// helper ใช้ตอน filter item ต่อ order
func (r *queueRepository) CountMenuItemBelongsToRestaurant(menuItemID, restaurantID string) (bool, error) {
	var count int64
	err := r.db.Table("menu_items").
		Where("id = ? AND restaurant_id = ?", menuItemID, restaurantID).
		Count(&count).Error
	return count > 0, err
}
