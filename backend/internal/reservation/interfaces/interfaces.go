package interfaces

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/db_model"
	"backend/internal/reservation/dto"
)

type TableReservationHandler interface {
	CreateTableReservation() gin.HandlerFunc
	CancelTableReservationMember() gin.HandlerFunc
	GetAllTableReservationHistory() gin.HandlerFunc
	GetTableReservationDetail() gin.HandlerFunc
	DeleteTableReservation() gin.HandlerFunc
	ConfirmTableReservation() gin.HandlerFunc
	ConfirmMemberInTableReservation() gin.HandlerFunc
}

type TableReservationRepository interface {
	// Table Reservation Repository
	CreateTableReservation(reservation *models.TableReservation) (*models.TableReservation, error)
	GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
	UpdateTableReservation(reservation *models.TableReservation) error
	DeleteTableReservation(reservationID uuid.UUID) error

	// Table Reservation Members Repository
	CreateTableReservationMember(member *models.TableReservationMembers) error
	IsCustomerInReservation(reservationID uuid.UUID, customerID uuid.UUID) (bool, error)
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	DeleteReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error
	GetAllReservationsByCustomerID(customerID uuid.UUID) ([]models.TableReservationMembers, error)
	GetTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) (*models.TableReservationMembers, error)
	UpdateTableReservationMember(member *models.TableReservationMembers) error

}

type TableReservationUsecase interface {
	CreateTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error)
	CreateNotRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error)
	CreateRandomTableReservation(request *dto.CreateTableReservationRequest, customerID uuid.UUID) (*dto.ReservationDetail, error)
	CancelTableReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error
	ConfirmTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	ConfirmMemberInTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	GetAllTableReservationHistory(customerID uuid.UUID) ([]dto.ReservationDetail, error)
	GetTableReservationDetail(reservationID uuid.UUID, customerID uuid.UUID) (*dto.ReservationDetail, error)
	DeleteTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	// UpdateTableReservation(reservation *models.TableReservation) error
	// GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
}