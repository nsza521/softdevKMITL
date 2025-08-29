package models

import (
	// "github.com/google/uuid"
)

type Payment struct {
	Base
	OrderID     string  `gorm:"foreignKey:OrderID;not null;type:char(36)"`
	Amount      float32    `gorm:"not null"`
	Method      string     `gorm:"not null"` // e.g., "credit_card", "wallet", "true_money"
	Status      string     `gorm:"not null"` // e.g., "pending", "paid", "refund"
}