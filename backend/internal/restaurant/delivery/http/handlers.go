package http

import (
	// "net/http"
	// "regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/restaurant/dto"
	"backend/internal/restaurant/interfaces"
	user "backend/internal/user/dto"
)

type RestaurantHandler struct {
	restaurantUsecase interfaces.RestaurantUsecase
}

func NewRestaurantHandler(restaurantUsecase interfaces.RestaurantUsecase) interfaces.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restaurantUsecase,
	}
}

func (h *RestaurantHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerRequest *dto.RegisterRestaurantRequest
		if err := c.ShouldBindJSON(&registerRequest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		err := h.restaurantUsecase.Register(registerRequest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "restaurant registered successfully"})
	}
}

func (h *RestaurantHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var request user.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		token, err := h.restaurantUsecase.Login(&request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": token})
	}
}

func (h *RestaurantHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "User unauthorized"})
			return
		}
		// role, exist := c.Get("role")
		// if !exist || role.(string) != "customer" {
		// 	c.JSON(401, gin.H{"error": "User unauthorized"})
		// 	return
		// }

		restaurants, err := h.restaurantUsecase.GetAll()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"restaurants": restaurants})
	}
}

func (h *RestaurantHandler) UploadProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		file, err := c.FormFile("restaurant_picture")
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed to get file: " + err.Error()})
			return
		}

		const maxFileSize = 3 << 20 // limit to 3MB
		if file.Size > maxFileSize {
			c.JSON(400, gin.H{"error": "File too large. Max allowed is 3MB"})
			return
		}

		url, err := h.restaurantUsecase.UploadProfilePicture(restaurantID, file)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Restaurant profile picture uploaded successfully", "url": url})
	}
}

func (h *RestaurantHandler) ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "User unauthorized"})
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

		var request *dto.ChangeStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = h.restaurantUsecase.ChangeStatus(restaurantID, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Restaurant status changed successfully"})
	}
}

func (h *RestaurantHandler) EditInfo() gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        restID, err := uuid.Parse(idStr)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid restaurant id"})
            return
        }
        var req dto.EditRestaurantRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(401, gin.H{"error": "invalid json body"})
            return
        }
        resp, err := h.restaurantUsecase.EditInfo(restID, &req)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        // c.JSON(200, gin.H{"message": "Edit Restaurant Info changed successfully"})
		c.JSON(200, resp)
    }
}
