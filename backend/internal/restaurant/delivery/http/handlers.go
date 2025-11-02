package http

import (
	"time"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/restaurant/dto"
	"backend/internal/restaurant/interfaces"
)

type RestaurantHandler struct {
	restaurantUsecase interfaces.RestaurantUsecase
}

func NewRestaurantHandler(restaurantUsecase interfaces.RestaurantUsecase) interfaces.RestaurantHandler {
	return &RestaurantHandler{
		restaurantUsecase: restaurantUsecase,
	}
}

func getRestaurantIDAndValidateRole(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}

	role, exist := c.Get("role")
	if !exist || role.(string) != "restaurant" {
		c.JSON(401, gin.H{"error": "restaurant unauthorized"})
		return uuid.Nil, false
	}

	restaurantID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": "invalid user id"})
		return uuid.Nil, false
	}

	return restaurantID, true
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

		var request dto.LoginRequest

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
			c.JSON(401, gin.H{"error": "user unauthorized"})
			return
		}
		role, exist := c.Get("role")
		if !exist || role.(string) != "customer" {
			c.JSON(401, gin.H{"error": "customer unauthorized"})
			return
		}

		restaurants, err := h.restaurantUsecase.GetAll()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"restaurants": restaurants})
	}
}

func (h *RestaurantHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "user unauthorized"})
			return
		}
		role, exist := c.Get("role")
		if !exist || (role.(string) != "customer") {
			c.JSON(401, gin.H{"error": "user unauthorized"})
			return
		}
		
		restaurantIDParam := c.Param("restaurant_id")
		restaurantID, err := uuid.Parse(restaurantIDParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid restaurant id"})
			return
		}	
		detail, err := h.restaurantUsecase.GetByID(restaurantID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"restaurant_detail": detail})
	}
}

func (h *RestaurantHandler) UploadProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {

		restaurantID, ok := getRestaurantIDAndValidateRole(c)
		if !ok {
			return
		}

		file, err := c.FormFile("restaurant_picture")
		if err != nil {
			c.JSON(400, gin.H{"error": "failed to get file: " + err.Error()})
			return
		}

		const maxFileSize = 3 << 20 // limit to 3MB
		if file.Size > maxFileSize {
			c.JSON(400, gin.H{"error": "file too large. Max allowed is 3MB"})
			return
		}

		url, err := h.restaurantUsecase.UploadProfilePicture(restaurantID, file)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "restaurant profile picture uploaded successfully", "url": url})
	}
}

func (h *RestaurantHandler) GetProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID, ok := getRestaurantIDAndValidateRole(c)
		if !ok {
			return
		}
		restaurant, err := h.restaurantUsecase.GetProfilePicture(restaurantID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"profile_picture": restaurant.ProfilePic})
	}
}

func (h *RestaurantHandler) ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID, ok := getRestaurantIDAndValidateRole(c)
		if !ok {
			return
		}

		var request *dto.ChangeStatusRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err := h.restaurantUsecase.ChangeStatus(restaurantID, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "restaurant status changed successfully"})
	}
}

func (h *RestaurantHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, ok := getRestaurantIDAndValidateRole(c)
		if !ok {
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		expClaim, exists := c.Get("exp")
		if !exists {
			c.JSON(500, gin.H{"error": "token expiry not found"})
			return
		}

		expFloat, ok := expClaim.(float64)
		if !ok {
			c.JSON(500, gin.H{"error": "invalid token expiry format"})
			return
		}

		expiry := time.Unix(int64(expFloat), 0)

		err := h.restaurantUsecase.Logout(tokenString, expiry)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "logged out successfully"})
	}
}

func (h *RestaurantHandler) EditInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// idStr := c.Param("id")
		restID, err := uuid.Parse(c.Param("id"))
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

func (h *RestaurantHandler) UpdateName() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := uuid.Parse(c.Param("id"))
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid restaurant id"})
            return
        }

        var req struct {
            Name string `json:"name" binding:"required,min=1"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(401, gin.H{"error": err.Error()})
            return
        }

        updated, err := h.restaurantUsecase.UpdateRestaurantName(c.Request.Context(), id, req.Name)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        c.JSON(200, updated)
    }
}

func (h *RestaurantHandler) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID, ok := getRestaurantIDAndValidateRole(c)
		if !ok {
			return
		}
		balance, err := h.restaurantUsecase.GetBalance(restaurantID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"balance": balance})
	}
}