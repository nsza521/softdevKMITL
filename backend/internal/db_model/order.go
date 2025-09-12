package models

import (
	"time"
	"github.com/google/uuid"
)

type FoodOrder struct {
	Base
	ReservationID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentID     uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExpectedTime  time.Time `gorm:"not null"`
	Status        string    `gorm:"not null;default:'pending'"` // e.g., "pending", "completed", "cancelled"

	FoodOrderItems []FoodOrderItem `gorm:"foreignKey:FoodOrderID"` // one-to-many
}

type FoodOrderItem struct {
	Base
	FoodOrderID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MenuItemID  uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity    int       `gorm:"not null"`
}
