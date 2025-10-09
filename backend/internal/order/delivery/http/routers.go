package http

import (
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

// MapFoodOrderRoutes กำหนดเส้นทาง API ของ FoodOrder
func MapFoodOrderRoutes(g *gin.RouterGroup, h *OrderHandler) {
	// ===== Order หลัก =====

	// เปิดออร์เดอร์ใหม่ (ลูกค้าหรือร้านสร้าง)
	g.POST("", middleware.AuthMiddleware(), h.Create)
	g.GET("/:orderID/detail", middleware.AuthMiddleware(), h.GetDetailForRestaurant)

}
