package repository

import (
	"context"
	"time"
	
	"gorm.io/gorm"
	"github.com/google/uuid"
	"backend/internal/db_model" // ปรับ path ให้ตรงของคุณ
)

type OrderHistoryRepository interface {
	ListServedOrdersByRestaurantAndDay(
		ctx context.Context,
		restaurantID uuid.UUID,
		day time.Time, // เราจะใช้ day.Date() เทียบแค่ yyyy-mm-dd
	) ([]models.FoodOrder, error)
}

type orderHistoryRepository struct {
	db *gorm.DB
}

func NewOrderHistoryRepository(db *gorm.DB) OrderHistoryRepository {
	return &orderHistoryRepository{db: db}
}

func (r *orderHistoryRepository) ListServedOrdersByRestaurantAndDay(
	ctx context.Context,
	restaurantID uuid.UUID,
	day time.Time,
) ([]models.FoodOrder, error) {

	var orders []models.FoodOrder

	// สมมติเราดูวันที่ตาม OrderDate (timestamp ตอนเริ่มออเดอร์)
	dateOnly := day.Format("2006-01-02") // "YYYY-MM-DD"

	// status ที่ถือว่า "เสิร์ฟเสร็จแล้ว"
	statuses := []string{"served", "paid"}

	err := r.db.WithContext(ctx).
		// preload children (snapshot เราต้องใช้)
		Preload("Items.Options").
		// join menu_items เพื่อกรองร้าน
		Joins(`
			JOIN food_order_items foi ON foi.food_order_id = food_orders.id
			JOIN menu_items mi ON mi.id = foi.menu_item_id
		`).
		Where("mi.restaurant_id = ?", restaurantID).
		Where("food_orders.status IN ?", statuses).
		Where("DATE(food_orders.order_date) = ?", dateOnly).
		Group("food_orders.id").
		Order("food_orders.order_date DESC").
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	return orders, nil
}
