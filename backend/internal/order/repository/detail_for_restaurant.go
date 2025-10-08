package repository

import (
	"context"

	"github.com/google/uuid"
	// "gorm.io/gorm"

	models "backend/internal/db_model"
)

type OrderDetailForRestaurant struct {
	Order           models.FoodOrder
	Items           []models.FoodOrderItem
	Options         []models.FoodOrderItemOption
	TableNumber     *string
	CustomerDisplay *string
}

// เพิ่มลงใน interface OrderRepository ด้วย
func (r *orderRepository) GetOrderDetailForRestaurant(
	ctx context.Context,
	orderID, restaurantID uuid.UUID,
) (*OrderDetailForRestaurant, error) {

	var ord models.FoodOrder
	if err := r.db.Debug().WithContext(ctx).
		First(&ord, "id = ?", orderID).Error; err != nil {
		return nil, err
	}

	// ดึงเฉพาะ item ของร้านนี้ (join menu_items เพื่อกรอง)
	var items []models.FoodOrderItem
	if err := r.db.WithContext(ctx).
		Table("food_order_items AS i").
		Select("i.*").
		Joins("JOIN menu_items m ON m.id = i.menu_item_id").
		Where("i.food_order_id = ? AND m.restaurant_id = ?", orderID, restaurantID).
		Scan(&items).Error; err != nil {
		return nil, err
	}

	// options ของ items เหล่านี้
	opts := make([]models.FoodOrderItemOption, 0)
	if len(items) > 0 {
		itemIDs := make([]uuid.UUID, 0, len(items))
		for _, it := range items {
			itemIDs = append(itemIDs, it.ID)
		}
		if err := r.db.WithContext(ctx).
			Where("food_order_item_id IN ?", itemIDs).
			Find(&opts).Error; err != nil {
			return nil, err
		}
	}

	// (optional) table number ถ้ามีการจอง — NOTE: ReservationID เป็น uuid.UUID
	var tableNumber *string
	if ord.ReservationID != uuid.Nil {
		type Row struct{ Number *string }
		var row Row
		_ = r.db.WithContext(ctx).Raw(`
			SELECT t.number
			FROM table_reservations tr
			JOIN tables t ON t.id = tr.table_id
			WHERE tr.id = ? LIMIT 1
		`, ord.ReservationID).Scan(&row).Error
		tableNumber = row.Number
	}

	return &OrderDetailForRestaurant{
		Order:           ord,
		Items:           items,
		Options:         opts,
		TableNumber:     tableNumber,
		CustomerDisplay: nil,
	}, nil
}
