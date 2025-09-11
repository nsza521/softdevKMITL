package models

import (
	"github.com/google/uuid"
)

type Payment struct {
	Base
	OrderID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount  float32   `gorm:"not null"`
	Method  string    `gorm:"not null"` // e.g., "credit_card", "wallet", "true_money"
	Status  string    `gorm:"not null; default:'pending'"` // e.g., "pending", "paid", "refund"
}
