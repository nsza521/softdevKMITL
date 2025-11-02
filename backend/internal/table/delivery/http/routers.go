package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/table/interfaces"
)

func MapTableRoutes(tableGroup *gin.RouterGroup, tableHandler interfaces.TableHandler) {
	// Table routes
	tableGroup.POST("/create", tableHandler.CreateTable())
	tableGroup.GET("/all", tableHandler.GetAllTables())
	tableGroup.GET("/:table_id", tableHandler.GetTableByID())
	tableGroup.PUT("/:table_id", tableHandler.EditTableDetails())
	tableGroup.DELETE("/:table_id", tableHandler.DeleteTable())

	// Timeslot routes
	tableGroup.POST("/timeslot/create", tableHandler.CreateTimeslot())
	tableGroup.GET("/timeslot/all", tableHandler.GetAllTimeslots())
	tableGroup.GET("/timeslot/:timeslot_id/details", tableHandler.GetTimeslotByID())
	tableGroup.PUT("/timeslot/:timeslot_id", tableHandler.EditTimeslotDetails())
	tableGroup.DELETE("/timeslot/:timeslot_id", tableHandler.DeleteTimeslot())

	// TableTimeslot routes
	tableGroup.GET("/table_timeslot/all/:timeslot_id", tableHandler.GetTableTimeslotByTimeslotID())
	tableGroup.GET("/table_timeslot/now", tableHandler.GetNowTableTimeslots())

	tableGroup.Use(middleware.AuthMiddleware())

	// TableTimeslot routes
	tableGroup.GET("/table_timeslot/:table_timeslot_id", tableHandler.GetTableTimeslotByID())
}
