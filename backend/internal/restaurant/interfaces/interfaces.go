package interfaces

import (
	"time"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/restaurant/dto"
	"backend/internal/db_model"
	user "backend/internal/user/dto"
)

type RestaurantHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetAll() gin.HandlerFunc
	UploadProfilePicture() gin.HandlerFunc
	ChangeStatus() gin.HandlerFunc
}

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) (*models.Restaurant, error)
	IsRestaurantExists(username string, email string) (bool, error)
	GetByUsername(username string) (*models.Restaurant, error)
	GetByID(id uuid.UUID) (*models.Restaurant, error)
	GetAll() ([]*models.Restaurant, error)
	CreateBankAccount(bankAccount *models.BankAccount) error
	Update(restaurant *models.Restaurant) error
}

type RestaurantUsecase interface {
	Register(request *dto.RegisterRestaurantRequest) error
	Login(request *user.LoginRequest) (string, error)
	Logout(token string, expiry time.Time) error
	GetAll() ([]dto.RestaurantDetailResponse, error)
	UploadProfilePicture(restaurantID uuid.UUID, file *multipart.FileHeader) (string, error)
	ChangeStatus(restaurantID uuid.UUID, request *dto.ChangeStatusRequest) error
}