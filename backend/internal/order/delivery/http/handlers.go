package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/order/dto"
	"backend/internal/order/usecase"
)

type OrderHandler struct{ uc usecase.OrderUsecase }

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

// POST /restaurant/reservations/:reservationID/orders
func (h *OrderHandler) Create(c *gin.Context) {
	rid, err := uuid.Parse(c.Param("reservationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	var req dto.CreateFoodOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JWT middleware ควร inject "user_id" (customer) ใน context
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.uc.Create(c.Request.Context(), rid, req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
