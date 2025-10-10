// backend/internal/notifications/delivery/http/routers.go
package http

import (
	stdhttp "net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/notifications/dto"
	"backend/internal/notifications/interfaces"
	"backend/internal/middleware"
)

// MapNotiRoutes กำหนดเส้นทาง API สำหรับการแจ้งเตือน
func MapNotiRoutes(r *gin.RouterGroup, h interfaces.NotiHandler) {
	// ✅ GET: List notifications
// 	r.GET("", func(c *gin.Context) {
// 	// 1) รับ receiverId แบบ raw แล้ว strip []" กันเคส Postman แปลก ๆ
// 	raw := c.Query("receiverId")
// 	raw = strings.Trim(raw, "[]\" ")

// 	rid, err := uuid.Parse(raw)
// 	if err != nil {
// 		c.JSON(stdhttp.StatusBadRequest, gin.H{
// 			"error":      "invalid receiverId",
// 			"raw_value":  raw,
// 			"parseError": err.Error(),
// 		})
// 		return
// 	}

// 	// 2) ผูกพารามิเตอร์เอง (เพราะเราเลี่ยง ShouldBindQuery เพื่อกัน UUID พัง)
// 	var q dto.ListQuery
// 	q.ReceiverID = rid

// 	// receiverType จำเป็นต่อ filter (ถ้า repository Where ด้วย receiver_type ด้วย)
// 	q.ReceiverType = c.Query("receiverType") // ใส่มาด้วยใน Postman: customer/restaurant

// 	// isRead เป็น optional (true/false) ถ้า dto ของคุณเป็น *bool
// 	if v := c.Query("isRead"); v != "" {
// 		if b, err := strconv.ParseBool(v); err == nil {
// 			// กรณี dto เป็น *bool:
// 			if p := new(bool); true {
// 				*p = b
// 				q.IsRead = p
// 			}
// 			// ถ้า dto ของคุณเป็น bool ธรรมดา ก็ใช้ q.IsRead = b
// 		}
// 	}

// 	// page/pageSize/sort (default ปลอดภัย)
// 	if q.Page == 0 {
// 		q.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
// 	}
// 	if q.PageSize == 0 {
// 		q.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
// 	}
// 	if q.Sort == "" {
// 		q.Sort = c.DefaultQuery("sort", "created_at_desc")
// 	}

// 	// 3) เรียก usecase
// 	if resp, err := h.List(c.Request.Context(), q); err != nil {
// 		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
// 	} else {
// 		c.JSON(stdhttp.StatusOK, resp)
// 	}
// })

	r.GET("/:receiverId/:page", middleware.AuthMiddleware(), func(c *gin.Context) {
        // Parse receiverId from URL parameter
        receiverIdStr := c.Param("receiverId")
        receiverId, err := uuid.Parse(receiverIdStr)
        if err != nil {
            c.JSON(stdhttp.StatusBadRequest, gin.H{
                "error": "invalid receiverId format",
                "receiverId": receiverIdStr,
            })
            return
        }

        // Parse page from URL parameter
        pageStr := c.Param("page")
        page, err := strconv.Atoi(pageStr)
        if err != nil || page < 1 {
            c.JSON(stdhttp.StatusBadRequest, gin.H{
                "error": "invalid page number",
                "page": pageStr,
            })
            return
        }

        // Get receiverType from query parameter (default: customer)
        receiverType := c.DefaultQuery("type", "customer")
        if receiverType != "customer" && receiverType != "restaurant" {
            c.JSON(stdhttp.StatusBadRequest, gin.H{
                "error": "receiverType must be 'customer' or 'restaurant'",
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
            ReceiverID:   receiverId,
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
	r.POST("/mock", middleware.AuthMiddleware(), func(c *gin.Context) {
		var req dto.MockCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if resp, err := h.MockCreate(c.Request.Context(), req); err != nil {
			c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(stdhttp.StatusCreated, resp)
		}
	})

	// ✅ PATCH: Mark single as read/unread
	r.PATCH("/:id/read", middleware.AuthMiddleware(), func(c *gin.Context) {
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
		} else {
			// c.Status(stdhttp.StatusNoContent)
			c.JSON(stdhttp.StatusOK, gin.H{
				"message": "notification IsRead",
				"noti_Id": id,
			})
		}

	})

	// ✅ PATCH: Mark all as read
	r.PATCH("/read-all", middleware.AuthMiddleware(), func(c *gin.Context) {
	type RequestReadAll struct {
		ReceiverID   string `json:"receiverId" binding:"required"`
		ReceiverType string `json:"receiverType" binding:"required"`
	}

	var req RequestReadAll
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rid, err := uuid.Parse(strings.Trim(req.ReceiverID, "[]\" "))
	if err != nil {
		c.JSON(stdhttp.StatusBadRequest, gin.H{
			"error":      "invalid receiverId",
			"raw_value":  req.ReceiverID,
			"parseError": err.Error(),
		})
		return
	}

	if req.ReceiverType != "customer" && req.ReceiverType != "restaurant" {
		c.JSON(stdhttp.StatusBadRequest, gin.H{"error": "invalid receiverType"})
		return
	}

	updated, err := h.MarkAllRead(c.Request.Context(), rid, req.ReceiverType)
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, gin.H{
		"message": "All notifications marked as read successfully",
		"updated": updated,
		// "rid": rid,
	})
	})

	r.POST("/event", middleware.AuthMiddleware(), func(c *gin.Context) {
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

