// backend/internal/notifications/delivery/http/routers.go
package http

import (
	stdhttp "net/http"
	"strconv"

	// "strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/middleware"
	"backend/internal/notifications/dto"
	"backend/internal/notifications/interfaces"
)

func getCustomerIDAndValidateRole(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	role, exist := c.Get("role")
	if !exist || role.(string) != "customer" {
		c.JSON(401, gin.H{"error": "customer unauthorized"})
		return uuid.Nil, false
	}

	customerID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "invalid user id"})
		return uuid.Nil, false
	}

	return customerID, true
}

// MapNotiRoutes กำหนดเส้นทาง API สำหรับการแจ้งเตือน
func MapNotiRoutes(r *gin.RouterGroup, h interfaces.NotiHandler) {

	r.Use(middleware.AuthMiddleware())

	// r.GET("/:receiverId/:page", func(c *gin.Context) {
	r.GET("/:page", func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}
		// // Parse receiverId from URL parameter
		// receiverIdStr := c.Param("receiverId")
		// receiverId, err := uuid.Parse(receiverIdStr)
		// if err != nil {
		//     c.JSON(stdhttp.StatusBadRequest, gin.H{
		//         "error": "invalid receiverId format",
		//         "receiverId": receiverIdStr,
		//     })
		//     return
		// }

		// Parse page from URL parameter
		pageStr := c.Param("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			c.JSON(stdhttp.StatusBadRequest, gin.H{
				"error": "invalid page number",
				"page":  pageStr,
			})
			return
		}

		// Get receiverType from query parameter (default: customer)
		receiverType := c.DefaultQuery("type", "customer")
		if receiverType != "customer" && receiverType != "restaurant" {
			c.JSON(stdhttp.StatusBadRequest, gin.H{
				"error":        "receiverType must be 'customer' or 'restaurant'",
				"receiverType": receiverType,
			})
			return
		}

		// Optional: filter by read status
		var isRead *bool
		if readStatus := c.Query("read"); readStatus != "" {
			if b, err := strconv.ParseBool(readStatus); err == nil {
				isRead = &b
			}
		}

		// Set page size (default: 20, max: 100)
		pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
		if pageSize > 100 {
			pageSize = 100
		}

		// Build query
		q := dto.ListQuery{
			ReceiverID:   customerID,
			ReceiverType: receiverType,
			IsRead:       isRead,
			Page:         page,
			PageSize:     pageSize,
			Sort:         c.DefaultQuery("sort", "created_at_desc"),
		}

		// Call usecase
		resp, err := h.List(c.Request.Context(), q)
		if err != nil {
			c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(stdhttp.StatusOK, resp)
	})

	// ✅ POST: Mock create notifications
	r.POST("/mock", func(c *gin.Context) {

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		var body struct {
			Count        int    `json:"count" binding:"required"`
			ReceiverType string `json:"receiverType" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// var req dto.MockCreateRequest
		req := dto.MockCreateRequest{
			Count:        body.Count,
			ReceiverID:   customerID,
			ReceiverType: body.ReceiverType,
		}
		resp, err := h.MockCreate(c.Request.Context(), req)
		if err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(stdhttp.StatusCreated, resp)
		// if err := c.ShouldBindJSON(&req); err != nil {
		// 	c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// if resp, err := h.MockCreate(c.Request.Context(), req); err != nil {
		// 	c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		// } else {
		// 	c.JSON(stdhttp.StatusCreated, resp)
		// }
	})

	// ✅ PATCH: Mark single as read/unread
	r.PATCH("/:id/read", func(c *gin.Context) {
		// require customer token
		_, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req dto.MarkReadRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.MarkRead(c.Request.Context(), id, req.IsRead); err != nil {
			c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			// c.Status(stdhttp.StatusNoContent)
			c.JSON(stdhttp.StatusOK, gin.H{
				"message": "notification IsRead",
				"noti_Id": id,
			})
		}

	})

	// ✅ PATCH: Mark all as read
	r.PATCH("/read-all", func(c *gin.Context) {
		type RequestReadAll struct {
			// ReceiverID   string `json:"receiverId" binding:"required"`
			ReceiverType string `json:"receiverType" binding:"required"`
		}

		var req RequestReadAll
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// rid, err := uuid.Parse(strings.Trim(req.ReceiverID, "[]\" "))
		// if err != nil {
		// 	c.JSON(stdhttp.StatusBadRequest, gin.H{
		// 		"error":      "invalid receiverId",
		// 		"raw_value":  req.ReceiverID,
		// 		"parseError": err.Error(),
		// 	})
		// 	return
		// }

		customerID, ok := getCustomerIDAndValidateRole(c)
		if !ok {
			return
		}

		// if req.ReceiverType != "customer" && req.ReceiverType != "restaurant" {
		// 	c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "invalid receiverType"})
		// 	return
		// }

		updated, err := h.MarkAllRead(c.Request.Context(), customerID, req.ReceiverType)
		if err != nil {
			c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(stdhttp.StatusOK, gin.H{
			"message": "All notifications marked as read successfully",
			"updated": updated,
		})
	})

	r.POST("/event", func(c *gin.Context) {
		var req dto.CreateEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := h.CreateFromEvent(c.Request.Context(), req)
		if err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(stdhttp.StatusCreated, resp)
	})

}
