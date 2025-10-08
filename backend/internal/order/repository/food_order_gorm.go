package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"backend/internal/db_model"
)

// ใช้ guard แบบ scoped โดยผูก reservation กับ restaurant_id เสมอเวลา query
type Reservation struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:char(36);index;not null"`
	CustomerID   uuid.UUID `gorm:"type:char(36);index;not null"`
	Status       string    `gorm:"type:varchar(32);not null"`
	ReserveDate  time.Time
}

func (Reservation) TableName() string { return "table_reservations" }

type OrderRepository interface {
	LoadReservation(ctx context.Context, id uuid.UUID) (*Reservation, error)
	CreateOrderTx(ctx context.Context, rsv *Reservation, order *models.FoodOrder, items []models.FoodOrderItem, opts []models.FoodOrderItemOption) error
}

type orderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) LoadReservation(ctx context.Context, id uuid.UUID) (*Reservation, error) {
	var res Reservation
	if err := r.db.WithContext(ctx).First(&res, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// สร้างออเดอร์แบบ transactional และผูก guard ตาม restaurant_id ของ reservation
func (r *orderRepository) CreateOrderTx(ctx context.Context, rsv *Reservation, order *models.FoodOrder, items []models.FoodOrderItem, opts []models.FoodOrderItemOption) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// double-check reservation ยังอยู่และเป็นร้านเดียวกัน
		var chk Reservation
		if err := tx.First(&chk, "id=? AND restaurant_id=?", rsv.ID, rsv.RestaurantID).Error; err != nil {
			return err
		}
		if chk.Status == "cancelled" {
			return errors.New("reservation cancelled")
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].FoodOrderID = order.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		if len(opts) > 0 {
			// ต้องใส่ FoodOrderItemID ให้ครบก่อนเรียก Create (usecase จะเซ็ตมาแล้ว)
			if err := tx.Create(&opts).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
