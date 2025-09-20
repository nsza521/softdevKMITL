// delivery/http/handlers.go
package http

import (
	"net/http"

	menuuc "backend/internal/menu/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// "log"
)

type MenuHandler struct{ uc menuuc.MenuUsecase }

func NewMenuHandler(uc menuuc.MenuUsecase) *MenuHandler { return &MenuHandler{uc: uc} }

func getAuth(c *gin.Context) (uuid.UUID, string) {
	uid, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var id uuid.UUID
	if uidStr, ok := uid.(string); ok {
		id, _ = uuid.Parse(uidStr)
	} else if uidUUID, ok := uid.(uuid.UUID); ok {
		id = uidUUID
	}

	var roleStr string
	if r, ok := role.(string); ok {
		roleStr = r
	}

	// log ชั่วคราว
	// log.Printf("[handler] user_id=%s role=%s", id.String(), roleStr)
	return id, roleStr
}


// GET /restaurant/menu/:restaurantID/items
func (h *MenuHandler) ListByRestaurant(c *gin.Context) {

	rid, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	CheckRestaurantExistsErr := h.uc.CheckRestaurantExists(c.Request.Context(), rid)
	if CheckRestaurantExistsErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}

	items, err := h.uc.ListByRestaurant(c.Request.Context(), rid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, role := getAuth(c)
	canEdit := (role == "restaurant" && userID != uuid.Nil && userID == rid)

	c.JSON(http.StatusOK, gin.H{
		"canEdit":     canEdit,
		"items":      items,
	})
}

func (h *MenuHandler) Create(c *gin.Context) {
	rid, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	// auth check
	userID, role := getAuth(c)
	
	if role != "restaurant" || userID == uuid.Nil || userID != rid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not restaurant owner",
	})
		return
	}

	var body menuuc.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	resp, err := h.uc.CreateMenuItem(c.Request.Context(), rid, &body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemID"})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, exists := c.Get("role")
	if !exists || role.(string) != "restaurant" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	restaurantID, err := uuid.Parse(userID.(string))

	if role.(string) != "restaurant" || userID == uuid.Nil || err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not restaurant owner"})
		return
	}


	var body menuuc.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	_, err = h.uc.UpdateMenuItem(c.Request.Context(), restaurantID, id, &body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	role, exists := c.Get("role")
	if !exists || role.(string) != "restaurant" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	restaurantID, err := uuid.Parse(userID.(string))

	if role.(string) != "restaurant" || userID == uuid.Nil || err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not restaurant owner"})
		return
	}

	if err := h.uc.DeleteMenuItem(c.Request.Context(), restaurantID, id); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MenuHandler) UploadMenuItemPicture(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Restaurant unauthorized"})
		return
	}
	role, exists := c.Get("role")
	if !exists || role.(string) != "restaurant" {
		c.JSON(401, gin.H{"error": "Restaurant unauthorized"})
		return
	}

	restaurantID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid restaurant ID"})
		return
	}
	menuItemID, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid menu item ID"})
		return
	}

	ctx := c.Request.Context()

	file, err := c.FormFile("menu_item_picture")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}

	const maxFileSize = 3 << 20 // limit to 3MB
	if file.Size > maxFileSize {
		c.JSON(400, gin.H{"error": "File too large. Max allowed is 3MB"})
		return
	}

	url, err := h.uc.UploadMenuItemPicture(ctx, restaurantID, menuItemID, file)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Restaurant profile picture uploaded successfully", "url": url})
}