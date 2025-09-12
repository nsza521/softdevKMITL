package http

import (
	"github.com/gin-gonic/gin"

	"backend/internal/middleware"
	"backend/internal/table/interfaces"
)

func MapTableRoutes(tableGroup *gin.RouterGroup, tableHandler interfaces.TableHandler) {
	tableGroup.POST("/create", tableHandler.CreateTable())
	tableGroup.GET("/all", middleware.AuthMiddleware(), tableHandler.GetAllTables())
	tableGroup.POST("/timeslot/create", tableHandler.CreateTimeslot())
	tableGroup.GET("/timeslot/all", tableHandler.GetAllTimeslots())
	tableGroup.GET("/timeslot/:timeslot_id", middleware.AuthMiddleware(), tableHandler.GetTableTimeslotByTimeslotID())
}