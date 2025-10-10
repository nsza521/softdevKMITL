package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/order/dto"
	"backend/internal/order/usecase"
)

type OrderHandler struct{ uc usecase.OrderUsecase }

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler { return &OrderHandler{uc: uc} }

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

	resp, err := h.uc.Create(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}
