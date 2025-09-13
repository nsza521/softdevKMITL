package interfaces

import (
	// "backend/internal/db_model"
	"github.com/google/uuid"
	"backend/internal/customer/dto"
	user "backend/internal/user/dto"
)

type CustomerUsecase interface {
	Register(request *dto.RegisterCustomerRequest) error
	Login(request *user.LoginRequest) (string, error)
	GetProfile(customerID uuid.UUID) (*dto.ProfileResponse, error)
	EditProfile(customerID uuid.UUID, request *dto.EditProfileRequest) error
	Logout() error
}