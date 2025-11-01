package http

import (
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

// MapFoodOrderRoutes กำหนดเส้นทาง API ของ FoodOrder
func MapFoodOrderRoutes(g *gin.RouterGroup, h *OrderHandler) {
	// ===== Order หลัก =====
	g.POST("", middleware.AuthMiddleware(), h.Create)
	g.GET("/:orderID/detail", middleware.AuthMiddleware(), h.GetDetailForRestaurant)

	// ===== Queue (รายการที่กำลังทำ) =====
	g.GET("/queue", middleware.AuthMiddleware(), h.GetQueue)

	// ===== History (ออเดอร์ที่เสิร์ฟหรือจ่ายเงินแล้ว) =====
	g.GET("/history", middleware.AuthMiddleware(), h.GetHistoryForDay)

	g.PATCH("/orders/:orderID/status", middleware.AuthMiddleware(), h.UpdateStatus)
}
