package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/reservation/interfaces"
	"backend/internal/middleware"
)

func MapTableReservationRoutes(tableReservationGroup *gin.RouterGroup, tableReservationHandler interfaces.TableReservationHandler) {
	tableReservationGroup.Use(middleware.AuthMiddleware())
	tableReservationGroup.POST("/create", tableReservationHandler.CreateTableReservation())
	tableReservationGroup.GET("/history", tableReservationHandler.GetAllReservationHistory())
	tableReservationGroup.DELETE("/cancel/:reservation_id", tableReservationHandler.CancelTableReservation())
}