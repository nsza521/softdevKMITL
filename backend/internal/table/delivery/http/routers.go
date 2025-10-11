package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/table/interfaces"
)

func MapTableRoutes(tableGroup *gin.RouterGroup, tableHandler interfaces.TableHandler) {
	// Table routes
	tableGroup.POST("/create", tableHandler.CreateTable())

	// Timeslot routes
	tableGroup.POST("/timeslot/create", tableHandler.CreateTimeslot())
	tableGroup.GET("/timeslot/all", tableHandler.GetAllTimeslots())


	tableGroup.Use(middleware.AuthMiddleware())
	// Table routes
	tableGroup.GET("/all", tableHandler.GetAllTables())
	
	// Timeslot routes
	tableGroup.GET("/timeslot/:timeslot_id", tableHandler.GetTableTimeslotByTimeslotID())

	// TableTimeslot routes
	tableGroup.GET("/tabletimeslot/now", tableHandler.GetNowTableTimeslots())
	tableGroup.GET("/tabletimeslot/:table_timeslot_id", tableHandler.GetTableTimeslotByID())
}