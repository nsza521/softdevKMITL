package dto

import (
	"github.com/google/uuid"
)

type ReservationDetail struct {
	ReservationID       uuid.UUID 	`json:"reservation_id"`
	TableTimeslotID     uuid.UUID 	`json:"table_timeslot_id"`
	ReservePeople       int       	`json:"reserve_people"`
	Random           	bool      	`json:"random"`
	Status          	string      `json:"status"`
	Members 	  		[]Username	`json:"members"`
}

type ReservationMemberDetail struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	Members      []Username `json:"members"`
}

