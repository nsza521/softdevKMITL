package dto

import (
	"github.com/google/uuid"
)

type ConfirmedStatusDetail struct {
	ReservationID 	uuid.UUID 		`json:"reservation_id"`
	Members       	[]MemberStatus 	`json:"members"`
}

type MemberStatus struct {
	Username      string `json:"username"`
	Status        string `json:"status"`
}