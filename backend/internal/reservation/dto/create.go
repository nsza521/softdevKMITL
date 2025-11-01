package dto

import(
	"github.com/google/uuid"
)

type CreateTableReservationRequest struct {
	TableTimeslotID     uuid.UUID 	`json:"table_timeslot_id"`
	// CustomerID 			uuid.UUID 	`json:"customer_id" binding:"required"`
	Random           	bool      	`json:"random"` // e.g., true for "random", false for "not random"
	// Status          	string      `json:"status" binding:"omitempty"` // e.g., "pending", "confirmed", "canceled"
	Members 	  		[]Username	`json:"members" binding:"required"`
}

type CreateRandomTableReservationRequest struct {
	TimeslotID     uuid.UUID 	`json:"timeslot_id"`
}

type Username struct {
	Username string `json:"username"`
}