package http

import (
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func MapMenuRoutes(g *gin.RouterGroup, h *MenuHandler) {
	g.GET("/:restaurantID/items", middleware.AuthMiddleware(), h.ListByRestaurant)

	//	g.GET("/items/:itemID/detail", middleware.AuthMiddleware(), h.GetDetail)
	g.GET("/:restaurantID/:itemID/detail", middleware.AuthMiddleware(), h.GetDetail)

	g.POST("/:restaurantID/items", middleware.AuthMiddleware(), h.Create)                     // middleware.RequireRole("restaurant_owner"),
	g.PATCH("/items/:itemID", middleware.AuthMiddleware(), h.Update)                          // middleware.RequireRole("restaurant_owner"),
	g.DELETE("/items/:itemID", middleware.AuthMiddleware(), h.Delete)                         // middleware.RequireRole("restaurant_owner"),
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

func MapAddOnRoutes(g *gin.RouterGroup, h *AddOnHandler) {
	// ===== AddOnGroup =====
	// List/Create group โดยใช้ restaurantID
	g.GET("/:restaurantID/addon-groups", middleware.AuthMiddleware(), h.ListGroups)
	g.POST("/:restaurantID/addon-groups", middleware.AuthMiddleware(), h.CreateGroup) // + middleware.RequireRole("restaurant_owner")

	// Update/Delete group ใช้ groupID โดยตรง
	g.PATCH("/addon-groups/:groupID", middleware.AuthMiddleware(), h.UpdateGroup)
	g.DELETE("/addon-groups/:groupID", middleware.AuthMiddleware(), h.DeleteGroup)

	// POST   /restaurant/menu/addon-groups/:groupID/types/:typeID   (link)
	g.POST("/addon-groups/:groupID/types/:typeID", middleware.AuthMiddleware(), h.LinkGroupToType)
	// DELETE /restaurant/menu/addon-groups/:groupID/types/:typeID   (unlink)
	g.DELETE("/addon-groups/:groupID/types/:typeID", middleware.AuthMiddleware(), h.UnlinkGroupFromType)

	// ===== AddOnOption =====
	// Create option โดยอิง groupID
	g.POST("/addon-groups/:groupID/options", middleware.AuthMiddleware(), h.CreateOption)

	// Update/Delete option ใช้ optionID โดยตรง
	g.PATCH("/options/:optionID", middleware.AuthMiddleware(), h.UpdateOption)
	g.GET("/addon-groups/:groupID/options", middleware.AuthMiddleware(), h.ListOptions) // ✅
	g.GET("/options/:optionID", middleware.AuthMiddleware(), h.GetOption)

	g.DELETE("/options/:optionID", middleware.AuthMiddleware(), h.DeleteOption)
}
