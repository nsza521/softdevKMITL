package models

import (
	"time"

	"github.com/google/uuid"
)

type FoodOrder struct {
	ID              uuid.UUID       `gorm:"type:char(36);primaryKey"`
	ReservationID   uuid.UUID       `gorm:"type:char(36);index;"` // อนุญาตให้ null ได้
	CustomerID      uuid.UUID       `gorm:"type:char(36);index;"`
	// Customer   *Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedByUserID uuid.UUID       `gorm:"type:char(36);not null;index"`
	Status          string          `gorm:"type:enum('pending','preparing','served','paid','cancelled');default:'pending';not null"`
	Channel         string     		`gorm:"type:enum('walk_in','reservation','delivery');not null;default:'walk_in'"`
	OrderDate       time.Time       `gorm:"not null"`
	ExpectedReceive *time.Time
	TotalAmount     float64         `gorm:"not null;default:0"`
	Note            *string
	Items           []FoodOrderItem `gorm:"foreignKey:FoodOrderID"`
}



func (FoodOrder) TableName() string { return "food_orders" }

type FoodOrderItem struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	FoodOrderID  uuid.UUID `gorm:"type:char(36);index;not null"`
	CustomerID      uuid.UUID       `gorm:"type:char(36);index;not null"`
	MenuItemID   uuid.UUID `gorm:"type:char(36);index;not null"`

	// Snapshot จากเมนูตอนสั่ง
	MenuName     string  `gorm:"type:varchar(255);not null"`
	UnitPrice    float64 `gorm:"not null"`
	TimeTakenMin int     `gorm:"not null"`

	Quantity int     `gorm:"not null"`
	Subtotal float64 `gorm:"not null;default:0"`
	Note     *string

	Options []FoodOrderItemOption `gorm:"foreignKey:FoodOrderItemID"`
}

func (FoodOrderItem) TableName() string { return "food_order_items" }

type FoodOrderItemOption struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey"`
	FoodOrderItemID uuid.UUID `gorm:"type:char(36);index;not null"`
	AddOnOptionID   uuid.UUID `gorm:"type:char(36);index;not null"`

	// Snapshot
	GroupID    uuid.UUID `gorm:"type:char(36);index;not null"`
	GroupName  string    `gorm:"type:varchar(255);not null"`
	OptionName string    `gorm:"type:varchar(255);not null"`
	PriceDelta float64   `gorm:"not null"`
	Qty        int       `gorm:"not null;default:1"`
}

func (FoodOrderItemOption) TableName() string { return "food_order_item_options" }
