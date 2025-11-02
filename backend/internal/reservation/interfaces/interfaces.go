package interfaces

import (
	"time"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/db_model"
	"backend/internal/reservation/dto"
)

type TableReservationHandler interface {
	CreateNotRandomTableReservation() gin.HandlerFunc
	CreateRandomTableReservation() gin.HandlerFunc
	CancelTableReservationMember() gin.HandlerFunc
	GetAllTableReservationHistory() gin.HandlerFunc
	GetAllTableReservationByCustomerID() gin.HandlerFunc
	GetTableReservationDetail() gin.HandlerFunc
	GetTableReservationOwnerDetail() gin.HandlerFunc
	DeleteTableReservation() gin.HandlerFunc
	ConfirmTableReservation() gin.HandlerFunc
	ConfirmMemberInTableReservation() gin.HandlerFunc
}

type TableReservationRepository interface {
	// Table Reservation Repository
	CreateTableReservation(reservation *models.TableReservation) (*models.TableReservation, error)
	GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
	GetAllTableReservationByTableTimeslotID(tableTimeslotID uuid.UUID) ([]models.TableReservation, error)
	CountReservationsByCustomerAndDate(customerID uuid.UUID, date time.Time) (int64, error)
	UpdateTableReservation(reservation *models.TableReservation) error
	DeleteTableReservation(reservationID uuid.UUID) error


	// Table Reservation Members Repository
	CreateTableReservationMember(member *models.TableReservationMembers) error
	IsCustomerInReservation(reservationID uuid.UUID, customerID uuid.UUID) (bool, error)
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	DeleteReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error
	GetAllTableReservationsByCustomerID(customerID uuid.UUID) ([]models.TableReservationMembers, error)
	GetTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) (*models.TableReservationMembers, error)
	UpdateTableReservationMember(member *models.TableReservationMembers) error

}

type TableReservationUsecase interface {
	CreateTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error)
	CreateNotRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error)
	CreateRandomTableReservation(request *dto.CreateRandomTableReservationRequest, customerID uuid.UUID) (*dto.RandomReservationDetail, error)
	CancelTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error
	ConfirmTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	ConfirmMemberInTableReservation(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ConfirmedStatusDetail, error)
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	GetAllTableReservationHistory(customerID uuid.UUID) ([]dto.ReservationDetail, error)
	GetAllTableReservationByCustomerID(customerID uuid.UUID) ([]dto.ReservationDetail, error)
	GetTableReservationDetail(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ReservationDetail, error)
	GetTableReservationOwnerDetail(reservationID uuid.UUID) (*dto.OwnerDetail, error)
	DeleteTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	// UpdateTableReservation(reservation *models.TableReservation) error
	// GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
}