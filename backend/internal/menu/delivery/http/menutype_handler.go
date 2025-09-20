// delivery/http/menutype_handler.go
package http

import (
	"net/http"
	"strings"

	menu "backend/internal/menu/dto"
	iface "backend/internal/menu/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuTypeHandler struct{ uc iface.MenuTypeUsecase }

func NewMenuTypeHandler(uc iface.MenuTypeUsecase) *MenuTypeHandler {
	return &MenuTypeHandler{uc: uc}
}


// GET /food/menu/restaurant/:restaurantID/types
func (h *MenuTypeHandler) ListByRestaurant(c *gin.Context) {
	rid, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	list, err := h.uc.ListByRestaurant(c.Request.Context(), rid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]menu.MenuTypeResponse, 0, len(list))
	for _, m := range list {
		resp = append(resp, menu.MenuTypeResponse{
			ID:           m.ID,
			RestaurantID: m.RestaurantID,
			Type:         m.Type,
		})
	}

	userID, role := getAuth(c)
	canEdit := (role == "restaurant" && userID != uuid.Nil && userID == rid)

	c.JSON(http.StatusOK, gin.H{
		"can_edit": canEdit,
		"types": resp,
	})
}

// POST /food/menu/restaurant/:restaurantID/types
func (h *MenuTypeHandler) Create(c *gin.Context) {
	roleID, role := getAuth(c)
	if role != "restaurant" || roleID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	rid, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	// ต้องเป็นเจ้าของร้านเท่านั้น
	if roleID != rid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not restaurant owner"})
		return
	}

	var in menu.CreateMenuTypeRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mt, err := h.uc.Create(c.Request.Context(), rid, &in)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, menu.MenuTypeResponse{
		ID:           mt.ID,
		RestaurantID: mt.RestaurantID,
		Type:         mt.Type,
	})
}

// PATCH /food/menu/types/:typeID
func (h *MenuTypeHandler) Update(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type id"})
		return
	}

	actorID, role := getAuth(c)
	if role != "restaurant" || actorID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var in menu.UpdateMenuTypeRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mt, err := h.uc.UpdateOwned(c.Request.Context(), actorID, typeID, &in)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "forbidden"):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, menu.MenuTypeResponse{
		ID:           mt.ID,
		RestaurantID: mt.RestaurantID,
		Type:         mt.Type,
	})
}

// DELETE /food/menu/types/:typeID
func (h *MenuTypeHandler) Delete(c *gin.Context) {
	typeID, err := uuid.Parse(c.Param("typeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type id"})
		return
	}

	actorID, role := getAuth(c)
	if role != "restaurant" || actorID == uuid.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.uc.DeleteOwned(c.Request.Context(), actorID, typeID); err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "forbidden"):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
