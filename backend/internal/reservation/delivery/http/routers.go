package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/reservation/interfaces"
	"backend/internal/middleware"
)

func MapTableReservationRoutes(tableReservationGroup *gin.RouterGroup, tableReservationHandler interfaces.TableReservationHandler) {
	tableReservationGroup.Use(middleware.AuthMiddleware())
	tableReservationGroup.POST("/create", tableReservationHandler.CreateNotRandomTableReservation())
	tableReservationGroup.POST("/create/random", tableReservationHandler.CreateRandomTableReservation())
	tableReservationGroup.GET("/history", tableReservationHandler.GetAllTableReservationHistory())
	tableReservationGroup.GET("/history/all", tableReservationHandler.GetAlltableReservationByCustomerID())
	tableReservationGroup.GET("/:reservation_id/detail", tableReservationHandler.GetTableReservationDetail())
	tableReservationGroup.DELETE("/:reservation_id/cancel", tableReservationHandler.CancelTableReservationMember())
	tableReservationGroup.POST("/:reservation_id/confirm", tableReservationHandler.ConfirmTableReservation())
	tableReservationGroup.POST("/:reservation_id/confirm_member", tableReservationHandler.ConfirmMemberInTableReservation())
	tableReservationGroup.DELETE("/:reservation_id", tableReservationHandler.DeleteTableReservation())
}