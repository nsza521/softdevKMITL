package http

import (
	orderif "backend/internal/order/interfaces"
	"backend/internal/order/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FoodOrderHandler struct {
	uc orderif.FoodOrderUsecase
}

func NewFoodOrderHandler(uc orderif.FoodOrderUsecase) *FoodOrderHandler {
	return &FoodOrderHandler{uc: uc}
}

// POST /orders
func (h *FoodOrderHandler) Create(c *gin.Context) {
	var req dto.CreateFoodOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentUser := getUserIDFromJWT(c)
	resp, err := h.uc.Create(c, req, currentUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, resp)
}

// GET /orders/:orderID
func (h *FoodOrderHandler) GetDetail(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid order id"}); return }
	resp, err := h.uc.GetDetail(c, orderID)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, resp)
}

// POST /orders/:orderID/items
func (h *FoodOrderHandler) AppendItems(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid order id"}); return }

	var req dto.AppendItemsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	currentUser := getUserIDFromJWT(c)

	resp, err := h.uc.AppendItems(c, orderID, req, currentUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, resp)
}

// DELETE /orders/:orderID/items/:itemID
func (h *FoodOrderHandler) RemoveItem(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid order id"}); return }
	itemID, err := uuid.Parse(c.Param("itemID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid item id"}); return }
	currentUser := getUserIDFromJWT(c)

	resp, err := h.uc.RemoveItem(c, orderID, itemID, currentUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, resp)
}

// POST /orders/:orderID/customers
func (h *FoodOrderHandler) AttachCustomer(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid order id"}); return }
	var req dto.AttachCustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	currentUser := getUserIDFromJWT(c)

	resp, err := h.uc.AttachCustomer(c, orderID, req, currentUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, resp)
}

// GET /orders/:orderID/customers
func (h *FoodOrderHandler) ListCustomers(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid order id"}); return }
	currentUser := getUserIDFromJWT(c)

	resp, err := h.uc.ListCustomers(c, orderID, currentUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, resp)
}

// helper สมมติ: ดึง user id จาก JWT ที่ middleware เซ็ตไว้
func getUserIDFromJWT(c *gin.Context) uuid.UUID {
	// ตัวอย่างง่าย ๆ: สมมติ middleware เซ็ต "user_id" เป็น uuid.UUID ลงใน context
	if v, ok := c.Get("user_id"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 { return id }
	}
	// fallback (ไม่ควรใช้จริง)
	return uuid.Nil
}
