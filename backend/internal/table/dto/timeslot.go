package dto

import (
	"github.com/google/uuid"
)

type TimeslotDetail struct {
	ID        uuid.UUID `json:"timeslot_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}