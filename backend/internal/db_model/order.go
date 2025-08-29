package models

import (
	"time"
	// "github.com/google/uuid"
)

type FoodOrder struct {
	Base
	ReservationID 	string 			`gorm:"foreignKey:ReservationID;not null;type:char(36)"`
	PaymentID     	string   		`gorm:"foreignKey:PaymentID;not null;type:char(36)"`
	ExpectedTime  	time.Time   	`gorm:"not null"`
	Status        	string      	`gorm:"not null"` // e.g., "pending", "completed", "cancelled"
	// FoodOrderItems  []FoodOrderItem `gorm:"foreignKey:OrderID;not null"`
}

type FoodOrderItem struct {
	Base
	FoodOrderID string `gorm:"foreignKey:OrderID;not null;type:char(36)"`
	MenuItemID  string `gorm:"foreignKey:MenuItemID;not null;type:char(36)"`
	Quantity    int    `gorm:"not null"`
}