package http

import (
	"backend/internal/menu/dto"
	"backend/internal/menu/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddOnHandler struct {
	uc interfaces.AddOnUsecase
}

func NewAddOnHandler(uc interfaces.AddOnUsecase) *AddOnHandler {
	return &AddOnHandler{uc: uc}
}

// POST /restaurants/:restaurantID/addon-groups
func (h *AddOnHandler) CreateGroup(c *gin.Context) {
	restID, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	UserID, userIDExists := c.Get("user_id")
	Role, roleExists := c.Get("role")

	if !userIDExists || !roleExists || Role != "restaurant" || UserID == nil || UserID != restID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}		

	var req menu.CreateAddOnGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group, err := h.uc.CreateGroup(restID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, group)
}

// GET /restaurants/:restaurantID/addon-groups
func (h *AddOnHandler) ListGroups(c *gin.Context) {
	restID, err := uuid.Parse(c.Param("restaurantID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id"})
		return
	}

	UserID, userIDExists := c.Get("user_id")
	Role, roleExists := c.Get("role")

	if !userIDExists || !roleExists || Role != "restaurant" || UserID == nil || UserID != restID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}		


	groups, err := h.uc.ListGroups(restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// PUT /addon-groups/:id
func (h *AddOnHandler) UpdateGroup(c *gin.Context) {
	
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req menu.UpdateAddOnGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.UpdateGroup(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE /addon-groups/:id
func (h *AddOnHandler) DeleteGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.uc.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GET /options/:optionID
func (h *AddOnHandler) GetOption(c *gin.Context) {
    id, err := uuid.Parse(c.Param("optionID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid option id"})
        return
    }
    opt, err := h.uc.GetOption(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, opt)
}

// GET /addon-groups/:groupID/options
func (h *AddOnHandler) ListOptions(c *gin.Context) {
    gid, err := uuid.Parse(c.Param("groupID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
        return
    }
    opts, err := h.uc.ListOptions(gid)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, opts)
}



// POST /addon-groups/:groupID/options
func (h *AddOnHandler) CreateOption(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	var req menu.CreateAddOnOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	opt, err := h.uc.CreateOption(groupID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, opt)
}

// PUT /options/:optionID
func (h *AddOnHandler) UpdateOption(c *gin.Context) {
	id, err := uuid.Parse(c.Param("optionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req menu.UpdateAddOnOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.UpdateOption(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE /options/:optionID
func (h *AddOnHandler) DeleteOption(c *gin.Context) {
	id, err := uuid.Parse(c.Param("optionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.uc.DeleteOption(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
