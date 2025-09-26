package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/db_model"
	"backend/internal/table/dto"
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
	GetAllTables() ([]models.Table, error)
	GetTableByID(id uuid.UUID) (*models.Table, error)

	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]models.Timeslot, error)
	GetTimeslotByID(id uuid.UUID) (*models.Timeslot, error)
	IsTimeslotAvailable(id uuid.UUID) (bool, error)

	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]models.TableTimeslot, error)
	GetTableTimeslotByID(id uuid.UUID) (*models.TableTimeslot, error)
	UpdateTableTimeslot(tableTimeslot *models.TableTimeslot) error
}

type TableUsecase interface {
	CreateTable(table *models.Table) error
	GetAllTables() ([]dto.TableDetail, error)

	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]dto.TimeslotDetail, error)
	// GetTimeslotByID(id uuid.UUID) (*dto.TimeslotDetail, error)
	// IsTimeslotAvailable(id uuid.UUID) (bool, error)

	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]dto.TableTimeslotDetail, error)
	// GetTableByID(id uint) (*models.Table, error)
}