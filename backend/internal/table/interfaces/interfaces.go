package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/db_model"
)

type TableHandler interface {
	CreateTable() gin.HandlerFunc
	GetAllTables() gin.HandlerFunc
	CreateTimeslot() gin.HandlerFunc
	GetAllTimeslots() gin.HandlerFunc
	GetTableTimeslotByTimeslotID() gin.HandlerFunc
	// GetTableByID() gin.HandlerFunc
}

type TableRepository interface {
	CreateTable(table *models.Table) error
	GetAllTables() ([]*models.Table, error)
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]*models.Timeslot, error)
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]*models.TableTimeslot, error)
	// GetTableByID(id uint) (*models.Table, error)
}

type TableUsecase interface {
	CreateTable(table *models.Table) error
	GetAllTables() ([]*models.Table, error)
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]*models.Timeslot, error)
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]*models.TableTimeslot, error)
	// GetTableByID(id uint) (*models.Table, error)
}