package dto

import (
	"github.com/google/uuid"
)

type TopupRequest struct {
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Amount          float32    `json:"amount"`
}