package models

import (
	"github.com/google/uuid"
)

type Payment struct {
	Base
	// OrderID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount  float32   `gorm:"not null"`
	Method  uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status  string    `gorm:"not null; default:'pending'"` // e.g., "pending", "paid", "refund"
}

type PaymentMethod struct {
	Base
	Name        string `gorm:"not null; unique"` // e.g., "Promtpay", "KBANK", "SCB", "Wallet"
	Type        string `gorm:"not null"` // e.g., for "topup", "withdraw", "both", "paid", "refund", "all"
	Description string `gorm:"type:text"`
	ImageURL    *string `gorm:"type:text"`
}

type Transaction struct {
	Base
	UserID         uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PaymentMethodID uuid.UUID `gorm:"type:char(36);not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount         float32   `gorm:"not null"`
	Type           string    `gorm:"not null"` // e.g., "topup", "paid", "refund", "received"
}