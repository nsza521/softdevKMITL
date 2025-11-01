package dto

import (
	"github.com/google/uuid"
)

type TransactionDetail struct {
	TransactionID   uuid.UUID  `json:"transaction_id"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id"`
	Amount          float32    `json:"amount"`
	Type            string     `json:"type"`
	CreatedAt       string     `json:"created_at"`
}