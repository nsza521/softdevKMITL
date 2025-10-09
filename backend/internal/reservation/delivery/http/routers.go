package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/reservation/interfaces"
	"backend/internal/middleware"
)

func MapTableReservationRoutes(tableReservationGroup *gin.RouterGroup, tableReservationHandler interfaces.TableReservationHandler) {
	tableReservationGroup.Use(middleware.AuthMiddleware())
	tableReservationGroup.POST("/create", tableReservationHandler.CreateTableReservation())
	tableReservationGroup.GET("/history", tableReservationHandler.GetAllTableReservationHistory())
	tableReservationGroup.GET("/:reservation_id/detail", tableReservationHandler.GetTableReservationDetail())
	tableReservationGroup.DELETE("/:reservation_id/cancel", tableReservationHandler.CancelTableReservation())
}