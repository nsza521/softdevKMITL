package interfaces

import (
	"time"
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
	GetTableTimeslotByID() gin.HandlerFunc
	GetNowTableTimeslots() gin.HandlerFunc
}

type TableRepository interface {
	// Table Repository
	CreateTable(table *models.Table) error
	GetAllTables() ([]models.Table, error)
	GetTableByID(id uuid.UUID) (*models.Table, error)

	// Timeslot Repository
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]models.Timeslot, error)
	GetTimeslotByID(id uuid.UUID) (*models.Timeslot, error)
	GetActiveTimeslot(now time.Time) (*models.Timeslot, error)

	// TableTimeslot Repository
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]models.TableTimeslot, error)
	GetTableTimeslotByID(id uuid.UUID) (*models.TableTimeslot, error)
	UpdateTableTimeslot(tableTimeslot *models.TableTimeslot) error
	GetAvailableTableTimeslot(timeslotID uuid.UUID) (*models.TableTimeslot, error)
}

type TableUsecase interface {
	// Table Usecase
	CreateTable(table *models.Table) error
	GetAllTables() ([]dto.TableDetail, error)

	// Timeslot Usecase
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]dto.TimeslotDetail, error)
	// GetTimeslotByID(id uuid.UUID) (*dto.TimeslotDetail, error)
	// IsTimeslotAvailable(id uuid.UUID) (bool, error)

	// TableTimeslot Usecase
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) (*dto.TableTimeslotResponse, error)
	GetTableTimeslotByID(id uuid.UUID) (*dto.TableTimeslotDetail, error)
	GetNowTableTimeslots() (*dto.TableTimeslotResponse, error)
}