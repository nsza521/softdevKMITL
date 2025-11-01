package interfaces

import (
	"time"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/db_model"
	"backend/internal/customer/dto"
	// user "backend/internal/user/dto"
)

type CustomerHandler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetProfile() gin.HandlerFunc
	EditProfile() gin.HandlerFunc
	GetFullnameByUsername() gin.HandlerFunc
	GetFirstNameByUsername() gin.HandlerFunc
	GetQRCode() gin.HandlerFunc
	GetMyHistory() gin.HandlerFunc
	
}

type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	Create(customer *models.Customer) error
	IsCustomerExists(username string, email string) (bool, error)
	GetByID(id uuid.UUID) (*models.Customer, error)
	Update(customer *models.Customer) error
	ListServedOrdersByCustomer(context.Context, string) ([]models.FoodOrder, error)
}

type CustomerUsecase interface {
	Register(request *dto.RegisterCustomerRequest) error
	Login(request *dto.LoginRequest) (string, error)
	GetProfile(customerID uuid.UUID) (*dto.ProfileResponse, error)
	EditProfile(customerID uuid.UUID, request *dto.EditProfileRequest) error
	GetFullnameByUsername(customerID uuid.UUID, request *dto.GetFullnameRequest) (*dto.GetFullnameResponse, error)
	GetFirstnameByUsername(customerID uuid.UUID, request *dto.GetFullnameRequest) (*dto.GetFirstnameResponse, error)
	GetQRCode(customerID uuid.UUID) (string, error)
	Logout(token string, expiry time.Time) error
	GetMyOrderHistory(ctx *gin.Context) ([]dto.OrderHistoryDay, error)
	ListServedOrdersByCustomer(ctx context.Context, customerID string) ([]models.FoodOrder, error)

}