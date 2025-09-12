package interfaces

import (
	"github.com/gin-gonic/gin"
)

type TableHandler interface {
	CreateTable() gin.HandlerFunc
	GetAllTables() gin.HandlerFunc
	CreateTimeslot() gin.HandlerFunc
	GetAllTimeslots() gin.HandlerFunc
	GetTableTimeslotByTimeslotID() gin.HandlerFunc
	// GetTableByID() gin.HandlerFunc
}