package dto

import (
	"github.com/google/uuid"
)

type TableDetail struct {
	ID       uuid.UUID `json:"table_id"`
	TableRow string    `json:"table_row"`
	TableCol string    `json:"table_col"`
	MaxSeats int       `json:"max_seats"`
}

type CreateTableRequest struct {
	TableRow string `json:"table_row" binding:"required"`
	TableCol string `json:"table_col" binding:"required"`
	MaxSeats int    `json:"max_seats" binding:"required"`
}

type EditTableRequest struct {
	TableRow string `json:"table_row"`
	TableCol string `json:"table_col"`
	MaxSeats int    `json:"max_seats"`
}