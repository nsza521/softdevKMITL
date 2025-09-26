package dto

import (
	"github.com/google/uuid"
)

type TableTimeslotDetail struct {
	ID            uuid.UUID 	`json:"id"`
	Table         TableDetail   `json:"table"`
	// TimeslotID    uuid.UUID 	`json:"timeslot_id"`
	// Timeslot    TimeslotDetail	`json:"timeslot"`
	Status        string 		`json:"status"`
	ReservedSeats int    		`json:"reserved_seats"`
	// MaxSeats      int    		`json:"max_seats"`
}

