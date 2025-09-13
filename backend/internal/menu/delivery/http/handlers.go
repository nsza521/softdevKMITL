package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	menuuc "backend/internal/menu/interfaces"
)

type MenuHandler struct{ uc menuuc.MenuUsecase }

func NewMenuHandler(uc menuuc.MenuUsecase) *MenuHandler { return &MenuHandler{uc: uc} }

// GET /food/menu/restaurant/:restaurantID/items
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

	out, err := h.uc.ListByRestaurant(c.Request.Context(), rid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}


func (h *MenuHandler) Create(c *gin.Context) {
	rid, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid restaurant id"}); return }

	var body menuuc.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"})
		return
	}
	resp, err := h.uc.CreateMenuItem(c.Request.Context(), rid, &body)
	if err != nil { c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, resp)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("itemID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }

	var body menuuc.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid body"})
		return
	}
	_, err = h.uc.UpdateMenuItem(c.Request.Context(), id, &body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("itemID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }
	if err := h.uc.DeleteMenuItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

