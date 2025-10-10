// internal/order/delivery/http/handler_detail.go
package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"backend/internal/order/usecase"
)

func (h *OrderHandler) GetDetailForRestaurant(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	restID, err := uuid.Parse(c.GetString("user_id")) // ต้องมีใน JWT middleware
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized" , "details": err.Error()})
		return
	}

	resp, err := h.uc.GetDetailForRestaurant(
		c.Request.Context(),
		usecase.GetDetailForRestaurantInput{
			OrderID:      orderID,
			RestaurantID: restID,
		},
	)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}