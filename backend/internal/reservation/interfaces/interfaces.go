package interfaces

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"backend/internal/db_model"
	"backend/internal/reservation/dto"
)

type TableReservationHandler interface {
	CreateTableReservation() gin.HandlerFunc
	CancelTableReservation() gin.HandlerFunc
	GetAllReservationHistory() gin.HandlerFunc
	// ConfirmTableReservation() gin.HandlerFunc
	// GetTableReservationByID() gin.HandlerFunc
}

type TableReservationRepository interface {
	// Table Reservation Repository
	CreateTableReservation(reservation *models.TableReservation) (*models.TableReservation, error)
	GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
	UpdateTableReservation(reservation *models.TableReservation) error

	// Table Reservation Members Repository
	CreateTableReservationMember(member *models.TableReservationMembers) error
	IsCustomerInReservation(reservationID uuid.UUID, customerID uuid.UUID) (bool, error)
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	DeleteReservationMember(reservationID uuid.UUID, customerID uuid.UUID) error
	GetAllReservationsByCustomerID(customerID uuid.UUID) ([]models.TableReservationMembers, error)
}

type TableReservationUsecase interface {
	CreateTableReservation(request *dto.CreateTableReservationRequest) (*dto.ReservationDetail, error)
	CancelTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	ConfirmTableReservation(reservationID uuid.UUID, customerID uuid.UUID) error
	GetAllMembersByReservationID(reservationID uuid.UUID) ([]models.TableReservationMembers, error)
	GetAllReservationHistory(customerID uuid.UUID) ([]dto.ReservationDetail, error)
	// UpdateTableReservation(reservation *models.TableReservation) error
	// GetTableReservationByID(id uuid.UUID) (*models.TableReservation, error)
}