package http

import (
	"github.com/gin-gonic/gin"
	"backend/internal/middleware"
)

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/:restaurantID/items", middleware.AuthMiddleware(),h.ListByRestaurant)
	g.POST("/:restaurantID/items", middleware.AuthMiddleware(), h.Create) // middleware.RequireRole("restaurant_owner"),
	g.PATCH("/items/:itemID", middleware.AuthMiddleware(), h.Update)      // middleware.RequireRole("restaurant_owner"),
	g.DELETE("/items/:itemID", middleware.AuthMiddleware(), h.Delete)    // middleware.RequireRole("restaurant_owner"),
	g.POST("/items/:itemID/upload_pic", middleware.AuthMiddleware(), h.UploadMenuItemPicture) // middleware.RequireRole("restaurant_owner"),

}


// MapMenuTypeRoutes กำหนดเส้นทาง CRUD ของ MenuType
func MapMenuTypeRoutes(g *gin.RouterGroup, h *MenuTypeHandler) {
	// List/Create ใช้ restaurantID ใน path
	g.GET("/:restaurantID/types", middleware.AuthMiddleware(), h.ListByRestaurant)
	g.POST("/:restaurantID/types", middleware.AuthMiddleware(), h.Create)

	// Update/Delete ใช้ typeID โดยตรง ไม่ต้อง query restaurantID
	g.PATCH("/types/:typeID", middleware.AuthMiddleware(), h.Update)
	g.DELETE("/types/:typeID", middleware.AuthMiddleware(), h.Delete)
}
