package dto

import (
	"github.com/google/uuid"
)

type TimeslotDetail struct {
	ID        uuid.UUID `json:"timeslot_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}

type CreateTimeslotRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type EditTimeslotRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}