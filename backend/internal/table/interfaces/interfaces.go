package interfaces

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/db_model"
	"backend/internal/table/dto"
)

type TableHandler interface {
	// Table Handler
	CreateTable() gin.HandlerFunc
	GetAllTables() gin.HandlerFunc
	GetTableByID() gin.HandlerFunc
	EditTableDetails() gin.HandlerFunc
	DeleteTable() gin.HandlerFunc
	
	// Timeslot Handler
	CreateTimeslot() gin.HandlerFunc
	GetAllTimeslots() gin.HandlerFunc
	GetTimeslotByID() gin.HandlerFunc
	EditTimeslotDetails() gin.HandlerFunc
	DeleteTimeslot() gin.HandlerFunc

	// TableTimeslot Handler
	GetTableTimeslotByTimeslotID() gin.HandlerFunc
	GetTableTimeslotByID() gin.HandlerFunc
	GetNowTableTimeslots() gin.HandlerFunc
}

type TableRepository interface {
	// Table Repository
	CreateTable(table *models.Table) error
	GetAllTables() ([]models.Table, error)
	GetTableByID(id uuid.UUID) (*models.Table, error)
	UpdateTable(table *models.Table) error
	DeleteTable(table *models.Table) error

	// Timeslot Repository
	CreateTimeslot(timeslot *models.Timeslot) error
	GetAllTimeslots() ([]models.Timeslot, error)
	GetTimeslotByID(id uuid.UUID) (*models.Timeslot, error)
	UpdateTimeslot(timeslot *models.Timeslot) error
	DeleteTimeslot(timeslot *models.Timeslot) error
	GetActiveTimeslot(now time.Time) (*models.Timeslot, error)

	// TableTimeslot Repository
	CreateTableTimeslot(tableTimeslot *models.TableTimeslot) error
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) ([]models.TableTimeslot, error)
	GetTableTimeslotByID(id uuid.UUID) (*models.TableTimeslot, error)
	UpdateTableTimeslot(tableTimeslot *models.TableTimeslot) error
	DeleteTableTimeslotByTimeslotID(timeslotID uuid.UUID) error
	GetAllAvailableTableTimeslot(timeslotID uuid.UUID) ([]models.TableTimeslot, error)
}

type TableUsecase interface {
	// Table Usecase
	CreateTable(request dto.CreateTableRequest) (*dto.TableDetail, error)
	GetTableByID(id uuid.UUID) (*dto.TableDetail, error)
	GetAllTables() ([]dto.TableDetail, error)
	EditTableDetails(id uuid.UUID, request *dto.EditTableRequest) error
	DeleteTable(id uuid.UUID) error

	// Timeslot Usecase
	CreateTimeslot(timeslot *dto.CreateTimeslotRequest) (*dto.TimeslotDetail, error)
	GetAllTimeslots() ([]dto.TimeslotDetail, error)
	GetTimeslotByID(id uuid.UUID) (*dto.TimeslotDetail, error)
	EditTimeslotDetails(id uuid.UUID, request *dto.EditTimeslotRequest) error
	DeleteTimeslot(id uuid.UUID) error

	// TableTimeslot Usecase
	GetTableTimeslotByTimeslotID(timeslotID uuid.UUID) (*dto.TableTimeslotResponse, error)
	GetTableTimeslotByID(id uuid.UUID) (*dto.TableTimeslotDetail, error)
	GetNowTableTimeslots() (*dto.TableTimeslotResponse, error)
}