package dto

import (
	"github.com/google/uuid"
)

type TableDetail struct {
	ID        uuid.UUID `json:"table_id"`
	Row       string    `json:"row"`
	Col       string    `json:"col"`
	MaxSeats  int       `json:"max_seats"`
}