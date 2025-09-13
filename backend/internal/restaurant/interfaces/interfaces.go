package interfaces

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/restaurant/dto"
	"backend/internal/db_model"
	user "backend/internal/user/dto"
)

type RestaurantHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	GetAll() gin.HandlerFunc
}

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) (*models.Restaurant, error)
	IsRestaurantExists(username string, email string) (bool, error)
	GetByUsername(username string) (*models.Restaurant, error)
	GetByID(id uuid.UUID) (*models.Restaurant, error)
	GetAll() ([]*models.Restaurant, error)
	CreateBankAccount(bankAccount *models.BankAccount) error

	ExistsByID(ctx context.Context, id string) (bool, error)
	IsOwner(ctx context.Context, restaurantID, userID string) (bool, error)
	PartialUpdate(ctx context.Context, id string, changes map[string]any) error
}

type RestaurantUsecase interface {
	Register(request *dto.RegisterRestaurantRequest) error
	Login(request *user.LoginRequest) (string, error)
	GetAll() ([]dto.RestaurantDetailResponse, error)
	EditRestaurant(ctx context.Context, id, userID, role string, req dto.EditRestaurantRequest) (dto.RestaurantDetailResponse, error)
}