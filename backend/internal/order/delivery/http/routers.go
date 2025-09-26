package http

import (
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

// MapFoodOrderRoutes กำหนดเส้นทาง API ของ FoodOrder
func MapFoodOrderRoutes(g *gin.RouterGroup, h *FoodOrderHandler) {
	// ===== Order หลัก =====

	// เปิดออร์เดอร์ใหม่ (ลูกค้าหรือร้านสร้าง)
	g.POST("", middleware.AuthMiddleware(), h.Create)

	// อ่านรายละเอียดออร์เดอร์ (ลูกค้าใน reservation หรือร้านดูได้)
	g.GET("/:orderID", middleware.AuthMiddleware(), h.GetDetail)

	// ===== Order Items =====

	// เพิ่มเมนูเข้าออร์เดอร์ที่มีอยู่แล้ว
	g.POST("/:orderID/items", middleware.AuthMiddleware(), h.AppendItems)

	// (ถ้าต้องการ) ลบเมนูออกจากออร์เดอร์
	g.DELETE("/:orderID/items/:itemID", middleware.AuthMiddleware(), h.RemoveItem)

	// ===== Customers (สำหรับ join reservation) =====

	// เพิ่มลูกค้าเข้าร่วมออร์เดอร์ (ใน reservation เดียวกัน)
	g.POST("/:orderID/customers", middleware.AuthMiddleware(), h.AttachCustomer)

	// ลิสต์ลูกค้าที่อยู่ในออร์เดอร์
	g.GET("/:orderID/customers", middleware.AuthMiddleware(), h.ListCustomers)
}
