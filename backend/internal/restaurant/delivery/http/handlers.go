package http

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"

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

		// userID, exists := c.Get("user_id")
		// if !exists {
		// 	c.JSON(401, gin.H{"error": "User unauthorized"})
		// 	return
		// }
		// _, err := uuid.Parse(userID.(string))
		// if err != nil {
		// 	c.JSON(401, gin.H{"error": "Invalid user ID"})
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

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type EditHandler struct {
	Uc interfaces.RestaurantUsecase
}

func NewEditHandler(uc interfaces.RestaurantUsecase) *EditHandler { return &EditHandler{Uc: uc} }

func (h *EditHandler) EditRestaurant(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	var req dto.EditRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	// at least it should had 1 field
	if req.Username == nil && req.Email == nil && req.BankName == nil && req.AccountNumber == nil && req.AccountName == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if req.Email != nil && !emailRe.MatchString(*req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	resp, err := h.Uc.EditRestaurant(c, id, userID, role, req)
	if err != nil {
		switch interfaces.AsHTTPStatus(err) {
		case http.StatusForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case http.StatusNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}