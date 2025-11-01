package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/table/dto"
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

func parseTableTimeslotID(c *gin.Context, param string) (uuid.UUID, error) {
	id := c.Param(param)
	if id == "" {
		return uuid.Nil, nil
	}

	parseID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %v", param, err)
	}

	return parseID, nil
}


// Table Handler
func (h *TableHandler) CreateTable() gin.HandlerFunc {
	return func (c *gin.Context) {
		var request dto.CreateTableRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		table, err := h.tableUsecase.CreateTable(request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"table": table})
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

func (h *TableHandler) GetTableByID() gin.HandlerFunc {
	return func (c *gin.Context) {

		tableID, err := parseTableTimeslotID(c, "table_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		table, err := h.tableUsecase.GetTableByID(tableID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"table": table})
	}
}

func (h *TableHandler) EditTableDetails() gin.HandlerFunc {
	return func (c *gin.Context) {

		tableID, err := parseTableTimeslotID(c, "table_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var request dto.EditTableRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.tableUsecase.EditTableDetails(tableID, &request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "table updated successfully"})
	}
}

func (h *TableHandler) DeleteTable() gin.HandlerFunc {
	return func (c *gin.Context) {

		tableID, err := parseTableTimeslotID(c, "table_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err := h.tableUsecase.DeleteTable(tableID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "table deleted successfully"})
	}
}

// Timeslot Handler
func (h *TableHandler) CreateTimeslot() gin.HandlerFunc {
	return func (c *gin.Context) {
		var request dto.CreateTimeslotRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		timeslot, err := h.tableUsecase.CreateTimeslot(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"timeslot": timeslot})
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

func (h *TableHandler) GetTimeslotByID() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslotID, err := parseTableTimeslotID(c, "timeslot_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		timeslot, err := h.tableUsecase.GetTimeslotByID(timeslotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"timeslot": timeslot})
	}
}

func (h *TableHandler) EditTimeslotDetails() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslotID, err := parseTableTimeslotID(c, "timeslot_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var request dto.EditTimeslotRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := h.tableUsecase.EditTimeslotDetails(timeslotID, &request); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "timeslot updated successfully"})
	}
}

func (h *TableHandler) DeleteTimeslot() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslotID, err := parseTableTimeslotID(c, "timeslot_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err := h.tableUsecase.DeleteTimeslot(timeslotID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "timeslot deleted successfully"})
	}
}


// TableTimeslot Handler
func (h *TableHandler) GetTableTimeslotByTimeslotID() gin.HandlerFunc {
	return func (c *gin.Context) {

		timeslotID, err := parseTableTimeslotID(c, "timeslot_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if timeslotID == uuid.Nil {
			c.JSON(400, gin.H{"error": "timeslot id is required"})
			return
		}

		response, err := h.tableUsecase.GetTableTimeslotByTimeslotID(timeslotID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, response)
	}
}

func (h *TableHandler) GetTableTimeslotByID() gin.HandlerFunc {
	return func (c *gin.Context) {

		tableTimeslotID, err := parseTableTimeslotID(c, "table_timeslot_id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if tableTimeslotID == uuid.Nil {
			c.JSON(400, gin.H{"error": "table timeslot id is required"})
			return
		}

		tableTimeslot, err := h.tableUsecase.GetTableTimeslotByID(tableTimeslotID)
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