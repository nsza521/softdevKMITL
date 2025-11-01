package dto

import (
	"github.com/google/uuid"
)

type PaymentMethodDetail struct {
	PaymentMethodID uuid.UUID `json:"payment_method_id"`
	Name          	string    `json:"name"`
	// ImageURL        *string   `json:"image_url"`
}