package interfaces

import (
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
}

type RestaurantUsecase interface {
	Register(request *dto.RegisterRestaurantRequest) error
	Login(request *user.LoginRequest) (string, error)
	GetAll() ([]dto.RestaurantDetailResponse, error)
}