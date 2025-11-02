package dto

import (
	"github.com/google/uuid"
)

type ReservationDetail struct {
	CreateAt 	  		string      `json:"create_at"`
	ReservationID       uuid.UUID 	`json:"reservation_id"`
	TableTimeslotID     uuid.UUID 	`json:"table_timeslot_id"`
	ReservePeople       int       	`json:"reserve_people"`
	// Random           	bool      	`json:"random"`
	Status          	string      `json:"status"`
	Members 	  		[]Username	`json:"members"`
	TableRow 	   		string      `json:"table_row"`
	TableCol 	   		string      `json:"table_col"`
	StartTime   		string      `json:"start_time"`
	EndTime     		string      `json:"end_time"`
}

type ReservationMemberDetail struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	Members      []Username `json:"members"`
}

type RandomReservationDetail struct {
	ReservationID   uuid.UUID `json:"reservation_id"`
	TableTimeslotID uuid.UUID `json:"table_timeslot_id"`
}

type OwnerDetail struct {
	OwnerUsername   string    `json:"owner_username"`
	OwnerFirstname  string    `json:"owner_firstname"`
	TableTimeslotID uuid.UUID `json:"table_timeslot_id"`
}
