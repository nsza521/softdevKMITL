package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/order/dto"
	"backend/internal/order/usecase"
)

type OrderHandler struct {
	orderUC usecase.OrderUsecase
	queueUC usecase.QueueUsecase
}

// สร้าง Handler โดยรองรับทั้ง order usecase และ queue usecase
func NewOrderHandler(orderUC usecase.OrderUsecase, queueUC usecase.QueueUsecase) *OrderHandler {
	return &OrderHandler{
		orderUC: orderUC,
		queueUC: queueUC,
	}
}


// POST /orders  (reservation_id เป็น optional ใน body)
func (h *OrderHandler) Create(c *gin.Context) {
	
	var req dto.CreateFoodOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id") // JWT middleware inject
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.orderUC.Create(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GET /orders/queue
func (h *OrderHandler) GetQueue(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	var restaurantID string
	if role != "restaurant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden1"})
		return
	}else {
		restaurantID = userID
	}
	resp, err := h.queueUC.GetQueue(c.Request.Context(), userID, role, restaurantID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden2"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

