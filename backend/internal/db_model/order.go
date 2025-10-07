package models

import (
	"time"
	"github.com/google/uuid"
)

// หลัก: ออร์เดอร์ถูกผูกกับ Reservation (ถ้ามี) และร้าน
type FoodOrder struct {
	Base
	RestaurantID  uuid.UUID  `gorm:"type:char(36);not null;index"`
	ReservationID *uuid.UUID `gorm:"type:char(36);index"` // multi-customer ผ่าน reservation นี้
	PaymentID     *uuid.UUID `gorm:"type:char(36)"`
	Channel       string     `gorm:"type:enum('web','walk_in');not null"`
	ExpectedTime  time.Time  `gorm:"not null"` // เวลาเสร็จโดยรวมที่คาดหวัง
	Status        string     `gorm:"not null;default:'pending'"`
	Notes         string     `gorm:"type:text"`
	TotalAmount   float64    `gorm:"not null;default:0"`

	Items      []FoodOrderItem     `gorm:"foreignKey:FoodOrderID"`
	Customers  []Customer          `gorm:"many2many:food_order_customers"` // ผู้มีสิทธิ์ร่วม/ผู้ร่วมสั่ง
}
type FoodOrderItem struct {
  Base
  FoodOrderID uuid.UUID `gorm:"type:char(36);not null;index"`
  MenuItemID  uuid.UUID `gorm:"type:char(36);not null;index"`
  Quantity    int       `gorm:"not null"`
  UnitPrice   float64   `gorm:"not null"`                 // snapshot menu
  AddOnSubtotal float64 `gorm:"not null;default:0"`       // Σ option ของบรรทัด (ก่อนคูณ qty)
  LineTotal   float64   `gorm:"not null"`                 // (UnitPrice + AddOnSubtotal) * Quantity
  CreatedByCustomerID *uuid.UUID `gorm:"type:char(36);index"`

  Options []FoodOrderItemOption `gorm:"foreignKey:FoodOrderItemID"`
}

type FoodOrderItemOption struct {
  Base
  FoodOrderItemID uuid.UUID `gorm:"type:char(36);not null;index"`
  AddOnOptionID   uuid.UUID `gorm:"type:char(36);not null;index"`
  Qty             int       `gorm:"not null;default:1"`
  PriceDelta      float64   `gorm:"not null"` // snapshot
}

// ผู้ร่วมสั่งในออร์เดอร์นี้ (many-to-many)
type FoodOrderCustomer struct {
	FoodOrderID uuid.UUID `gorm:"type:char(36);primaryKey"`
	CustomerID  uuid.UUID `gorm:"type:char(36);primaryKey"`
	Role        string    `gorm:"type:enum('owner','contributor');not null;default:'contributor'"`
	// owner = คนที่เปิดออร์เดอร์, contributor = ร่วมสั่ง
}

type FoodOrderHistory struct {
	Base
	CustomerID    uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RestaurantID  uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FoodOrderID   uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentID     uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalAmount   float32   `gorm:"not null"`
}