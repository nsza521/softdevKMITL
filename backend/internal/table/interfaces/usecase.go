package interfaces

import (
	"github.com/google/uuid"

	"backend/internal/db_model"
)

type TableUsecase interface {
	CreateTable(table *models.Table) error
	GetAllTables() ([]*models.Table, error)
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]*models.Timeslot, error)
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]*models.TableTimeslot, error)
	// GetTableByID(id uint) (*models.Table, error)
}