package interfaces

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/db_model"
	"backend/internal/customer/dto"
	user "backend/internal/user/dto"
)

type CustomerHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetProfile() gin.HandlerFunc
	EditProfile() gin.HandlerFunc
	GetFullnameByUsername() gin.HandlerFunc
}

type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	Create(customer *models.Customer) error
	IsCustomerExists(username string, email string) (bool, error)
	GetByID(id uuid.UUID) (*models.Customer, error)
	Update(customer *models.Customer) error
}

type CustomerUsecase interface {
	Register(request *dto.RegisterCustomerRequest) error
	Login(request *user.LoginRequest) (string, error)
	GetProfile(customerID uuid.UUID) (*dto.ProfileResponse, error)
	EditProfile(customerID uuid.UUID, request *dto.EditProfileRequest) error
	GetFullnameByUsername(customerID uuid.UUID, request *dto.GetFullnameRequest) (*dto.GetFullnameResponse, error)
	Logout(token string, expiry time.Time) error
}