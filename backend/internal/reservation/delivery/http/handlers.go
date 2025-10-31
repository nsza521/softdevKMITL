package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/reservation/dto"
	"backend/internal/reservation/interfaces"
)

type TableReservationHandler struct {
	tableReservationUsecase interfaces.TableReservationUsecase
}

func NewTableReservationHandler(tableReservationUsecase interfaces.TableReservationUsecase) interfaces.TableReservationHandler {
	return &TableReservationHandler{
		tableReservationUsecase: tableReservationUsecase,
	}
}

func getCustomerIDAndValidateRole(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	role, exist := c.Get("role")
	if !exist || role.(string) != "customer" {
		c.JSON(401, gin.H{"error": "customer unauthorized"})
		return uuid.Nil, false
	}

	customerID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "invalid user id"})
		return uuid.Nil, false
	}

	return customerID, true
}

func getReservationID(c *gin.Context) (uuid.UUID, error) {
	reservationIDParam := c.Param("reservation_id")
	reservationID, err := uuid.Parse(reservationIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid reservation id"})
		return uuid.Nil, err
	}
	return reservationID, nil
}

func (h *TableReservationHandler) CreateTableReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}
		var request dto.CreateTableReservationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		reservation, err := h.tableReservationUsecase.CreateTableReservation(&request, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"reservation": reservation, "message": "Table reservation created successfully"})
	}
}

func (h *TableReservationHandler) CancelTableReservationMember() gin.HandlerFunc {
	return func(c *gin.Context) {

		reservationID, err := getReservationID(c)
		if err != nil {
			return
		}

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		err = h.tableReservationUsecase.CancelTableReservationMember(reservationID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Reservation cancelled successfully"})
	}
}

func (h *TableReservationHandler) GetAllTableReservationHistory() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		reservations, err := h.tableReservationUsecase.GetAllTableReservationHistory(customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"reservations": reservations})
	}
}

func (h *TableReservationHandler) GetTableReservationDetail() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		reservationID, err := getReservationID(c)
		if err != nil {
			return
		}

		reservation, err := h.tableReservationUsecase.GetTableReservationDetail(reservationID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"reservation": reservation})
	}
}

func (h *TableReservationHandler) DeleteTableReservation() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		reservationID, err := getReservationID(c)
		if err != nil {
			return
		}

		err = h.tableReservationUsecase.DeleteTableReservation(reservationID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Reservation deleted successfully"})
	}
}

func (h *TableReservationHandler) ConfirmTableReservation() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		reservationID, err := getReservationID(c)
		if err != nil {
			return
		}

		err = h.tableReservationUsecase.ConfirmTableReservation(reservationID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Reservation confirmed successfully"})
	}
}

func (h *TableReservationHandler) ConfirmMemberInTableReservation() gin.HandlerFunc {
	return func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		reservationID, err := getReservationID(c)
		if err != nil {
			return
		}

		err = h.tableReservationUsecase.ConfirmMemberInTableReservation(reservationID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Member confirmed in reservation successfully"})
	}
}