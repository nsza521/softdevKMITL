package http

import (
	// "github.com/gin-gonic/gin"

	"backend/internal/reservation/interfaces"
	// "backend/internal/customer/dto"
)

type TableReservationHandler struct {
	tableReservationUsecase interfaces.TableReservationUsecase
}

func NewTableReservationHandler(tableReservationUsecase interfaces.TableReservationUsecase) interfaces.TableReservationHandler {
	return &TableReservationHandler{
		tableReservationUsecase: tableReservationUsecase,
	}
}
