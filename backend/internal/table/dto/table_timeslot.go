package dto

import (
	"github.com/google/uuid"
)

type TableTimeslotDetail struct {
	ID            uuid.UUID 	`json:"id"`
	// Table         TableDetail    `json:"table"`
	// TimeslotID    uuid.UUID 	    `json:"timeslot_id"`
	// Timeslot      TimeslotDetail	`json:"timeslot"`
	TableRow	  string 		`json:"table_row"`
	TableCol	  string 		`json:"table_col"`
	Status        string 		`json:"status"`
	ReservedSeats int    		`json:"reserved_seats"`
	MaxSeats      int    		`json:"max_seats"`
}

type TableTimeslotResponse struct {
	StartTime	 string    `json:"start_time"`
	EndTime		 string    `json:"end_time"`
	TableTimeslots []TableTimeslotDetail `json:"table_timeslots"`
}

