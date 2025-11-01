package repository

import (
	"context"
	"github.com/google/uuid"
	"time"
	// "gorm.io/gorm"

	models "backend/internal/db_model"
)

type OrderDetailForRestaurant struct {
	Order           models.FoodOrder
	Items           []models.FoodOrderItem
	Options         []models.FoodOrderItemOption
	TableLabel      *string
	TimeslotStart   *time.Time
	TimeslotEnd     *time.Time
	CustomerDisplay *string
}

// เพิ่มลงใน interface OrderRepository ด้วย
func (r *orderRepository) GetOrderDetailForRestaurant(
	ctx context.Context,
	orderID, restaurantID uuid.UUID,
) (*OrderDetailForRestaurant, error) {

	var ord models.FoodOrder
	if err := r.db.WithContext(ctx).
		First(&ord, "id = ?", orderID).Error; err != nil {
		return nil, err
	}

	// ดึงเฉพาะ items ของ "ร้านนี้" (กรองด้วย menu_items.restaurant_id)
	var items []models.FoodOrderItem
	if err := r.db.WithContext(ctx).
		Table("food_order_items AS i").
		Select("i.*").
		Joins("JOIN menu_items m ON m.id = i.menu_item_id").
		Where("i.food_order_id = ?", orderID).
		Scan(&items).Error; err != nil {
		return nil, err
	}

	// options ของ items เหล่านี้
	var opts []models.FoodOrderItemOption
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

	// ===== ดึงข้อมูลโต๊ะ/เวลาจากสคีมาของคุณ =====
	var tableLabel *string
	var tStart, tEnd *time.Time

	type rr struct {
		TableLabel *string    `gorm:"column:table_label"`
		StartTime  *time.Time `gorm:"column:start_time"`
		EndTime    *time.Time `gorm:"column:end_time"`
	}

	var out rr
	err := r.db.Debug().WithContext(ctx).
		Table("`table_reservations` AS tr").
		Select("CONCAT(t.`row`, t.`col`) AS table_label, ts.`start_time` AS start_time, ts.`end_time` AS end_time").
		Joins("JOIN `table_timeslots` tt ON tt.id = tr.table_timeslot_id").
		Joins("JOIN `tables` t ON t.id = tt.table_id").
		Joins("JOIN `timeslots` ts ON ts.id = tt.timeslot_id").
		Where("tr.id = ?", ord.ReservationID).
		Limit(1).
		Take(&out).Error

	if err == nil {
		tableLabel = out.TableLabel
		tStart = out.StartTime
		tEnd = out.EndTime
	}

	return &OrderDetailForRestaurant{
		Order:         ord,
		Items:         items,
		Options:       opts,
		TableLabel:    tableLabel,
		TimeslotStart: tStart,
		TimeslotEnd:   tEnd,
	}, nil
}
