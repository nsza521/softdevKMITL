package repository

import (
	"gorm.io/gorm"
)

type TableReservationRepository struct {
	db *gorm.DB
}

func NewTableReservationRepository(db *gorm.DB) *TableReservationRepository {
	return &TableReservationRepository{
		db: db,
	}
}