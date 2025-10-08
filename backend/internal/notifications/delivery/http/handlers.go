package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/notifications/dto"
	"backend/internal/notifications/interfaces"

)

type NotificationHandler struct {
	uc interfaces.NotiHandler
}

func NewNotificationHandler(r *gin.Engine, uc interfaces.NotiHandler) {
	h := &NotificationHandler{uc: uc}

	v1 := r.Group("/api/v1")
	{
		// v1.GET("/notifications", h.List)
		v1.POST("/notifications/mock", h.MockCreate)
		v1.PATCH("/notifications/:id/read", h.MarkRead)
		v1.PATCH("/notifications/read-all", h.MarkAllRead)
	}
}

// GET /api/v1/notifications?... (receiverId, receiverType, isRead?, page?, pageSize?, sort?)
func (h *NotificationHandler) List(c *gin.Context) {
	var q dto.ListQuery
	// if err := c.ShouldBindQuery(&q); err != nil {
	if err := c.ShouldBindUri(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.uc.List(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// POST /api/v1/notifications/mock
// body: { "receiverId": "...", "receiverType": "customer|restaurant", "count": 10 }
func (h *NotificationHandler) MockCreate(c *gin.Context) {
	var req dto.MockCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.uc.MockCreate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// PATCH /api/v1/notifications/:id/read
// body: { "isRead": true }
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req dto.MarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.MarkRead(c.Request.Context(), id, req.IsRead); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// PATCH /api/v1/notifications/read-all?receiverId=...&receiverType=...
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	receiverIdStr := c.Query("receiverId")
	receiverType := c.Query("receiverType")

	rid, err := uuid.Parse(receiverIdStr)
	if err != nil || (receiverType != "customer" && receiverType != "restaurant") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiverId or receiverType"})
		return
	}

	affected, err := h.uc.MarkAllRead(c.Request.Context(), rid, receiverType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": affected})
}
