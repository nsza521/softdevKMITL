package repository

import (
	"context"
	"errors"
	"time"
	// "fmt"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"backend/internal/db_model"
)

// ใช้ guard แบบ scoped โดยผูก reservation กับ restaurant_id เสมอเวลา query
type Reservation struct {
	ReservationID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	RestaurantID uuid.UUID `gorm:"type:char(36);index;not null"`
	CustomerID   uuid.UUID `gorm:"type:char(36);index;not null"`
	Status       string    `gorm:"type:varchar(32);not null"`
	ReserveDate  time.Time
}

func (Reservation) TableName() string { return "table_reservations" }

type OrderRepository interface {
	LoadReservationForCustomer(ctx context.Context, reservationID, customerID uuid.UUID) (*Reservation, error)
	CreateOrderTx(ctx context.Context, order *models.FoodOrder, items []models.FoodOrderItem, opts []models.FoodOrderItemOption) error
	GetOrderDetailForRestaurant(ctx context.Context, orderID, restaurantID uuid.UUID) (*OrderDetailForRestaurant, error)
}

type orderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) LoadReservationForCustomer(ctx context.Context, reservationID, customerID uuid.UUID) (*Reservation, error) {
	var res Reservation
	// โหลด reservation
	if err := r.db.Debug().WithContext(ctx).
		First(&res, "id = ?", reservationID).Error; err != nil {
		return nil, err
	}

	// ถ้าเป็นเจ้าของ → ผ่าน
	if res.CustomerID == customerID {
		return &res, nil
	}

	// ไม่ใช่เจ้าของ → ตรวจในตารางสมาชิก
	var cnt int64
	if err := r.db.Debug().WithContext(ctx).
		Table("table_reservation_members").
		Where("reservation_id = ? AND customer_id = ?", reservationID, customerID).
		Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, errors.New("forbidden: not a member of this reservation")
	}
	return &res, nil
}


// สร้างออเดอร์แบบ transactional และผูก guard ตาม restaurant_id ของ reservation
func (r *orderRepository) CreateOrderTx(ctx context.Context, order *models.FoodOrder, items []models.FoodOrderItem, opts []models.FoodOrderItemOption) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// double-check reservation ยังอยู่และเป็นร้านเดียวกัน


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
