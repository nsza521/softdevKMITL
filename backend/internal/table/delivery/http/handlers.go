package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	models "backend/internal/db_model"
	"backend/internal/table/interfaces"
)

type TableHandler struct {
	tableUsecase interfaces.TableUsecase
}

func NewTableHandler(tableUsecase interfaces.TableUsecase) interfaces.TableHandler {
	return &TableHandler{
		tableUsecase: tableUsecase,
	}
}

func (h *TableHandler) GetAllTables() gin.HandlerFunc {
	return func (c *gin.Context) {

		tables, err := h.tableUsecase.GetAllTables()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"tables": tables})
	}
}

func (h *TableHandler) CreateTable() gin.HandlerFunc {
	return func (c *gin.Context) {
		var table models.Table
		if err := c.ShouldBindJSON(&table); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := h.tableUsecase.CreateTable(&table)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Table created successfully"})
	}
}

func (h *TableHandler) CreateTimeslot() gin.HandlerFunc {
	return func (c *gin.Context) {
		var timeslot models.Timeslot
		if err := c.ShouldBindJSON(&timeslot); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := h.tableUsecase.CreateTimeslot(&timeslot)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Timeslot created successfully"})
	}
}

func (h *TableHandler) GetAllTimeslots() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslots, err := h.tableUsecase.GetAllTimeslots()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"timeslots": timeslots})
	}
}

func (h *TableHandler) GetTableTimeslotByTimeslotID() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslotID := c.Param("timeslot_id")
		if timeslotID == "" {
			c.JSON(400, gin.H{"error": "timeslot id is required"})
			return
		}

		parseTimeslotID, err := uuid.Parse(timeslotID)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid timeslot id"})
			return
		}

		response, err := h.tableUsecase.GetTableTimeslotByTimeslotID(parseTimeslotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, response)
	}
}

func (h *TableHandler) GetTableTimeslotByID() gin.HandlerFunc {
	return func (c *gin.Context) {

		tableTimeslotID := c.Param("table_timeslot_id")
		if tableTimeslotID == "" {
			c.JSON(400, gin.H{"error": "table timeslot id is required"})
			return
		}

		parseTableTimeslotID, err := uuid.Parse(tableTimeslotID)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid table timeslot id"})
			return
		}

		tableTimeslot, err := h.tableUsecase.GetTableTimeslotByID(parseTableTimeslotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"table_timeslot": tableTimeslot})
	}
}

func (h *TableHandler) GetNowTableTimeslots() gin.HandlerFunc {
	return func (c *gin.Context) {

		response, err := h.tableUsecase.GetNowTableTimeslots()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, response)
	}
}