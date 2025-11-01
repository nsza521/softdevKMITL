package http

import (
	"backend/internal/order/dto"
	"backend/internal/order/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
    orderUC        usecase.OrderUsecase
    queueUC        usecase.QueueUsecase
    historyUC      usecase.OrderHistoryUsecase
}

func NewOrderHandler(
    orderUC usecase.OrderUsecase,
    queueUC usecase.QueueUsecase,
    historyUC usecase.OrderHistoryUsecase,
) *OrderHandler {
    return &OrderHandler{
        orderUC:   orderUC,
        queueUC:   queueUC,
        historyUC: historyUC,
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

func (h *OrderHandler) GetHistoryForDay(c *gin.Context) {
	// ดึง restaurant_id จาก context (middleware.AuthMiddleware + RequireRole ควรใส่ให้แล้วใน router)
	userID := c.GetString("user_id")
	role := c.GetString("role")
	var restaurantID string
	if role != "restaurant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden1"})
		return
	}else {
		restaurantID = userID
	}
	if restaurantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing restaurant scope"})
		return
	}

	restaurantUUID, err := uuid.Parse(restaurantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant_id"})
		return
	}

	// parse query date=YYYY-MM-DD
	dayStr := c.Query("date")

	var day time.Time
	if dayStr == "" {
		day = time.Now()
	} else {
		parsed, parseErr := time.Parse("2006-01-02", dayStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date, use YYYY-MM-DD"})
			return
		}
		day = parsed
	}

	resp, err := h.historyUC.GetServedHistoryForDay(
		c.Request.Context(),
		restaurantUUID,
		day,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}


// PATCH /restaurant/orders/:orderID/status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("orderID")

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	res, err := h.orderUC.UpdateStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateStatusResponse{
		ID:     res.ID,
		Status: res.Status,
	})
}