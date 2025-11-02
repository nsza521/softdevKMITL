package dto

import (
	"github.com/google/uuid"
)

type PaymentSummary struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	FoodOrderID   uuid.UUID `json:"food_order_id"`
	TotalMembers  int       `json:"total_members"`
	PaidMembers   int       `json:"paid_members"`
}
